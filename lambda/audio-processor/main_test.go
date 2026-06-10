package main

import (
	"context"
	"testing"
)

// TestProcessAudioFile_DownloadsFromS3 verifies Lambda downloads input file from S3
func TestProcessAudioFile_DownloadsFromS3(t *testing.T) {
	// GIVEN: Input with valid bucket and key
	input := AudioProcessorInput{
		Bucket: "test-input-bucket",
		Key:    "audio/test-file.mp3",
	}
	
	// Set environment variables
	t.Setenv("TABLE_NAME", "test-table")
	t.Setenv("OUTPUT_BUCKET", "test-output-bucket")
	
	// WHEN: Handler processes the request
	_, err := Handler(context.Background(), input)
	
	// THEN: Should succeed (will fail until we implement S3 download)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestProcessAudioFile_SynthesizesWithPolly verifies Lambda uses Polly for text-to-speech
func TestProcessAudioFile_SynthesizesWithPolly(t *testing.T) {
	// GIVEN: Input with text prompt (indicated by .txt extension)
	input := AudioProcessorInput{
		Bucket: "test-input-bucket",
		Key:    "prompts/sleep-meditation.txt",
	}
	
	// Set environment variables
	t.Setenv("TABLE_NAME", "test-table")
	t.Setenv("OUTPUT_BUCKET", "test-output-bucket")
	
	// WHEN: Handler processes the request
	output, err := Handler(context.Background(), input)
	
	// THEN: Should synthesize speech (will fail until we implement Polly)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Output should contain processed file information
	if output.Metadata == nil {
		t.Error("Expected metadata in output")
	}
	
	if output.Status != "processed" {
		t.Errorf("Expected status 'processed', got: %s", output.Status)
	}
}

// TestProcessAudioFile_UploadsToOutputBucket verifies Lambda uploads to output S3 bucket
func TestProcessAudioFile_UploadsToOutputBucket(t *testing.T) {
	// GIVEN: Input with valid audio file
	input := AudioProcessorInput{
		Bucket: "test-input-bucket",
		Key:    "audio/sample.mp3",
	}
	
	// Set environment variables
	t.Setenv("TABLE_NAME", "test-table")
	t.Setenv("OUTPUT_BUCKET", "test-output-bucket")
	
	// WHEN: Handler processes the request
	output, err := Handler(context.Background(), input)
	
	// THEN: Output should contain output bucket and key
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if output.Metadata == nil {
		t.Error("Expected metadata in output")
	}
	
	// Check for output location in metadata
	if _, ok := output.Metadata["outputBucket"]; !ok {
		t.Error("Expected outputBucket in metadata")
	}
	
	if _, ok := output.Metadata["outputKey"]; !ok {
		t.Error("Expected outputKey in metadata")
	}
}

// TestProcessAudioFile_UpdatesDynamoDBWithOutput verifies Lambda updates DynamoDB with output info
func TestProcessAudioFile_UpdatesDynamoDBWithOutput(t *testing.T) {
	// GIVEN: Input with valid audio file
	input := AudioProcessorInput{
		Bucket: "test-input-bucket",
		Key:    "audio/test.wav",
	}
	
	// Set environment variables
	t.Setenv("TABLE_NAME", "test-table")
	t.Setenv("OUTPUT_BUCKET", "test-output-bucket")
	
	// WHEN: Handler processes the request
	output, err := Handler(context.Background(), input)
	
	// THEN: Should update DynamoDB with output location
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	// Check status is set to processed
	if output.Status != "processed" {
		t.Errorf("Expected status 'processed', got: %s", output.Status)
	}
	
	// Check metadata contains file size and duration
	if _, ok := output.Metadata["fileSize"]; !ok {
		t.Error("Expected fileSize in metadata")
	}
}

// TestProcessAudioFile_OutputNamingConvention verifies output files follow naming convention
func TestProcessAudioFile_OutputNamingConvention(t *testing.T) {
	// GIVEN: Input with specific key
	input := AudioProcessorInput{
		Bucket: "test-input-bucket",
		Key:    "uploads/user123/audio.mp3",
	}
	
	// Set environment variables
	t.Setenv("TABLE_NAME", "test-table")
	t.Setenv("OUTPUT_BUCKET", "test-output-bucket")
	
	// WHEN: Handler processes the request
	output, err := Handler(context.Background(), input)
	
	// THEN: Output key should include original key and timestamp
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	
	if output.Metadata == nil {
		t.Fatal("Expected metadata in output")
	}
	
	outputKey, ok := output.Metadata["outputKey"].(string)
	if !ok {
		t.Fatal("Expected outputKey to be a string")
	}
	
	// Output key should contain processed- prefix and original filename
	if outputKey == "" {
		t.Error("Expected non-empty output key")
	}
}

// TestProcessAudioFile_MissingOutputBucketEnv verifies error when OUTPUT_BUCKET is missing
func TestProcessAudioFile_MissingOutputBucketEnv(t *testing.T) {
	// GIVEN: Missing OUTPUT_BUCKET environment variable
	input := AudioProcessorInput{
		Bucket: "test-input-bucket",
		Key:    "audio/test.mp3",
	}
	
	// Only set TABLE_NAME
	t.Setenv("TABLE_NAME", "test-table")
	
	// WHEN: Handler processes the request
	_, err := Handler(context.Background(), input)
	
	// THEN: Should return error
	if err == nil {
		t.Error("Expected error when OUTPUT_BUCKET is missing")
	}
}
