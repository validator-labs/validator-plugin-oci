// Package constants contains the constants used in validator-plugin-oci
package constants

const (
	// PluginCode is the constant for the plugin code
	PluginCode string = "OCI"
	// OciRegistry is the OCI registry string
	OciRegistry string = "oci-registry"
	// EcrRegistry is the ECR registry string
	EcrRegistry string = "ecr-registry"

	// AwsAccessKey is the key for the AWS access key
	AwsAccessKey = "AWS_ACCESS_KEY_ID" // #nosec
	// AwsSecretAccessKey is the key for the AWS secret access key
	AwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY" // #nosec
	// AwsSessionToken is the key for the AWS session token
	AwsSessionToken = "AWS_SESSION_TOKEN" // #nosec
)
