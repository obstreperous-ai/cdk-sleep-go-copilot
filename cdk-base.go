package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssns"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsstepfunctionstasks"
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

	// Determine environment from CDK context (dev, stage, prod)
	// Default to dev if not specified
	envContext := "dev"
	if stack.Node().TryGetContext(jsii.String("env")) != nil {
		if envStr, ok := stack.Node().TryGetContext(jsii.String("env")).(string); ok {
			envContext = envStr
		}
	}
	
	// Set removal policy based on environment
	// dev: DESTROY for easy cleanup
	// stage/prod: RETAIN for data safety
	var removalPolicy awscdk.RemovalPolicy
	if envContext == "dev" {
		removalPolicy = awscdk.RemovalPolicy_DESTROY
	} else {
		removalPolicy = awscdk.RemovalPolicy_RETAIN
	}

	// Input S3 Bucket - receives raw audio uploads
	inputBucket := awss3.NewBucket(stack, jsii.String("SleepAudioInputBucket"), &awss3.BucketProps{
		Encryption:        awss3.BucketEncryption_S3_MANAGED,
		Versioned:         jsii.Bool(true),
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		EventBridgeEnabled: jsii.Bool(true),
		RemovalPolicy:     removalPolicy,
		AutoDeleteObjects: jsii.Bool(envContext == "dev"), // Only auto-delete in dev
	})

	// Output S3 Bucket - stores processed audio
	awss3.NewBucket(stack, jsii.String("SleepAudioOutputBucket"), &awss3.BucketProps{
		Encryption:        awss3.BucketEncryption_S3_MANAGED,
		Versioned:         jsii.Bool(true),
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		RemovalPolicy:     removalPolicy,
		AutoDeleteObjects: jsii.Bool(envContext == "dev"), // Only auto-delete in dev
	})

	// DynamoDB Table - stores audio pipeline metadata
	metadataTable := awsdynamodb.NewTable(stack, jsii.String("SleepAudioMetadataTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("audioId"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode:              awsdynamodb.BillingMode_PAY_PER_REQUEST,
		Encryption:               awsdynamodb.TableEncryption_AWS_MANAGED,
		PointInTimeRecoverySpecification: &awsdynamodb.PointInTimeRecoverySpecification{
			PointInTimeRecoveryEnabled: jsii.Bool(true),
		},
		RemovalPolicy:            removalPolicy,
	})

	// Lambda Function - Audio Processor
	// This Lambda validates audio files and extracts metadata before processing
	audioProcessorFunction := awslambda.NewFunction(stack, jsii.String("SleepAudioProcessor"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Handler: jsii.String("bootstrap"),
		Code:    awslambda.Code_FromAsset(jsii.String("lambda/audio-processor"), nil),
		Environment: &map[string]*string{
			"TABLE_NAME": metadataTable.TableName(),
		},
		Timeout: awscdk.Duration_Minutes(jsii.Number(1)),
		MemorySize: jsii.Number(256),
		Description: jsii.String("Validates audio files and extracts metadata for the sleep audio pipeline"),
	})

	// Grant the Lambda function read access to DynamoDB
	metadataTable.GrantReadData(audioProcessorFunction)

	// SNS Topics - for pipeline notifications

	// KMS Key for SNS topic encryption
	snsKmsKey := awskms.NewKey(stack, jsii.String("SnsTopicKey"), &awskms.KeyProps{
		Description:      jsii.String("KMS key for SNS topic encryption"),
		EnableKeyRotation: jsii.Bool(true),
		RemovalPolicy:    awscdk.RemovalPolicy_DESTROY, // For dev/test - change for production
	})

	// SNS Topic - pipeline completion notifications
	completedTopic := awssns.NewTopic(stack, jsii.String("SleepAudioPipelineCompletedTopic"), &awssns.TopicProps{
		DisplayName:    jsii.String("Sleep Audio Pipeline Completed"),
		MasterKey:      snsKmsKey,
	})

	// SNS Topic - pipeline failure notifications
	failedTopic := awssns.NewTopic(stack, jsii.String("SleepAudioPipelineFailedTopic"), &awssns.TopicProps{
		DisplayName:    jsii.String("Sleep Audio Pipeline Failed"),
		MasterKey:      snsKmsKey,
	})

	// CloudWatch Log Group for Step Functions state machine
	logGroup := awslogs.NewLogGroup(stack, jsii.String("StateMachineLogGroup"), &awslogs.LogGroupProps{
		Retention:         awslogs.RetentionDays_ONE_WEEK,
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
	})

	// Step Functions State Machine - Audio Processing Pipeline
	
	// Pass state to set error message for missing input validation
	invalidInputError := awsstepfunctions.NewPass(stack, jsii.String("InvalidInputError"), &awsstepfunctions.PassProps{
		Parameters: &map[string]interface{}{
			"errorMessage": jsii.String("Invalid input: bucket and key are required"),
			"bucket":       awsstepfunctions.JsonPath_StringAt(jsii.String("$.bucket")),
			"key":          awsstepfunctions.JsonPath_StringAt(jsii.String("$.key")),
		},
		ResultPath: jsii.String("$"),
	})
	
	// Choice state for input validation - check if bucket and key exist
	validateInputChoice := awsstepfunctions.NewChoice(stack, jsii.String("ValidateInput"), &awsstepfunctions.ChoiceProps{
		Comment: jsii.String("Validate that bucket and key are present in the input"),
	})
	
	// DynamoDB PutItem task - write initial metadata record when pipeline starts
	putItemTask := awsstepfunctionstasks.NewDynamoPutItem(stack, jsii.String("WriteInitialMetadata"), &awsstepfunctionstasks.DynamoPutItemProps{
		Table: metadataTable,
		Item: &map[string]awsstepfunctionstasks.DynamoAttributeValue{
			"audioId": awsstepfunctionstasks.DynamoAttributeValue_FromString(
				awsstepfunctions.JsonPath_StringAt(jsii.String("$.key")),
			),
			"inputBucket": awsstepfunctionstasks.DynamoAttributeValue_FromString(
				awsstepfunctions.JsonPath_StringAt(jsii.String("$.bucket")),
			),
			"inputKey": awsstepfunctionstasks.DynamoAttributeValue_FromString(
				awsstepfunctions.JsonPath_StringAt(jsii.String("$.key")),
			),
			"status": awsstepfunctionstasks.DynamoAttributeValue_FromString(jsii.String("PROCESSING")),
			"createdAt": awsstepfunctionstasks.DynamoAttributeValue_FromString(
				awsstepfunctions.JsonPath_StringAt(jsii.String("$$.State.EnteredTime")),
			),
		},
		ResultPath: jsii.String("$.dynamoResult"),
	})
	
	// Lambda invocation task - process and validate audio file
	processAudioTask := awsstepfunctionstasks.NewLambdaInvoke(stack, jsii.String("ProcessAudioFile"), &awsstepfunctionstasks.LambdaInvokeProps{
		LambdaFunction: audioProcessorFunction,
		ResultPath:     jsii.String("$.processorResult"),
		PayloadResponseOnly: jsii.Bool(true),
	})
	
	// Pass state to format Lambda errors for the error handling chain
	formatLambdaError := awsstepfunctions.NewPass(stack, jsii.String("FormatLambdaError"), &awsstepfunctions.PassProps{
		Parameters: &map[string]interface{}{
			"errorMessage": awsstepfunctions.JsonPath_StringAt(jsii.String("$.errorInfo.Error")),
		},
		ResultPath: jsii.String("$"),
	})
	
	// Polly task - synthesizes speech from text (placeholder configuration for now)
	pollyTask := awsstepfunctionstasks.NewCallAwsService(stack, jsii.String("PollyTask"), &awsstepfunctionstasks.CallAwsServiceProps{
		Service: jsii.String("polly"),
		Action:  jsii.String("synthesizeSpeech"),
		Parameters: &map[string]interface{}{
			"Text":         jsii.String("This is a placeholder text for sleep audio narration."),
			"VoiceId":      jsii.String("Joanna"),
			"OutputFormat": jsii.String("mp3"),
		},
		IamResources: jsii.Strings("*"),
		ResultPath:   jsii.String("$.pollyResult"),
	})

	// DynamoDB UpdateItem task - update status to COMPLETED on success
	updateStatusCompleted := awsstepfunctionstasks.NewDynamoUpdateItem(stack, jsii.String("UpdateStatusCompleted"), &awsstepfunctionstasks.DynamoUpdateItemProps{
		Table: metadataTable,
		Key: &map[string]awsstepfunctionstasks.DynamoAttributeValue{
			"audioId": awsstepfunctionstasks.DynamoAttributeValue_FromString(
				awsstepfunctions.JsonPath_StringAt(jsii.String("$.key")),
			),
		},
		UpdateExpression: jsii.String("SET #status = :completed, #updatedAt = :timestamp"),
		ExpressionAttributeNames: &map[string]*string{
			"#status":    jsii.String("status"),
			"#updatedAt": jsii.String("updatedAt"),
		},
		ExpressionAttributeValues: &map[string]awsstepfunctionstasks.DynamoAttributeValue{
			":completed": awsstepfunctionstasks.DynamoAttributeValue_FromString(jsii.String("COMPLETED")),
			":timestamp": awsstepfunctionstasks.DynamoAttributeValue_FromString(
				awsstepfunctions.JsonPath_StringAt(jsii.String("$$.State.EnteredTime")),
			),
		},
		ResultPath: jsii.String("$.updateResult"),
	})

	// SNS Publish task - send completion notification
	publishCompletedNotification := awsstepfunctionstasks.NewSnsPublish(stack, jsii.String("PublishCompletedNotification"), &awsstepfunctionstasks.SnsPublishProps{
		Topic:   completedTopic,
		Message: awsstepfunctions.TaskInput_FromObject(&map[string]interface{}{
			"status": "COMPLETED",
			"audioId": awsstepfunctions.JsonPath_StringAt(jsii.String("$.key")),
			"bucket":  awsstepfunctions.JsonPath_StringAt(jsii.String("$.bucket")),
			"message": "Audio processing pipeline completed successfully",
		}),
		Subject:    jsii.String("Audio Pipeline Completed"),
		ResultPath: jsii.String("$.notificationResult"),
	})

	// DynamoDB UpdateItem task - update status to FAILED on error
	updateStatusFailed := awsstepfunctionstasks.NewDynamoUpdateItem(stack, jsii.String("UpdateStatusFailed"), &awsstepfunctionstasks.DynamoUpdateItemProps{
		Table: metadataTable,
		Key: &map[string]awsstepfunctionstasks.DynamoAttributeValue{
			"audioId": awsstepfunctionstasks.DynamoAttributeValue_FromString(
				awsstepfunctions.JsonPath_StringAt(jsii.String("$.key")),
			),
		},
		UpdateExpression: jsii.String("SET #status = :failed, #updatedAt = :timestamp, #error = :errorMsg"),
		ExpressionAttributeNames: &map[string]*string{
			"#status":    jsii.String("status"),
			"#updatedAt": jsii.String("updatedAt"),
			"#error":     jsii.String("errorMessage"),
		},
		ExpressionAttributeValues: &map[string]awsstepfunctionstasks.DynamoAttributeValue{
			":failed":    awsstepfunctionstasks.DynamoAttributeValue_FromString(jsii.String("FAILED")),
			":timestamp": awsstepfunctionstasks.DynamoAttributeValue_FromString(
				awsstepfunctions.JsonPath_StringAt(jsii.String("$$.State.EnteredTime")),
			),
			":errorMsg": awsstepfunctionstasks.DynamoAttributeValue_FromString(
				awsstepfunctions.JsonPath_StringAt(jsii.String("$.errorMessage")),
			),
		},
		ResultPath: jsii.String("$.updateResult"),
	})

	// SNS Publish task - send failure notification
	publishFailedNotification := awsstepfunctionstasks.NewSnsPublish(stack, jsii.String("PublishFailedNotification"), &awsstepfunctionstasks.SnsPublishProps{
		Topic:   failedTopic,
		Message: awsstepfunctions.TaskInput_FromObject(&map[string]interface{}{
			"status": "FAILED",
			"audioId": awsstepfunctions.JsonPath_StringAt(jsii.String("$.key")),
			"bucket":  awsstepfunctions.JsonPath_StringAt(jsii.String("$.bucket")),
			"error":   awsstepfunctions.JsonPath_StringAt(jsii.String("$.errorMessage")),
			"message": "Audio processing pipeline failed",
		}),
		Subject:    jsii.String("Audio Pipeline Failed"),
		ResultPath: jsii.String("$.notificationResult"),
	})

	// Define error handling path: format error, update status to FAILED, then send notification
	lambdaErrorChain := formatLambdaError.Next(updateStatusFailed).Next(publishFailedNotification)
	
	// Wire up invalid input error path
	invalidInputError.Next(publishFailedNotification)
	
	// Add error handling to the Lambda task (for validation errors)
	processAudioTask.AddCatch(lambdaErrorChain, &awsstepfunctions.CatchProps{
		Errors: jsii.Strings("States.ALL"),
		ResultPath: jsii.String("$.errorInfo"),
	})
	
	// Add error handling to the Polly task
	pollyTask.AddCatch(lambdaErrorChain, &awsstepfunctions.CatchProps{
		Errors: jsii.Strings("States.ALL"),
		ResultPath: jsii.String("$.errorInfo"),
	}).Next(updateStatusCompleted).Next(publishCompletedNotification)
	
	// Wire up the Choice state for input validation
	// Check if both bucket and key are present (not null and have a length > 0)
	validateInputChoice.When(
		awsstepfunctions.Condition_And(
			awsstepfunctions.Condition_IsPresent(jsii.String("$.bucket")),
			awsstepfunctions.Condition_IsPresent(jsii.String("$.key")),
		),
		putItemTask.Next(processAudioTask).Next(pollyTask),
		nil,
	).Otherwise(invalidInputError)

	// Define the state machine starting with input validation
	stateMachineDefinition := validateInputChoice

	// Create the state machine
	stateMachine := awsstepfunctions.NewStateMachine(stack, jsii.String("AudioProcessingStateMachine"), &awsstepfunctions.StateMachineProps{
		DefinitionBody: awsstepfunctions.DefinitionBody_FromChainable(stateMachineDefinition),
		Logs: &awsstepfunctions.LogOptions{
			Destination: logGroup,
			Level:       awsstepfunctions.LogLevel_ALL,
		},
		TracingEnabled: jsii.Bool(true),
	})

	// EventBridge Rule - triggers on Object Created events from Input Bucket
	rule := awsevents.NewRule(stack, jsii.String("AudioUploadedRule"), &awsevents.RuleProps{
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

	// Add the state machine as a target for the EventBridge rule
	rule.AddTarget(awseventstargets.NewSfnStateMachine(stateMachine, &awseventstargets.SfnStateMachineProps{
		Input: awsevents.RuleTargetInput_FromObject(&map[string]interface{}{
			"bucket": awsevents.EventField_FromPath(jsii.String("$.detail.bucket.name")),
			"key":    awsevents.EventField_FromPath(jsii.String("$.detail.object.key")),
		}),
	}))

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
