package servicequota

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	sqtypes "github.com/aws/aws-sdk-go-v2/service/servicequotas/types"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"

	"github.com/spectrocloud-labs/validator-plugin-aws/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-aws/internal/utils/test"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	"github.com/spectrocloud-labs/validator/pkg/types"
	"github.com/spectrocloud-labs/validator/pkg/util/ptr"
)

type ec2ApiMock struct {
	addresses         *ec2.DescribeAddressesOutput
	images            *ec2.DescribeImagesOutput
	internetGateways  *ec2.DescribeInternetGatewaysOutput
	natGateways       *ec2.DescribeNatGatewaysOutput
	networkInterfaces *ec2.DescribeNetworkInterfacesOutput
	subnets           *ec2.DescribeSubnetsOutput
	vpcs              *ec2.DescribeVpcsOutput
}

func (m ec2ApiMock) DescribeAddresses(ctx context.Context, params *ec2.DescribeAddressesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeAddressesOutput, error) {
	return m.addresses, nil
}

func (m ec2ApiMock) DescribeImages(ctx context.Context, params *ec2.DescribeImagesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error) {
	return m.images, nil
}

func (m ec2ApiMock) DescribeInternetGateways(ctx context.Context, params *ec2.DescribeInternetGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInternetGatewaysOutput, error) {
	return m.internetGateways, nil
}

func (m ec2ApiMock) DescribeNetworkInterfaces(ctx context.Context, params *ec2.DescribeNetworkInterfacesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNetworkInterfacesOutput, error) {
	return m.networkInterfaces, nil
}

func (m ec2ApiMock) DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error) {
	return m.subnets, nil
}

func (m ec2ApiMock) DescribeVpcs(ctx context.Context, params *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error) {
	return m.vpcs, nil
}

func (m ec2ApiMock) DescribeNatGateways(ctx context.Context, params *ec2.DescribeNatGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeNatGatewaysOutput, error) {
	return m.natGateways, nil
}

type efsApiMock struct {
	filesystems *efs.DescribeFileSystemsOutput
}

func (m efsApiMock) DescribeFileSystems(ctx context.Context, params *efs.DescribeFileSystemsInput, optFns ...func(*efs.Options)) (*efs.DescribeFileSystemsOutput, error) {
	return m.filesystems, nil
}

type elbApiMock struct {
	loadBalancers *elasticloadbalancing.DescribeLoadBalancersOutput
}

func (m elbApiMock) DescribeLoadBalancers(context.Context, *elasticloadbalancing.DescribeLoadBalancersInput, ...func(*elasticloadbalancing.Options)) (*elasticloadbalancing.DescribeLoadBalancersOutput, error) {
	return m.loadBalancers, nil
}

type elbv2ApiMock struct {
	loadBalancers *elasticloadbalancingv2.DescribeLoadBalancersOutput
}

func (m elbv2ApiMock) DescribeLoadBalancers(context.Context, *elasticloadbalancingv2.DescribeLoadBalancersInput, ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeLoadBalancersOutput, error) {
	return m.loadBalancers, nil
}

type sqApiMock struct {
	serviceQuotas *servicequotas.ListServiceQuotasOutput
}

func (m sqApiMock) ListServiceQuotas(context.Context, *servicequotas.ListServiceQuotasInput, ...func(*servicequotas.Options)) (*servicequotas.ListServiceQuotasOutput, error) {
	return m.serviceQuotas, nil
}

