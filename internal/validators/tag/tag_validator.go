package tag

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"

	"github.com/spectrocloud-labs/validator-plugin-aws/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-aws/internal/constants"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	vapiconstants "github.com/spectrocloud-labs/validator/pkg/constants"
	vapitypes "github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/spectrocloud-labs/validator/pkg/util/ptr"
)

type tagApi interface {
	DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
}

type TagRuleService struct {
	log    logr.Logger
	tagSvc tagApi
}

func NewTagRuleService(log logr.Logger, tagSvc tagApi) *TagRuleService {
	return &TagRuleService{
		log:    log,
		tagSvc: tagSvc,
	}
}

// ReconcileTagRule reconciles an EC2 tagging validation rule from the AWSValidator config
func (s *TagRuleService) ReconcileTagRule(rule v1alpha1.TagRule) (*vapitypes.ValidationResult, error) {

	// Build the default latest condition for this tag rule
	state := vapi.ValidationSucceeded
	latestCondition := vapi.DefaultValidationCondition()
	latestCondition.Message = "All required subnet tags were found"
	latestCondition.ValidationRule = fmt.Sprintf("%s-%s-%s", vapiconstants.ValidationRulePrefix, rule.ResourceType, rule.Key)
	latestCondition.ValidationType = constants.ValidationTypeTag
	validationResult := &vapitypes.ValidationResult{Condition: &latestCondition, State: &state}

	switch rule.ResourceType {
	case "subnet":
		// match the tag rule's list of ARNs against the subnets with tag 'rule.Key=rule.ExpectedValue'
		failures := make([]string, 0)
		foundArns := make(map[string]bool)
		subnets, err := s.tagSvc.DescribeSubnets(context.Background(), &ec2.DescribeSubnetsInput{
			Filters: []ec2types.Filter{
				{
					Name:   ptr.Ptr(fmt.Sprintf("tag:%s", rule.Key)),
					Values: []string{rule.ExpectedValue},
				},
			},
		})
		if err != nil {
			s.log.V(0).Error(err, "failed to describe subnets", "region", rule.Region)
			return validationResult, err
		}
		for _, s := range subnets.Subnets {
			if s.SubnetArn != nil {
				foundArns[*s.SubnetArn] = true
			}
		}
		for _, arn := range rule.ARNs {
			_, ok := foundArns[arn]
			if !ok {
				failures = append(failures, fmt.Sprintf("Subnet with ARN %s missing tag %s=%s", arn, rule.Key, rule.ExpectedValue))
			}
		}
		if len(failures) > 0 {
			state = vapi.ValidationFailed
			latestCondition.Failures = failures
			latestCondition.Message = "One or more required subnet tags was not found"
			latestCondition.Status = corev1.ConditionFalse
		}
	default:
		return nil, fmt.Errorf("unsupported resourceType %s for TagRule", rule.ResourceType)
	}

	return validationResult, nil
}
