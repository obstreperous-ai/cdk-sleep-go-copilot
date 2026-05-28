package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

// TestCdkBaseStackSynthesizes verifies the base stack synthesizes without error.
func TestCdkBaseStackSynthesizes(t *testing.T) {
	defer jsii.Close()

	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewCdkBaseStack(app, "MyStack", nil)

	// THEN – obtaining the template is itself the synthesis assertion
	assertions.Template_FromStack(stack, nil)
}