var svcQuotaService = NewServiceQuotaRuleService(
	logr.Logger{},
	ec2ApiMock{
		addresses: &ec2.DescribeAddressesOutput{
			Addresses: []ec2types.Address{
				{
					AssociationId: ptr.Ptr("1"),
				},
				{
					AssociationId: nil,
				},
			},
		},
		images: &ec2.DescribeImagesOutput{
			Images: []ec2types.Image{
				{
					ImageId: ptr.Ptr("1"),
				},
			},
		},
		internetGateways: &ec2.DescribeInternetGatewaysOutput{
			InternetGateways: []ec2types.InternetGateway{
				{
					InternetGatewayId: ptr.Ptr("1"),
				},
			},
		},
		natGateways: &ec2.DescribeNatGatewaysOutput{
			NatGateways: []ec2types.NatGateway{
				{
					NatGatewayId: ptr.Ptr("1"),
					SubnetId:     ptr.Ptr("1"),
				},
			},
		},
		networkInterfaces: &ec2.DescribeNetworkInterfacesOutput{
			NetworkInterfaces: []ec2types.NetworkInterface{
				{
					NetworkInterfaceId: ptr.Ptr("1"),
				},
			},
		},
		subnets: &ec2.DescribeSubnetsOutput{
			Subnets: []ec2types.Subnet{
				{
					SubnetId:         ptr.Ptr("1"),
					AvailabilityZone: ptr.Ptr("us-west-1"),
					VpcId:            ptr.Ptr("1"),
				},
				{
					SubnetId:         ptr.Ptr("2"),
					AvailabilityZone: ptr.Ptr("us-west-2"),
					VpcId:            ptr.Ptr("2"),
				},
				{
					SubnetId:         ptr.Ptr("3"),
					AvailabilityZone: ptr.Ptr("us-west-1"),
					VpcId:            ptr.Ptr("1"),
				},
			},
		},
		vpcs: &ec2.DescribeVpcsOutput{
			Vpcs: []ec2types.Vpc{
				{
					VpcId: ptr.Ptr("1"),
				},
			},
		},
	},
	efsApiMock{
		filesystems: &efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{
				{
					FileSystemId: ptr.Ptr("1"),
				},
			},
		},
	},
	elbApiMock{
		loadBalancers: &elasticloadbalancing.DescribeLoadBalancersOutput{
			LoadBalancerDescriptions: []elbtypes.LoadBalancerDescription{
				{
					LoadBalancerName: ptr.Ptr("clb1"),
				},
			},
		},
	},
	elbv2ApiMock{
		loadBalancers: &elasticloadbalancingv2.DescribeLoadBalancersOutput{
			LoadBalancers: []elbv2types.LoadBalancer{
				{
					LoadBalancerName: ptr.Ptr("alb1"),
					Type:             elbv2types.LoadBalancerTypeEnumApplication,
				},
				{
					LoadBalancerName: ptr.Ptr("nlb1"),
					Type:             elbv2types.LoadBalancerTypeEnumNetwork,
				},
			},
		},
	},
	sqApiMock{
		serviceQuotas: mockQuotas,
	},
)

var mockQuotas = &servicequotas.ListServiceQuotasOutput{}

type testCase struct {
	name           string
	rule           v1alpha1.ServiceQuotaRule
	expectedResult types.ValidationResult
	expectedError  error
	mockQuotas     []sqtypes.ServiceQuota
}

