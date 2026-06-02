package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkBaseStackProps struct {
	awscdk.StackProps
}

func NewCdkBaseStack(scope constructs.Construct, id string, props *CdkBaseStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Input S3 Bucket - receives raw audio uploads
	inputBucket := awss3.NewBucket(stack, jsii.String("SleepAudioInputBucket"), &awss3.BucketProps{
		Encryption:        awss3.BucketEncryption_S3_MANAGED,
		Versioned:         jsii.Bool(true),
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		EventBridgeEnabled: jsii.Bool(true),
	})

	// Output S3 Bucket - stores processed audio
	awss3.NewBucket(stack, jsii.String("SleepAudioOutputBucket"), &awss3.BucketProps{
		Encryption:        awss3.BucketEncryption_S3_MANAGED,
		Versioned:         jsii.Bool(true),
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
	})

	// EventBridge Rule - triggers on Object Created events from Input Bucket
	awsevents.NewRule(stack, jsii.String("AudioUploadedRule"), &awsevents.RuleProps{
		EventPattern: &awsevents.EventPattern{
			Source:     jsii.Strings("aws.s3"),
			DetailType: jsii.Strings("Object Created"),
			Detail: &map[string]interface{}{
				"bucket": map[string]interface{}{
					"name": jsii.Strings(*inputBucket.BucketName()),
				},
			},
		},
		Enabled: jsii.Bool(true),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkBaseStack(app, "CdkBaseStack", &CdkBaseStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
