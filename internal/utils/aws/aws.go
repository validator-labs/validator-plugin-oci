package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
)

type AwsApi struct {
	IAM   *iam.Client
	EC2   *ec2.Client
	EFS   *efs.Client
	ELB   *elasticloadbalancing.Client
	ELBV2 *elasticloadbalancingv2.Client
	SQ    *servicequotas.Client
}

// NewAwsApi creates an AwsApi object that aggregates AWS service clients
func NewAwsApi(region string) (*AwsApi, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &AwsApi{
		IAM:   iam.NewFromConfig(cfg),
		EC2:   ec2.NewFromConfig(cfg),
		EFS:   efs.NewFromConfig(cfg),
		ELB:   elasticloadbalancing.NewFromConfig(cfg),
		ELBV2: elasticloadbalancingv2.NewFromConfig(cfg),
		SQ:    servicequotas.NewFromConfig(cfg),
	}, nil
}