func TestQuotaValidation(t *testing.T) {
	cs := []testCase{
		{
			name: "Fail (insufficient EIPs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "EC2-VPC Elastic IPs",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("EC2-VPC Elastic IPs"),
					Value:     ptr.Ptr(1.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"EC2-VPC Elastic IPs: quota: 1, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service ec2 and quota EC2-VPC Elastic IPs"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient EIPs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "EC2-VPC Elastic IPs",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("EC2-VPC Elastic IPs"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"EC2-VPC Elastic IPs: quota: 5, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (insufficient AMIs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Public AMIs",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Public AMIs"),
					Value:     ptr.Ptr(1.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"Public AMIs: quota: 1, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service ec2 and quota Public AMIs"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient AMIs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Public AMIs",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Public AMIs"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"Public AMIs: quota: 5, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (insufficient IGs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Internet gateways per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Internet gateways per Region"),
					Value:     ptr.Ptr(1.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"Internet gateways per Region: quota: 1, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service ec2 and quota Internet gateways per Region"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient IGs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Internet gateways per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Internet gateways per Region"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"Internet gateways per Region: quota: 5, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (insufficient NICs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Network interfaces per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Network interfaces per Region"),
					Value:     ptr.Ptr(1.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"Network interfaces per Region: quota: 1, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service ec2 and quota Network interfaces per Region"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient NICs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Network interfaces per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Network interfaces per Region"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"Network interfaces per Region: quota: 5, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (insufficient VPCs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "VPCs per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("VPCs per Region"),
					Value:     ptr.Ptr(1.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"VPCs per Region: quota: 1, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service ec2 and quota VPCs per Region"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient VPCs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "VPCs per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("VPCs per Region"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"VPCs per Region: quota: 5, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (insufficient Subnets)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Subnets per VPC",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Subnets per VPC"),
					Value:     ptr.Ptr(2.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"Subnets per VPC: quota: 2, buffer: 3, max. usage: 2, max. usage entity: 1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service ec2 and quota Subnets per VPC"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient Subnets)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Subnets per VPC",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Subnets per VPC"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"Subnets per VPC: quota: 5, buffer: 3, max. usage: 2, max. usage entity: 1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (insufficient NGs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "NAT gateways per Availability Zone",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("NAT gateways per Availability Zone"),
					Value:     ptr.Ptr(1.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"NAT gateways per Availability Zone: quota: 1, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service ec2 and quota NAT gateways per Availability Zone"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient NGs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "ec2",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "NAT gateways per Availability Zone",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("NAT gateways per Availability Zone"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-ec2",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"NAT gateways per Availability Zone: quota: 5, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (insufficient Filesystems)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "efs",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "File systems per account",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("File systems per account"),
					Value:     ptr.Ptr(1.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-efs",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"File systems per account: quota: 1, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service efs and quota File systems per account"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient Filesystems)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "efs",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "File systems per account",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("File systems per account"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-efs",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"File systems per account: quota: 5, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (insufficient CLBs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "elasticloadbalancing",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Classic Load Balancers per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Classic Load Balancers per Region"),
					Value:     ptr.Ptr(1.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-elasticloadbalancing",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"Classic Load Balancers per Region: quota: 1, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service elasticloadbalancing and quota Classic Load Balancers per Region"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient CLBs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "elasticloadbalancing",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Classic Load Balancers per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Classic Load Balancers per Region"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-elasticloadbalancing",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"Classic Load Balancers per Region: quota: 5, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (insufficient ALBs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "elasticloadbalancing",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Application Load Balancers per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Application Load Balancers per Region"),
					Value:     ptr.Ptr(1.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-elasticloadbalancing",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"Application Load Balancers per Region: quota: 1, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service elasticloadbalancing and quota Application Load Balancers per Region"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient ALBs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "elasticloadbalancing",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Application Load Balancers per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Application Load Balancers per Region"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-elasticloadbalancing",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"Application Load Balancers per Region: quota: 5, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
		{
			name: "Fail (insufficient NLBs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "elasticloadbalancing",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Network Load Balancers per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Network Load Balancers per Region"),
					Value:     ptr.Ptr(1.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-elasticloadbalancing",
					Message:        "Usage for one or more service quotas exceeded the specified buffer",
					Details:        []string{"Network Load Balancers per Region: quota: 1, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       []string{"Remaining quota 0, less than buffer 3, for service elasticloadbalancing and quota Network Load Balancers per Region"},
					Status:         corev1.ConditionFalse,
				},
				State: ptr.Ptr(vapi.ValidationFailed),
			},
		},
		{
			name: "Pass (sufficient NLBs)",
			rule: v1alpha1.ServiceQuotaRule{
				Region:      "us-west-1",
				ServiceCode: "elasticloadbalancing",
				ServiceQuotas: []v1alpha1.ServiceQuota{
					{
						Name:   "Network Load Balancers per Region",
						Buffer: 3,
					},
				},
			},
			mockQuotas: []sqtypes.ServiceQuota{
				{
					QuotaName: ptr.Ptr("Network Load Balancers per Region"),
					Value:     ptr.Ptr(5.0),
				},
			},
			expectedResult: types.ValidationResult{
				Condition: &vapi.ValidationCondition{
					ValidationType: "aws-service-quota",
					ValidationRule: "validation-elasticloadbalancing",
					Message:        "Usage for all service quotas is below specified buffer",
					Details:        []string{"Network Load Balancers per Region: quota: 5, buffer: 3, max. usage: 1, max. usage entity: us-west-1"},
					Failures:       nil,
					Status:         corev1.ConditionTrue,
				},
				State: ptr.Ptr(vapi.ValidationSucceeded),
			},
		},
	}
	for _, c := range cs {
		fmt.Printf("Executing test: %s\n", c.name)
		mockQuotas.Quotas = c.mockQuotas
		result, err := svcQuotaService.ReconcileServiceQuotaRule(c.rule)
		test.CheckTestCase(t, result, c.expectedResult, err, c.expectedError)
	}
}
