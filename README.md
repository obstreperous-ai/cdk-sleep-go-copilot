# cdk-sleep-go-copilot

[![CI](https://github.com/obstreperous-ai/cdk-sleep-go-copilot/actions/workflows/ci.yml/badge.svg)](https://github.com/obstreperous-ai/cdk-sleep-go-copilot/actions/workflows/ci.yml)
[![Tests](https://img.shields.io/badge/tests-54%20passing-success.svg)](cdk-base_test.go)
[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org/doc/go1.25)
[![CDK Version](https://img.shields.io/badge/CDK-2.1125.0-orange.svg)](https://docs.aws.amazon.com/cdk/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

**cdk-sleep-go-copilot** is a fully serverless, event-driven sleep audio pipeline built with the AWS CDK in Go. The pipeline automatically processes audio files uploaded to S3 through a sophisticated orchestration using EventBridge, Step Functions, Lambda, Amazon Polly, DynamoDB, and SNS — all without any always-on compute resources.

The project follows a strict **Test-Driven Development (TDD)** discipline: every infrastructure change begins with a failing `go test`, and no code is committed until both `go test ./...` and `cdk synth` pass locally.

---

## 📖 Table of Contents

- [Why This Project?](#-why-this-project)
- [Key Features](#-key-features)
- [Architecture Overview](#-architecture-overview)
- [Quick Start](#-quick-start)
  - [Prerequisites](#prerequisites)
  - [Installation & Setup](#installation--setup)
  - [Running Tests](#running-tests)
  - [Deploying the Stack](#deploying-the-stack)
  - [Testing the Pipeline](#testing-the-pipeline)
- [Test-Driven Development Rules](#-test-driven-development-rules)
- [Experiment Methodology](#-experiment-methodology)
- [Meta-Prompting & Agentic Development](#-meta-prompting--agentic-development)
- [Documentation](#-documentation)
- [Technology Stack](#-technology-stack)
- [Security Features](#-security-features)
- [Multi-Environment Support](#-multi-environment-support)
- [Observability & Monitoring](#-observability--monitoring)
- [Useful Commands](#-useful-commands)
- [Project Status](#-project-status)
- [Contributing](#-contributing)
- [License](#-license)
- [Acknowledgments](#-acknowledgments)

---

## 🎓 Why This Project?

This project is an **experiment in AI-powered Test-Driven Infrastructure as Code** development. Key objectives:

1. **Prove TDD Works for IaC** - Demonstrate that strict test-first development can be applied successfully to infrastructure code, not just application code.

2. **Establish Agentic Development Patterns** - Explore how AI agents (GitHub Copilot) can effectively build complex infrastructure when guided by clear personas, strict rules, and comprehensive tests.

3. **Create Reusable Meta-Prompts** - Extract and document patterns that can be applied to future IaC projects. See [META-PROMPTS.md](.github/META-PROMPTS.md) for reusable templates.

4. **Maintain Living Documentation** - Show that architecture documentation can stay perfectly synchronized with code through disciplined workflows.

5. **Validate Pure Issue-Driven Development** - Every change originates from a GitHub issue with clear success criteria, creating complete traceability.

**What Makes This Different:**
- 🤖 **AI-First Development** - Built entirely by GitHub Copilot with human oversight
- 📋 **100% Issue-Driven** - Zero ad-hoc changes, every commit traces to an issue
- ✅ **Test Coverage Before Features** - 54 infrastructure tests ensuring reliability
- 📚 **Living Architecture** - ARCHITECTURE.md stays perfectly in sync with deployed infrastructure
- 🔄 **Reproducible Patterns** - Documented meta-prompts enable replication

**Who Should Use This:**
- Platform engineers exploring AI-assisted IaC development
- Teams wanting to improve IaC testing discipline
- Projects requiring comprehensive infrastructure documentation
- Anyone interested in agentic software development patterns

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

---

## 🔬 Experiment Methodology

This project serves as a **controlled experiment** in AI-powered infrastructure development. The methodology:

### Pure Issue-Driven Development

**Every** code change originates from a GitHub issue. No ad-hoc commits, no "quick fixes" outside the issue workflow. This creates:
- **Complete Traceability** - Every line of code traces back to a documented requirement
- **Clear Scope** - Issues define boundaries, preventing scope creep
- **Progress Tracking** - Issue completion = measurable progress
- **Knowledge Capture** - Issue discussions preserve decision context

### Strict TDD Enforcement

The project enforces TDD through **mandatory workflow rules**:

```
1. Create GitHub Issue with clear success criteria
2. Write failing test(s) that validate the requirement
3. Commit the test (separate commit)
4. Write minimal implementation to pass the test
5. Run go test ./... (must pass)
6. Run cdk synth (must succeed)
7. Update ARCHITECTURE.md if topology changed
8. Commit with conventional message
9. Open PR and verify CI passes
```

**No exceptions.** This discipline ensures:
- Infrastructure behavior is tested before implementation
- Tests document intended behavior
- Refactoring is safe (tests catch regressions)
- Coverage remains comprehensive (currently 54 tests)

### Agent Persona & Guidelines

The project uses GitHub Copilot with a clearly defined persona:

> **"Senior AWS CDK Go TDD Specialist"**  
> - Enforces strict TDD workflow  
> - Keeps ARCHITECTURE.md synchronized  
> - Prefers L2/L3 constructs  
> - Applies AWS Well-Architected principles  
> - Never deploys until tests + synth succeed  

See [.github/AGENT_GUIDELINES.md](.github/AGENT_GUIDELINES.md) for the complete persona definition and workflow rules.

### Living Architecture Documentation

**ARCHITECTURE.md is the single source of truth.** Key principles:

1. **Before implementing any issue**, review ARCHITECTURE.md to understand context
2. **If an issue changes topology**, update ARCHITECTURE.md in the same PR
3. **Diagrams and prose stay in sync** - no divergence allowed
4. **Component inventory is complete** - every resource is documented

This creates documentation that developers actually trust and reference.

### Metrics & Success Criteria

The experiment tracks:

| Metric | Target | Actual |
|---|---|---|
| **Test Coverage** | Every resource tested | ✅ 54 tests, 100% resource coverage |
| **Documentation Sync** | Zero divergence | ✅ Perfect sync maintained |
| **TDD Compliance** | 100% test-first | ✅ Zero exceptions across 12 issues |
| **CI Pass Rate** | 100% before merge | ✅ All PRs green before merge |
| **Issue Completion** | All criteria met | ✅ 12/12 issues fully complete |

**Conclusion:** The experiment validates that strict TDD + AI agents + living docs = production-ready infrastructure.

---

## 🤖 Meta-Prompting & Agentic Development

One of the primary goals of this project is to **extract and document reusable patterns** for AI-powered IaC development.

### Key Meta-Patterns Discovered

1. **Agent Persona Definition** - Clear role definition with explicit rules produces consistent results
2. **Strict TDD Rules** - Non-negotiable workflow steps prevent quality degradation
3. **Source of Truth Pattern** - Single architecture document provides reliable context
4. **Conventional Commit Standards** - Structured git history enables automated workflows
5. **Multi-Environment Configuration** - Context-driven patterns avoid code duplication
6. **Security Checklists** - Comprehensive security validation prevents vulnerabilities

### Reusable Meta-Prompts

All patterns have been extracted into **[.github/META-PROMPTS.md](.github/META-PROMPTS.md)**, including:

- ✅ **Agent Persona Template** - Customize for your language/framework
- ✅ **10 Commandments of IaC TDD** - Apply to any infrastructure project
- ✅ **Issue Templates** - Structured requirements for AI agents
- ✅ **Testing Patterns** - CDK assertion examples for common scenarios
- ✅ **Security Checklist** - Comprehensive infrastructure security validation
- ✅ **Observability Standards** - Logging, tracing, and alarming patterns
- ✅ **Documentation Sync Pattern** - Keep architecture docs current

### How to Apply These Patterns

**To replicate this approach in your project:**

1. **Copy** `.github/AGENT_GUIDELINES.md` and customize the persona for your stack
2. **Adopt** the TDD workflow rules (test-first, conventional commits, synth checks)
3. **Create** ARCHITECTURE.md as your single source of truth
4. **Use** the issue templates from META-PROMPTS.md for structured requirements
5. **Apply** security and observability patterns to all resources
6. **Enforce** through CI (test + synth must pass before merge)

**Expected Benefits:**
- 🎯 **Higher Quality** - Tests catch issues before deployment
- 📚 **Better Documentation** - Stays current because it's part of the workflow
- 🤖 **Effective AI Collaboration** - Clear rules produce consistent agent behavior
- 🔍 **Complete Traceability** - Every change has clear origin and justification
- 🚀 **Faster Onboarding** - New team members understand system from docs + tests

### Lessons for Agentic Development

**What Worked:**
- Clear persona with explicit rules → consistent agent behavior
- Test-first discipline → comprehensive validation
- Living architecture docs → reliable context for AI
- Issue-driven workflow → clear scope and progress

**What to Avoid:**
- Vague or contradictory instructions confuse agents
- Skipping tests leads to incomplete implementations
- Allowing doc drift undermines agent context
- Ad-hoc changes break traceability

See [SUMMARY.md](SUMMARY.md) for detailed lessons learned and project retrospective.

---

## 📚 Documentation

Comprehensive documentation provides complete project context:

| Document | Purpose | Audience |
|---|---|---|
| **[README.md](README.md)** | Project overview, quick start, features | All users |
| **[ARCHITECTURE.md](ARCHITECTURE.md)** | Complete system design, data flow, Mermaid diagrams | Developers, architects |
| **[SUMMARY.md](SUMMARY.md)** | Project completion status, key decisions, lessons learned | Project managers, contributors |
| **[CONTRIBUTING.md](CONTRIBUTING.md)** | TDD workflow, contribution guidelines, PR checklist | Contributors |
| **[.github/AGENT_GUIDELINES.md](.github/AGENT_GUIDELINES.md)** | Agent persona, strict rules, development workflow | AI agents, automation |
| **[.github/META-PROMPTS.md](.github/META-PROMPTS.md)** | 🆕 **Reusable meta-prompting patterns and templates** | **Anyone building IaC with AI** |

**Special Note on META-PROMPTS.md:**  
This document extracts the **reusable patterns** from this project that can be applied to future TDD IaC projects. It includes agent persona templates, testing patterns, security checklists, and workflow templates. If you're interested in replicating this development approach, start here.

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
