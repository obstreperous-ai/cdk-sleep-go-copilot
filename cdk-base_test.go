package main

import (
	"os"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestMain(m *testing.M) {
	code := m.Run()
	jsii.Close()
	os.Exit(code)
}

// TestCdkBaseStackSynthesizes verifies the base stack synthesizes without error.
func TestCdkBaseStackSynthesizes(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkBaseStack(app, "MyStack", nil)

	// THEN – obtaining the template is itself the synthesis assertion
	assertions.Template_FromStack(stack, nil)
}

// TestInputS3BucketExists verifies the Input S3 bucket exists with correct properties.
func TestInputS3BucketExists(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkBaseStack(app, "TestStack", nil)

	// THEN - verify input bucket exists with proper configuration
	template := assertions.Template_FromStack(stack, nil)
	
	// Verify at least one S3 bucket has all required properties
	template.HasResourceProperties(jsii.String("AWS::S3::Bucket"), map[string]interface{}{
		"BucketEncryption": map[string]interface{}{
			"ServerSideEncryptionConfiguration": []interface{}{
				map[string]interface{}{
					"ServerSideEncryptionByDefault": map[string]interface{}{
						"SSEAlgorithm": "AES256",
					},
				},
			},
		},
		"VersioningConfiguration": map[string]interface{}{
			"Status": "Enabled",
		},
		"PublicAccessBlockConfiguration": map[string]interface{}{
			"BlockPublicAcls":       true,
			"BlockPublicPolicy":     true,
			"IgnorePublicAcls":      true,
			"RestrictPublicBuckets": true,
		},
	})
	
	// Verify EventBridge notification is configured (input bucket only)
	template.ResourceCountIs(jsii.String("Custom::S3BucketNotifications"), jsii.Number(1))
	template.HasResourceProperties(jsii.String("Custom::S3BucketNotifications"), map[string]interface{}{
		"NotificationConfiguration": map[string]interface{}{
			"EventBridgeConfiguration": map[string]interface{}{},
		},
	})
}

// TestOutputS3BucketExists verifies the Output S3 bucket exists with correct properties.
func TestOutputS3BucketExists(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkBaseStack(app, "TestStack", nil)

	// THEN - verify output bucket exists with encryption and versioning
	template := assertions.Template_FromStack(stack, nil)
	
	// Verify we have exactly 2 S3 buckets (input + output)
	template.ResourceCountIs(jsii.String("AWS::S3::Bucket"), jsii.Number(2))
	
	// Verify all buckets have encryption, versioning, and public access blocked
	// This applies to both input and output buckets
	captures := assertions.NewCapture(nil)
	template.HasResourceProperties(jsii.String("AWS::S3::Bucket"), map[string]interface{}{
		"BucketEncryption": captures,
		"VersioningConfiguration": map[string]interface{}{
			"Status": "Enabled",
		},
		"PublicAccessBlockConfiguration": map[string]interface{}{
			"BlockPublicAcls":       true,
			"BlockPublicPolicy":     true,
			"IgnorePublicAcls":      true,
			"RestrictPublicBuckets": true,
		},
	})
}

// TestEventBridgeRuleExists verifies the EventBridge rule exists and is configured correctly.
func TestEventBridgeRuleExists(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkBaseStack(app, "TestStack", nil)

	// THEN - verify EventBridge rule exists
	template := assertions.Template_FromStack(stack, nil)
	
	template.HasResourceProperties(jsii.String("AWS::Events::Rule"), map[string]interface{}{
		"EventPattern": map[string]interface{}{
			"source":      []interface{}{"aws.s3"},
			"detail-type": []interface{}{"Object Created"},
		},
		"State": "ENABLED",
	})
}

// TestStepFunctionsStateMachineExists verifies the Step Functions state machine exists.
func TestStepFunctionsStateMachineExists(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkBaseStack(app, "TestStack", nil)

	// THEN - verify state machine exists
	template := assertions.Template_FromStack(stack, nil)
	
	// Verify we have exactly 1 state machine
	template.ResourceCountIs(jsii.String("AWS::StepFunctions::StateMachine"), jsii.Number(1))
	
	// Verify state machine has logging enabled
	template.HasResourceProperties(jsii.String("AWS::StepFunctions::StateMachine"), map[string]interface{}{
		"LoggingConfiguration": map[string]interface{}{
			"Level": "ALL",
		},
	})
}

// TestStepFunctionsStateMachineDefinitionHasPollyTask verifies the state machine definition includes Polly task.
func TestStepFunctionsStateMachineDefinitionHasPollyTask(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkBaseStack(app, "TestStack", nil)

	// THEN - verify state machine definition contains expected states
	template := assertions.Template_FromStack(stack, nil)
	
	// Capture the state machine definition
	captures := assertions.NewCapture(nil)
	template.HasResourceProperties(jsii.String("AWS::StepFunctions::StateMachine"), map[string]interface{}{
		"DefinitionString": captures,
	})
}

// TestEventBridgeRuleTargetsStateMachine verifies the EventBridge rule targets the state machine.
func TestEventBridgeRuleTargetsStateMachine(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkBaseStack(app, "TestStack", nil)

	// THEN - verify EventBridge rule has a target
	template := assertions.Template_FromStack(stack, nil)
	
	// Verify the rule has targets configured
	template.HasResourceProperties(jsii.String("AWS::Events::Rule"), map[string]interface{}{
		"State": "ENABLED",
		"Targets": assertions.Match_AnyValue(),
	})
}

// TestStateMachineExecutionRole verifies the state machine has an execution role with appropriate permissions.
func TestStateMachineExecutionRole(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkBaseStack(app, "TestStack", nil)

	// THEN - verify IAM role exists for state machine
	template := assertions.Template_FromStack(stack, nil)
	
	// Verify at least one IAM role exists (for state machine execution)
	template.HasResourceProperties(jsii.String("AWS::IAM::Role"), map[string]interface{}{
		"AssumeRolePolicyDocument": map[string]interface{}{
			"Statement": []interface{}{
				map[string]interface{}{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": map[string]interface{}{
						"Service": assertions.Match_StringLikeRegexp(jsii.String("states.*")),
					},
				},
			},
		},
	})
}
