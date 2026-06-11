package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3"
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
	Bucket       string                 `json:"bucket"`
	Key          string                 `json:"key"`
	AudioID      string                 `json:"audioId"`
	OutputBucket string                 `json:"outputBucket,omitempty"`
	OutputKey    string                 `json:"outputKey,omitempty"`
	Status       string                 `json:"status"`
	Message      string                 `json:"message"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessingResult holds the result of audio processing
type ProcessingResult struct {
	OutputBucket string
	OutputKey    string
	FileSize     int64
	Duration     float64
	Format       string
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

// Constants for audio processing configuration
const (
	processorVersion = "2.0.0"
	pollyCharLimit   = 3000
	mp3BitrateKbps   = 128
	wavBitrateKbps   = 1411
)

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
	
	jsonLog, err := json.Marshal(entry)
	if err != nil {
		// Fallback to plain text if JSON marshaling fails
		fmt.Fprintf(os.Stderr, "[%s] %s: %s (marshal error: %v)\n", 
			time.Now().UTC().Format(time.RFC3339), level, message, err)
		return
	}
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
	outputBucket := os.Getenv("OUTPUT_BUCKET")
	
	if tableName == "" {
		logStructured("ERROR", "Environment variable not set", map[string]interface{}{
			"requestId": requestID,
			"error":     "TABLE_NAME environment variable not set",
		})
		return AudioProcessorOutput{}, fmt.Errorf("TABLE_NAME environment variable not set")
	}
	
	if outputBucket == "" {
		logStructured("ERROR", "Environment variable not set", map[string]interface{}{
			"requestId": requestID,
			"error":     "OUTPUT_BUCKET environment variable not set",
		})
		return AudioProcessorOutput{}, fmt.Errorf("OUTPUT_BUCKET environment variable not set")
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
	
	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		logStructured("ERROR", "Failed to create AWS session", map[string]interface{}{
			"requestId": requestID,
			"error":     err.Error(),
		})
		return AudioProcessorOutput{}, fmt.Errorf("failed to create AWS session: %w", err)
	}
	
	// Process the audio file
	result, err := processAudioFile(ctx, sess, input.Bucket, input.Key, outputBucket, requestID)
	if err != nil {
		logStructured("ERROR", "Failed to process audio file", map[string]interface{}{
			"requestId": requestID,
			"error":     err.Error(),
		})
		return AudioProcessorOutput{}, fmt.Errorf("failed to process audio file: %w", err)
	}
	
	// Update DynamoDB with processing results
	err = updateMetadata(sess, tableName, input.Key, result, requestID)
	if err != nil {
		logStructured("ERROR", "Failed to update DynamoDB metadata", map[string]interface{}{
			"requestId": requestID,
			"error":     err.Error(),
		})
		// Log error but don't fail the function - the file was already processed
	}
	
	// Prepare successful output
	output := AudioProcessorOutput{
		Bucket:       input.Bucket,
		Key:          input.Key,
		AudioID:      input.Key,
		OutputBucket: result.OutputBucket,
		OutputKey:    result.OutputKey,
		Status:       "processed",
		Message:      "Audio file processed successfully",
		Metadata: map[string]interface{}{
			"processor":    "SleepAudioProcessor",
			"version":      processorVersion,
			"extension":    ext,
			"processedAt":  time.Now().UTC().Format(time.RFC3339),
			"outputBucket": result.OutputBucket,
			"outputKey":    result.OutputKey,
			"fileSize":     result.FileSize,
			"duration":     result.Duration,
			"format":       result.Format,
		},
	}
	
	logStructured("INFO", "Audio file processing completed successfully", map[string]interface{}{
		"requestId":    requestID,
		"bucket":       input.Bucket,
		"key":          input.Key,
		"status":       output.Status,
		"outputBucket": result.OutputBucket,
		"outputKey":    result.OutputKey,
		"metadata":     output.Metadata,
	})
	
	return output, nil
}

// processAudioFile handles the core audio processing logic
func processAudioFile(ctx context.Context, sess *session.Session, inputBucket, inputKey, outputBucket, requestID string) (*ProcessingResult, error) {
	s3Client := s3.New(sess)
	pollyClient := polly.New(sess)
	
	// Download the input file from S3
	logStructured("INFO", "Downloading input file from S3", map[string]interface{}{
		"requestId": requestID,
		"bucket":    inputBucket,
		"key":       inputKey,
	})
	
	getObjectOutput, err := s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(inputBucket),
		Key:    aws.String(inputKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download from S3: %w", err)
	}
	defer getObjectOutput.Body.Close()
	
	// Read the file content
	inputData, err := io.ReadAll(getObjectOutput.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}
	
	inputSize := len(inputData)
	
	logStructured("INFO", "Downloaded input file", map[string]interface{}{
		"requestId": requestID,
		"size":      inputSize,
	})
	
	// Process the audio - for this implementation, we'll:
	// 1. For audio files: Pass through (or apply basic normalization placeholder)
	// 2. For text files: Use Polly to synthesize speech (not implemented in this version but structure is ready)
	
	var processedData []byte
	var format string
	
	ext := strings.ToLower(filepath.Ext(inputKey))
	
	// Check if this is a text file that needs Polly synthesis
	if ext == ".txt" {
		// Use Polly to synthesize soothing speech
		logStructured("INFO", "Synthesizing speech with Polly", map[string]interface{}{
			"requestId": requestID,
		})
		
		// Read text content
		textContent := string(inputData)
		if len(textContent) > pollyCharLimit {
			textContent = textContent[:pollyCharLimit] // Polly limit
		}
		
		synthesizeOutput, err := pollyClient.SynthesizeSpeechWithContext(ctx, &polly.SynthesizeSpeechInput{
			Text:         aws.String(textContent),
			OutputFormat: aws.String("mp3"),
			VoiceId:      aws.String("Joanna"),
			Engine:       aws.String("neural"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to synthesize speech with Polly: %w", err)
		}
		defer synthesizeOutput.AudioStream.Close()
		
		processedData, err = io.ReadAll(synthesizeOutput.AudioStream)
		if err != nil {
			return nil, fmt.Errorf("failed to read Polly output: %w", err)
		}
		format = "mp3"
	} else {
		// For audio files, pass through (in a real implementation, you might apply normalization, etc.)
		processedData = inputData
		format = strings.TrimPrefix(ext, ".")
		
		logStructured("INFO", "Processing audio file (passthrough mode)", map[string]interface{}{
			"requestId": requestID,
			"format":    format,
		})
	}
	
	// Generate output key with naming convention: processed-<timestamp>-<original-filename>
	timestamp := time.Now().UTC().Format("20060102-150405")
	filename := filepath.Base(inputKey)
	filenameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	outputKey := fmt.Sprintf("processed-%s-%s.%s", timestamp, filenameWithoutExt, format)
	
	// Upload processed file to output bucket
	logStructured("INFO", "Uploading processed file to output bucket", map[string]interface{}{
		"requestId":    requestID,
		"outputBucket": outputBucket,
		"outputKey":    outputKey,
		"size":         len(processedData),
	})
	
	_, err = s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(outputBucket),
		Key:         aws.String(outputKey),
		Body:        bytes.NewReader(processedData),
		ContentType: aws.String(fmt.Sprintf("audio/%s", format)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to output bucket: %w", err)
	}
	
	// Calculate approximate duration (simplified - real implementation would parse audio metadata)
	// For demonstration: assume ~128kbps for MP3, ~1411kbps for WAV
	var durationSeconds float64
	if format == "mp3" {
		durationSeconds = float64(len(processedData)) / (mp3BitrateKbps * 1024 / 8) // bytes / (bitrate in bytes/sec)
	} else if format == "wav" {
		durationSeconds = float64(len(processedData)) / (wavBitrateKbps * 1024 / 8)
	} else {
		durationSeconds = float64(len(processedData)) / (mp3BitrateKbps * 1024 / 8) // default estimate
	}
	
	result := &ProcessingResult{
		OutputBucket: outputBucket,
		OutputKey:    outputKey,
		FileSize:     int64(len(processedData)),
		Duration:     durationSeconds,
		Format:       format,
	}
	
	logStructured("INFO", "File processing completed", map[string]interface{}{
		"requestId":    requestID,
		"outputBucket": outputBucket,
		"outputKey":    outputKey,
		"fileSize":     result.FileSize,
		"duration":     result.Duration,
	})
	
	return result, nil
}

// updateMetadata updates DynamoDB with processing results
func updateMetadata(sess *session.Session, tableName, audioID string, result *ProcessingResult, requestID string) error {
	ddbClient := dynamodb.New(sess)
	
	logStructured("INFO", "Updating DynamoDB metadata", map[string]interface{}{
		"requestId": requestID,
		"tableName": tableName,
		"audioId":   audioID,
	})
	
	_, err := ddbClient.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"audioId": {
				S: aws.String(audioID),
			},
		},
		UpdateExpression: aws.String("SET outputBucket = :outputBucket, outputKey = :outputKey, fileSize = :fileSize, #dur = :duration, #fmt = :format, updatedAt = :updatedAt"),
		ExpressionAttributeNames: map[string]*string{
			"#dur": aws.String("duration"),
			"#fmt": aws.String("format"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":outputBucket": {
				S: aws.String(result.OutputBucket),
			},
			":outputKey": {
				S: aws.String(result.OutputKey),
			},
			":fileSize": {
				N: aws.String(fmt.Sprintf("%d", result.FileSize)),
			},
			":duration": {
				N: aws.String(fmt.Sprintf("%.2f", result.Duration)),
			},
			":format": {
				S: aws.String(result.Format),
			},
			":updatedAt": {
				S: aws.String(time.Now().UTC().Format(time.RFC3339)),
			},
		},
	})
	
	if err != nil {
		return fmt.Errorf("failed to update DynamoDB: %w", err)
	}
	
	logStructured("INFO", "DynamoDB metadata updated successfully", map[string]interface{}{
		"requestId": requestID,
		"audioId":   audioID,
	})
	
	return nil
}

func main() {
	lambda.Start(Handler)
}
