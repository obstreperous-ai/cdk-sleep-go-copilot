package main

import (
	"os"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	_jsii_ "github.com/aws/jsii-runtime-go"
)

func TestMain(m *testing.M) {
	code := m.Run()
	_jsii_.Close()
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
