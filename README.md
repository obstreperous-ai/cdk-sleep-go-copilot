# cdk-sleep-go-copilot

[![CI](https://github.com/obstreperous-ai/cdk-sleep-go-copilot/actions/workflows/ci.yml/badge.svg)](https://github.com/obstreperous-ai/cdk-sleep-go-copilot/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org/doc/go1.25)
[![CDK Version](https://img.shields.io/badge/CDK-2.1125.0-orange.svg)](https://docs.aws.amazon.com/cdk/)

**cdk-sleep-go-copilot** is a fully serverless, event-driven sleep audio pipeline built with the AWS CDK in Go. The pipeline automatically processes audio files uploaded to S3 through a sophisticated orchestration using EventBridge, Step Functions, Lambda, Amazon Polly, DynamoDB, and SNS — all without any always-on compute resources.

The project follows a strict **Test-Driven Development (TDD)** discipline: every infrastructure change begins with a failing `go test`, and no code is committed until both `go test ./...` and `cdk synth` pass locally.

## 🎯 Key Features

- **📤 Event-Driven Architecture:** S3 uploads automatically trigger processing via EventBridge
- **🔀 Workflow Orchestration:** AWS Step Functions manages multi-step audio processing
- **🎙️ Audio Processing:** Lambda function with Amazon Polly text-to-speech synthesis
- **💾 Metadata Storage:** DynamoDB tracks processing status and output locations
- **🔔 Notifications:** SNS topics for success and failure events with KMS encryption
- **📊 Observability:** X-Ray tracing, CloudWatch Logs, and CloudWatch Alarms
- **🌍 Multi-Environment:** Context-driven deployment to dev/stage/prod environments
- **🔒 Security First:** Encrypted storage, private buckets, least-privilege IAM, input validation
- **✅ Comprehensive Testing:** 54 CDK infrastructure tests ensuring pipeline integrity

## 📐 Architecture Overview

```
┌─────────┐      ┌──────────┐      ┌───────────────┐      ┌────────┐      ┌─────────┐
│ Client  │─────▶│  S3 Input │─────▶│ EventBridge   │─────▶│  Step  │─────▶│ Lambda  │
│ Upload  │      │  Bucket   │      │ Rule (Event)  │      │Functions│      │Processor│
└─────────┘      └──────────┘      └───────────────┘      └────────┘      └─────────┘
                                                                │                  │
                                                                │                  │
                                                                ▼                  ▼
                                                          ┌─────────┐        ┌─────────┐
                                                          │   SNS   │        │   S3    │
                                                          │ Topics  │        │ Output  │
                                                          └─────────┘        └─────────┘
                                                                │                  │
                                                                │                  ▼
                                                                │            ┌─────────┐
                                                                │            │DynamoDB │
                                                                │            │Metadata │
                                                                │            └─────────┘
                                                                ▼
                                                          ┌─────────┐
                                                          │Subscribers│
                                                          └─────────┘
```

For detailed architecture information including complete data flow, component inventory, and Mermaid diagrams, see [ARCHITECTURE.md](ARCHITECTURE.md).

## 🚀 Quick Start

### Prerequisites

- **Go:** 1.21 or later ([install](https://golang.org/doc/install))
- **Node.js:** 18+ ([install](https://nodejs.org/))
- **AWS CDK CLI:** 2.1125.0 (`npm install -g aws-cdk@2.1125.0`)
- **AWS Account:** With configured credentials (`aws configure`)

### Installation & Setup

```bash
# Clone the repository
git clone https://github.com/obstreperous-ai/cdk-sleep-go-copilot.git
cd cdk-sleep-go-copilot

# Install Go dependencies
go mod download

# Install AWS CDK CLI (if not already installed)
npm install -g aws-cdk@2.1125.0

# Verify installation
go version
cdk --version
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestEndToEndPipelineIntegration
```

### Deploying the Stack

```bash
# Synthesize CloudFormation template (dry run)
cdk synth

# Deploy to development environment
cdk deploy -c env=dev

# Deploy to staging environment
cdk deploy -c env=stage

# Deploy to production environment
cdk deploy -c env=prod

# View differences before deployment
cdk diff -c env=dev
```

### Testing the Pipeline

1. **Upload a file to the input bucket:**
   ```bash
   aws s3 cp test-audio.mp3 s3://<input-bucket-name>/test-audio.mp3
   ```

2. **Monitor the execution:**
   - View Step Functions execution in AWS Console
   - Check CloudWatch Logs for detailed processing logs
   - Review X-Ray traces for performance analysis

3. **Verify processing results:**
   ```bash
   # Check DynamoDB for metadata
   aws dynamodb get-item --table-name <table-name> --key '{"audioId": {"S": "test-audio.mp3"}}'
   
   # Verify output file exists
   aws s3 ls s3://<output-bucket-name>/processed-*
   ```

4. **Subscribe to notifications:**
   ```bash
   # Subscribe to success notifications
   aws sns subscribe --topic-arn <completed-topic-arn> --protocol email --notification-endpoint your@email.com
   
   # Subscribe to failure notifications
   aws sns subscribe --topic-arn <failed-topic-arn> --protocol email --notification-endpoint your@email.com
   ```

## 🧪 Test-Driven Development Rules

This project follows **strict TDD discipline**. All contributions must adhere to these rules:

1. **Write a failing test first** — Always commit the test before the implementation
2. **Write the minimum code** to make the test pass — No speculative logic
3. **`go test ./...` must pass** before any Go source commit
4. **`cdk synth` must succeed** before any CDK stack commit
5. **Update `ARCHITECTURE.md`** whenever the infrastructure topology changes
6. **Conventional commits only** (`feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:`)

See [CONTRIBUTING.md](CONTRIBUTING.md) for the complete TDD workflow and contribution guidelines.

## 📚 Documentation

- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Complete system architecture, data flow, and Mermaid diagrams
- **[SUMMARY.md](SUMMARY.md)** - Project completion status, key decisions, and lessons learned
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - TDD workflow and contribution guidelines
- **[.github/AGENT_GUIDELINES.md](.github/AGENT_GUIDELINES.md)** - Agent persona and development rules

## 🛠️ Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Infrastructure as Code** | AWS CDK (Go) | Define cloud resources programmatically |
| **Compute** | AWS Lambda (Go) | Serverless audio processing function |
| **Orchestration** | AWS Step Functions | Multi-step workflow coordination |
| **Storage** | Amazon S3 | Input/output audio file storage |
| **Database** | Amazon DynamoDB | Audio metadata and processing status |
| **Event Bus** | Amazon EventBridge | Event-driven triggers |
| **Notifications** | Amazon SNS | Success/failure notifications |
| **AI/ML** | Amazon Polly | Text-to-speech synthesis |
| **Observability** | CloudWatch + X-Ray | Logging, metrics, and distributed tracing |
| **Security** | AWS KMS, IAM | Encryption and access control |
| **CI/CD** | GitHub Actions | Automated testing and validation |

## 🔐 Security Features

- **Private S3 Buckets:** Public access blocked at bucket level
- **Encryption at Rest:** S3 (S3-managed), DynamoDB (AWS-managed), SNS (KMS)
- **Encryption in Transit:** TLS for all AWS service communications
- **Least-Privilege IAM:** Granular permissions for each service
- **Input Validation:** Multi-layer validation (Step Functions + Lambda)
- **Path Traversal Protection:** File key sanitization
- **Key Rotation:** Automatic KMS key rotation enabled
- **Audit Logging:** CloudWatch Logs for all operations

## 📊 Multi-Environment Support

The stack supports multiple deployment environments with different configurations:

| Environment | Removal Policy | Auto-Delete S3 | Use Case |
|-------------|---------------|----------------|----------|
| **dev** | DESTROY | ✅ Enabled | Rapid iteration and testing |
| **stage** | RETAIN | ❌ Disabled | Pre-production validation |
| **prod** | RETAIN | ❌ Disabled | Production workloads |

Deploy to specific environments using CDK context:

```bash
cdk deploy -c env=dev    # Development
cdk deploy -c env=stage  # Staging
cdk deploy -c env=prod   # Production
```

## 📈 Observability & Monitoring

### Logs
- **Step Functions:** ALL level logging to CloudWatch
- **Lambda:** Structured JSON logs with request IDs
- **Retention:** 7 days (configurable per environment)

### Tracing
- **X-Ray:** Active tracing on Lambda and Step Functions
- **Service Map:** Visualize request flow across services
- **Performance Analysis:** Identify bottlenecks and latency

### Alarms
- **State Machine Failures:** Triggers on execution failures
- **Lambda Errors:** Monitors invocation errors
- **SNS Integration:** Failure alarms publish to failed topic

## 🧮 Useful Commands

| Command | Description |
|---------|-------------|
| `go test ./...` | Run all CDK infrastructure tests |
| `go test -v ./...` | Run tests with verbose output |
| `cdk synth` | Synthesize CloudFormation template |
| `cdk deploy` | Deploy stack to AWS |
| `cdk deploy -c env=dev` | Deploy to development environment |
| `cdk diff` | Show differences with deployed stack |
| `cdk destroy` | Delete the stack from AWS |
| `go mod download` | Download Go dependencies |
| `go mod tidy` | Clean up Go module dependencies |

### Building Lambda Function Locally

```bash
cd lambda/audio-processor
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
cd ../..
```

## 📝 Project Status

✅ **COMPLETE** - All core functionality implemented and validated (Issue #12)

- [x] S3 Input/Output Buckets with encryption and versioning
- [x] EventBridge Rule for Object Created events
- [x] Step Functions State Machine with complete orchestration
- [x] Lambda Function with full audio processing
- [x] Amazon Polly integration for text-to-speech
- [x] DynamoDB metadata storage and tracking
- [x] SNS notification topics (success/failure) with KMS encryption
- [x] Comprehensive error handling and retry policies
- [x] Multi-environment support (dev/stage/prod)
- [x] X-Ray tracing and CloudWatch monitoring
- [x] CloudWatch Alarms for proactive monitoring
- [x] 54 comprehensive infrastructure tests
- [x] Complete documentation with architecture diagrams

See [SUMMARY.md](SUMMARY.md) for detailed project completion status and metrics.

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for the TDD workflow and guidelines.

**Key Points:**
- All changes must follow strict TDD (test-first development)
- Tests must pass before submitting PR
- Update ARCHITECTURE.md for infrastructure changes
- Use conventional commit messages
- Keep code clean and well-documented

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [AWS CDK](https://aws.amazon.com/cdk/) in Go
- Developed using strict Test-Driven Development principles
- Powered by GitHub Copilot as Senior AWS CDK TDD Specialist
- Follows [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/) principles

---

**Ready to deploy?** Follow the Quick Start guide above or see [ARCHITECTURE.md](ARCHITECTURE.md) for detailed system documentation.
