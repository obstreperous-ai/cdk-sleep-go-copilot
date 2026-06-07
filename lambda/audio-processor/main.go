package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
)

// AudioProcessorInput represents the input from Step Functions
type AudioProcessorInput struct {
	Bucket       string `json:"bucket"`
	Key          string `json:"key"`
	DynamoResult *struct {
		SdkResponseMetadata map[string]interface{} `json:"SdkResponseMetadata,omitempty"`
	} `json:"dynamoResult,omitempty"`
}

// AudioProcessorOutput represents the output returned to Step Functions
type AudioProcessorOutput struct {
	Bucket    string                 `json:"bucket"`
	Key       string                 `json:"key"`
	AudioID   string                 `json:"audioId"`
	Status    string                 `json:"status"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Supported audio file extensions
var supportedExtensions = map[string]bool{
	".mp3":  true,
	".wav":  true,
	".m4a":  true,
	".flac": true,
	".ogg":  true,
	".aac":  true,
}

// Handler is the Lambda function handler
func Handler(ctx context.Context, input AudioProcessorInput) (AudioProcessorOutput, error) {
	// Log the input for debugging
	inputJSON, _ := json.Marshal(input)
	log.Printf("Received input: %s", string(inputJSON))
	
	// Get environment variables
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		return AudioProcessorOutput{}, fmt.Errorf("TABLE_NAME environment variable not set")
	}
	
	log.Printf("Processing audio file: bucket=%s, key=%s, table=%s", input.Bucket, input.Key, tableName)
	
	// Basic validation - check required fields
	if input.Bucket == "" || input.Key == "" {
		return AudioProcessorOutput{}, fmt.Errorf("bucket and key are required")
	}
	
	// File extension validation
	ext := strings.ToLower(filepath.Ext(input.Key))
	if !supportedExtensions[ext] {
		return AudioProcessorOutput{}, fmt.Errorf("unsupported file extension: %s. Supported formats: .mp3, .wav, .m4a, .flac, .ogg, .aac", ext)
	}
	
	log.Printf("File extension validation passed: %s", ext)
	
	// Additional validation - check key doesn't contain path traversal
	if strings.Contains(input.Key, "..") {
		return AudioProcessorOutput{}, fmt.Errorf("invalid file key: path traversal detected")
	}
	
	// Placeholder: In the future, this will validate audio format, extract metadata, etc.
	// For now, just return a success response
	
	output := AudioProcessorOutput{
		Bucket:  input.Bucket,
		Key:     input.Key,
		AudioID: input.Key, // Using key as audioID for now
		Status:  "validated",
		Message: "Audio file received and validation passed",
		Metadata: map[string]interface{}{
			"processor":  "SleepAudioProcessor",
			"version":    "1.0.0",
			"extension":  ext,
			"validatedAt": "now", // In production, use actual timestamp
		},
	}
	
	outputJSON, _ := json.Marshal(output)
	log.Printf("Returning output: %s", string(outputJSON))
	
	return output, nil
}

func main() {
	lambda.Start(Handler)
}
