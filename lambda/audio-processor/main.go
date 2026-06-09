package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
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

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp  string                 `json:"timestamp"`
	Level      string                 `json:"level"`
	Message    string                 `json:"message"`
	RequestID  string                 `json:"requestId,omitempty"`
	Bucket     string                 `json:"bucket,omitempty"`
	Key        string                 `json:"key,omitempty"`
	Status     string                 `json:"status,omitempty"`
	Error      string                 `json:"error,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// logStructured outputs a structured JSON log entry
func logStructured(level, message string, data map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   message,
	}
	
	// Copy data fields to the log entry
	if requestID, ok := data["requestId"].(string); ok {
		entry.RequestID = requestID
	}
	if bucket, ok := data["bucket"].(string); ok {
		entry.Bucket = bucket
	}
	if key, ok := data["key"].(string); ok {
		entry.Key = key
	}
	if status, ok := data["status"].(string); ok {
		entry.Status = status
	}
	if err, ok := data["error"].(string); ok {
		entry.Error = err
	}
	if metadata, ok := data["metadata"].(map[string]interface{}); ok {
		entry.Metadata = metadata
	}
	
	jsonLog, _ := json.Marshal(entry)
	fmt.Println(string(jsonLog))
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
	// Extract request ID from Lambda context
	requestID := ""
	if lc, ok := lambdacontext.FromContext(ctx); ok {
		requestID = lc.AwsRequestID
	}
	
	// Log the input for debugging with structured logging
	logStructured("INFO", "Processing audio file request", map[string]interface{}{
		"requestId": requestID,
		"bucket":    input.Bucket,
		"key":       input.Key,
	})
	
	// Get environment variables
	tableName := os.Getenv("TABLE_NAME")
	if tableName == "" {
		logStructured("ERROR", "Environment variable not set", map[string]interface{}{
			"requestId": requestID,
			"error":     "TABLE_NAME environment variable not set",
		})
		return AudioProcessorOutput{}, fmt.Errorf("TABLE_NAME environment variable not set")
	}
	
	// Basic validation - check required fields
	if input.Bucket == "" || input.Key == "" {
		logStructured("ERROR", "Validation failed: missing required fields", map[string]interface{}{
			"requestId": requestID,
			"bucket":    input.Bucket,
			"key":       input.Key,
			"error":     "bucket and key are required",
		})
		return AudioProcessorOutput{}, fmt.Errorf("bucket and key are required")
	}
	
	// File extension validation
	ext := strings.ToLower(filepath.Ext(input.Key))
	if !supportedExtensions[ext] {
		logStructured("ERROR", "Validation failed: unsupported file extension", map[string]interface{}{
			"requestId": requestID,
			"bucket":    input.Bucket,
			"key":       input.Key,
			"extension": ext,
			"error":     fmt.Sprintf("unsupported file extension: %s. Supported formats: .mp3, .wav, .m4a, .flac, .ogg, .aac", ext),
		})
		return AudioProcessorOutput{}, fmt.Errorf("unsupported file extension: %s. Supported formats: .mp3, .wav, .m4a, .flac, .ogg, .aac", ext)
	}
	
	logStructured("INFO", "File extension validation passed", map[string]interface{}{
		"requestId": requestID,
		"extension": ext,
	})
	
	// Additional validation - check key doesn't contain path traversal
	if strings.Contains(input.Key, "..") {
		logStructured("ERROR", "Validation failed: path traversal detected", map[string]interface{}{
			"requestId": requestID,
			"bucket":    input.Bucket,
			"key":       input.Key,
			"error":     "invalid file key: path traversal detected",
		})
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
			"processor":   "SleepAudioProcessor",
			"version":     "1.0.0",
			"extension":   ext,
			"validatedAt": time.Now().UTC().Format(time.RFC3339),
		},
	}
	
	logStructured("INFO", "Audio file validation completed successfully", map[string]interface{}{
		"requestId": requestID,
		"bucket":    input.Bucket,
		"key":       input.Key,
		"status":    output.Status,
		"metadata":  output.Metadata,
	})
	
	return output, nil
}

func main() {
	lambda.Start(Handler)
}
