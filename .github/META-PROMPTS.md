# Meta-Prompts & Reusable Patterns for Agentic TDD IaC Projects

> **Purpose:** This document extracts the reusable meta-prompting patterns, agent personas, and workflow templates from the `cdk-sleep-go-copilot` project. These patterns can be applied to future agentic Test-Driven Development (TDD) Infrastructure as Code (IaC) projects using AWS CDK or similar tools.

---

## Table of Contents

- [Overview](#overview)
- [Agent Persona Template](#agent-persona-template)
- [Strict TDD Rules (Reusable)](#strict-tdd-rules-reusable)
- [Source of Truth Pattern](#source-of-truth-pattern)
- [Issue-Driven Development Workflow](#issue-driven-development-workflow)
- [CDK Construct Hierarchy Guidelines](#cdk-construct-hierarchy-guidelines)
- [Conventional Commit Standards](#conventional-commit-standards)
- [Testing Patterns for IaC](#testing-patterns-for-iac)
- [Documentation Synchronization Pattern](#documentation-synchronization-pattern)
- [Security & Compliance Checklist](#security--compliance-checklist)
- [Multi-Environment Configuration Pattern](#multi-environment-configuration-pattern)
- [Observability Standards](#observability-standards)
- [Meta-Prompt Templates](#meta-prompt-templates)

---

## Overview

The `cdk-sleep-go-copilot` project demonstrated a highly disciplined approach to building AWS infrastructure using:

1. **Pure Issue-Driven Development** - Every change starts with a GitHub issue
2. **Strict TDD** - Every implementation begins with a failing test
3. **Living Documentation** - ARCHITECTURE.md stays perfectly in sync with code
4. **AI-Powered Development** - GitHub Copilot as Senior AWS CDK TDD Specialist

This document captures the **meta-patterns** that made this approach successful, so they can be replicated in other projects.

---

## Agent Persona Template

### Recommended Persona Definition

```markdown
You are a **Senior AWS CDK [LANGUAGE] TDD Specialist**. 

**Core Characteristics:**
- Use clean [LANGUAGE] idioms and best practices
- Write tests first, then minimal code to pass them
- Always follow strict TDD: failing test → minimal implementation → refactor
- Keep architecture documentation perfectly in sync after every change
- Prefer L2/L3 CDK constructs over L1 CloudFormation resources
- Follow AWS Well-Architected Framework principles
- Never deploy until tests + synth succeed locally
- Apply least-privilege IAM policies
- Encrypt data at rest and in transit by default
```

**Customization Points:**
- Replace `[LANGUAGE]` with: Go, TypeScript, Python, Java, or C#
- Add domain-specific expertise (e.g., "with expertise in event-driven architectures")
- Include project-specific frameworks or tools

---

## Strict TDD Rules (Reusable)

Apply these rules to any IaC project to enforce TDD discipline:

### The 10 Commandments of IaC TDD

1. **Test First, Always**  
   Every implementation must be preceded by a failing test in `*_test.[ext]`. Commit the failing test before the implementation.

2. **Conventional Commits Only**  
   Use `feat:`, `fix:`, `chore:`, `docs:`, `test:`, or `refactor:` prefixes for all commits.

3. **Tests Must Pass**  
   All unit tests must pass before any commit that touches source files.

4. **Synthesis Must Succeed**  
   `cdk synth` (or equivalent) must succeed before any commit that modifies infrastructure code.

5. **Documentation Stays in Sync**  
   Update architecture documentation whenever the infrastructure topology changes. Never let docs diverge from deployed infrastructure.

6. **Minimal Code Only**  
   Write the smallest implementation that makes tests pass. Avoid speculative abstractions or "future-proofing."

7. **High-Level Constructs Preferred**  
   Use L2/L3 constructs (or equivalent high-level abstractions). Resort to L1/CloudFormation only when no alternative exists.

8. **Least-Privilege IAM**  
   Never attach `*` actions or resources in IAM policies. Grant only the minimum permissions required.

9. **No Secrets in Source**  
   Never commit credentials, tokens, or account IDs. Use CDK context, SSM Parameter Store, or Secrets Manager.

10. **Reference Official Docs**  
    When uncertain, consult official framework documentation and best practices.

---

## Source of Truth Pattern

### Architecture Document as Single Source of Truth

**Pattern:** Designate a single architecture document (e.g., `ARCHITECTURE.md`) as the canonical reference for:
- System design and topology
- Data flow and component inventory
- Visual diagrams (Mermaid, PlantUML, etc.)
- Service choices and integration points

**Implementation Rules:**

```markdown
## Before Starting Any Issue

1. Read `ARCHITECTURE.md` to understand where your change fits
2. Implement only what the current issue requires, consistent with the documented design
3. If an issue changes the topology, update `ARCHITECTURE.md` (description AND diagram) in the SAME pull request

## Never Allow Divergence

- Architecture doc updates are NOT optional
- Code and docs must be committed together
- CI should enforce this (e.g., require ARCHITECTURE.md update for infra changes)
```

**Benefits:**
- Prevents architectural drift
- Provides clear context for AI agents and human developers
- Creates reliable system documentation that stays current
- Enables confident refactoring with documented baseline

---

## Issue-Driven Development Workflow

### Pure Issue-Based Development

**Pattern:** Every code change originates from a GitHub issue with clear requirements.

**Issue Template Structure:**

```markdown
# [Issue Number] Title: Concise Problem Statement

## Goal
Clear, one-sentence objective.

## Requirements
1. Specific requirement 1
2. Specific requirement 2
3. ...

## Strict Discipline (must follow)
- Start with failing test
- Implement minimal code
- Update ARCHITECTURE.md if topology changes
- Use conventional commits
- Verify tests + synth pass

## Success Criteria
- [ ] Criterion 1 met
- [ ] Criterion 2 met
- [ ] Tests passing
- [ ] Documentation updated
```

**Workflow:**

```
GitHub Issue → Plan → Failing Test(s) → Minimal Implementation → 
Tests Pass → Synth Succeeds → Update Docs → Commit → PR → Review → Merge
```

---

## CDK Construct Hierarchy Guidelines

### Construct Preference Pattern

Apply this decision hierarchy when choosing CDK constructs:

| Priority | Level | Examples | When to Use |
|---|---|---|---|
| ✅ **First Choice** | **L3 Patterns** | `awss3deployments`, `awsecs_patterns` | Multi-resource patterns with opinionated defaults |
| ✅ **Default** | **L2 Constructs** | `s3.Bucket`, `lambda.Function` | Standard AWS resources with sensible defaults |
| ⚠️ **Last Resort** | **L1 Cfn Resources** | `CfnBucket`, `CfnFunction` | Only when L2/L3 unavailable or insufficient |

**Decision Algorithm:**

```
Need to provision resource?
  └─> Check: Does L3 pattern exist for this use case?
      ├─> YES → Use L3 pattern
      └─> NO → Check: Does L2 construct exist?
          ├─> YES → Use L2 construct
          └─> NO → Check: Is L2 insufficient for requirements?
              ├─> YES → Use L1, document reason in code comments
              └─> NO → Reconsider requirements
```

---

## Conventional Commit Standards

### Commit Message Format

```
<type>(<optional-scope>): <description>

[optional body]

[optional footer]
```

### Standard Types

| Type | Use When | Examples |
|---|---|---|
| `feat:` | Adding new feature/capability | `feat: add S3 input bucket with encryption` |
| `fix:` | Fixing a bug or issue | `fix: correct IAM policy for Lambda execution` |
| `test:` | Adding or updating tests | `test: add DynamoDB integration tests` |
| `docs:` | Documentation only | `docs: update ARCHITECTURE.md with new flow` |
| `chore:` | Tooling, dependencies, CI | `chore: upgrade CDK to v2.120.0` |
| `refactor:` | Code restructuring, no behavior change | `refactor: extract S3 bucket construct` |
| `perf:` | Performance improvements | `perf: increase Lambda memory to 512MB` |
| `style:` | Formatting, whitespace | `style: fix indentation in state machine` |

### Examples from Project

```bash
feat: add Step Functions state machine with input validation
test: add EventBridge rule target verification
fix: update DynamoDB PointInTimeRecovery to use PointInTimeRecoverySpecification
docs: add Mermaid diagram for complete data flow
chore: add CI workflow with go test and cdk synth
refactor: extract error handling to reusable Pass states
```

---

## Testing Patterns for IaC

### CDK Assertion Testing Patterns

**Pattern 1: Resource Existence**
```go
// Verify critical resources exist
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
})
```

**Pattern 2: Security Validation**
```go
// Ensure security best practices
template.HasResourceProperties(jsii.String("AWS::S3::Bucket"), map[string]interface{}{
    "PublicAccessBlockConfiguration": map[string]interface{}{
        "BlockPublicAcls": true,
        "BlockPublicPolicy": true,
        "IgnorePublicAcls": true,
        "RestrictPublicBuckets": true,
    },
})
```

**Pattern 3: End-to-End Flow Validation**
```go
// Verify complete integration
func TestEndToEndPipelineIntegration(t *testing.T) {
    // Verify S3 → EventBridge → Step Functions → Lambda flow
    // Check all connections exist
    // Validate IAM permissions chain
    // Confirm error handling paths
}
```

**Pattern 4: Resource Count Validation**
```go
// Prevent resource proliferation
template.ResourceCountIs(jsii.String("AWS::Lambda::Function"), jsii.Number(3))
```

---

## Documentation Synchronization Pattern

### Living Architecture Documentation

**File Structure:**

```
docs/
  ARCHITECTURE.md       # System design, data flow, Mermaid diagrams
  SUMMARY.md            # Project completion status, decisions, lessons learned
  CONTRIBUTING.md       # TDD workflow, contribution guidelines
  .github/
    AGENT_GUIDELINES.md # Agent persona, strict rules, workflow
    META-PROMPTS.md     # Reusable patterns (this file)
```

**Update Triggers:**

| Change Type | Required Doc Updates |
|---|---|
| New AWS resource | ARCHITECTURE.md (prose + diagram) |
| Data flow modification | ARCHITECTURE.md (data flow section) |
| New integration | ARCHITECTURE.md (component inventory) |
| Testing pattern | CONTRIBUTING.md (if novel pattern) |
| Issue completion | SUMMARY.md (lessons learned) |
| Agent behavior update | AGENT_GUIDELINES.md |

**Enforcement:**

```yaml
# CI check example
- name: Check Architecture Documentation
  run: |
    # Verify ARCHITECTURE.md was updated if infra changed
    if git diff --name-only HEAD~1 | grep -E '(cdk-|construct|stack)'; then
      if ! git diff --name-only HEAD~1 | grep 'ARCHITECTURE.md'; then
        echo "ERROR: Infrastructure changed but ARCHITECTURE.md not updated"
        exit 1
      fi
    fi
```

---

## Security & Compliance Checklist

### Infrastructure Security Patterns

Apply these security principles to all IaC projects:

**Storage Security:**
- [ ] S3 buckets have public access blocked
- [ ] S3 buckets use encryption at rest (S3-managed or KMS)
- [ ] S3 buckets have versioning enabled for data protection
- [ ] DynamoDB tables use encryption (AWS-managed or customer-managed)

**Network Security:**
- [ ] Resources deployed in private subnets where applicable
- [ ] Security groups follow least-privilege rules
- [ ] VPC endpoints used for AWS service access (avoid internet routing)

**IAM Security:**
- [ ] All roles follow least-privilege principle
- [ ] No wildcard (`*`) permissions in production
- [ ] Service-specific roles (not shared across services)
- [ ] Assume role policies restrict to specific services

**Data Security:**
- [ ] SNS topics encrypted with KMS
- [ ] SQS queues encrypted at rest
- [ ] Secrets stored in Secrets Manager or Parameter Store
- [ ] KMS keys have rotation enabled

**Observability Security:**
- [ ] CloudWatch Logs have retention policies
- [ ] X-Ray tracing enabled for distributed systems
- [ ] Alarms configured for security events
- [ ] Audit logging enabled (CloudTrail)

**Input Validation:**
- [ ] Step Functions validate inputs with Choice states
- [ ] Lambda functions sanitize file paths
- [ ] API Gateway validates request schemas
- [ ] DynamoDB uses schema validation where applicable

---

## Multi-Environment Configuration Pattern

### Context-Driven Environment Pattern

**Pattern:** Use CDK context to drive environment-specific behavior without code duplication.

**Implementation:**

```go
// Read environment context
env := "dev" // default
if envContext := stack.Node().TryGetContext(jsii.String("env")); envContext != nil {
    if envStr, ok := envContext.(string); ok {
        env = envStr
    }
}

// Apply environment-specific policies
var removalPolicy awscdk.RemovalPolicy
var autoDeleteObjects *bool

switch env {
case "dev":
    removalPolicy = awscdk.RemovalPolicy_DESTROY
    autoDeleteObjects = jsii.Bool(true)
case "stage", "prod":
    removalPolicy = awscdk.RemovalPolicy_RETAIN
    autoDeleteObjects = jsii.Bool(false)
default:
    removalPolicy = awscdk.RemovalPolicy_DESTROY
    autoDeleteObjects = jsii.Bool(true)
}
```

**Environment Configuration Matrix:**

| Environment | Removal Policy | Auto-Delete S3 | CloudWatch Retention | Cost Profile |
|---|---|---|---|---|
| **dev** | DESTROY | ✅ Enabled | 3 days | Minimal |
| **stage** | RETAIN | ❌ Disabled | 7 days | Moderate |
| **prod** | RETAIN | ❌ Disabled | 30 days | Full |

**Deployment Commands:**

```bash
cdk deploy -c env=dev    # Development environment
cdk deploy -c env=stage  # Staging environment
cdk deploy -c env=prod   # Production environment
```

---

## Observability Standards

### Comprehensive Observability Pattern

Apply these observability patterns to all serverless infrastructure:

**1. Structured Logging**
```go
// Lambda structured logging example
log.Printf(`{"level":"INFO","requestId":"%s","event":"processing_started","bucket":"%s","key":"%s"}`,
    requestID, bucket, key)
```

**2. X-Ray Tracing**
```go
// Enable X-Ray on Lambda
lambda.NewFunction(stack, jsii.String("Processor"), &lambda.FunctionProps{
    Tracing: lambda.Tracing_ACTIVE,
    // ...
})

// Enable X-Ray on Step Functions
sfn.NewStateMachine(stack, jsii.String("StateMachine"), &sfn.StateMachineProps{
    TracingEnabled: jsii.Bool(true),
    // ...
})
```

**3. CloudWatch Alarms**
```go
// State Machine failure alarm
cloudwatch.NewAlarm(stack, jsii.String("StateMachineFailureAlarm"), &cloudwatch.AlarmProps{
    Metric: stateMachine.MetricFailed(&awscloudwatch.MetricOptions{
        Statistic: "Sum",
        Period: awscdk.Duration_Minutes(jsii.Number(5)),
    }),
    Threshold: jsii.Number(1),
    EvaluationPeriods: jsii.Number(1),
    AlarmDescription: jsii.String("Alert when state machine execution fails"),
})
```

**4. Log Retention Policies**
```go
logs.NewLogGroup(stack, jsii.String("StateMachineLogGroup"), &logs.LogGroupProps{
    Retention: logs.RetentionDays_ONE_WEEK,
    RemovalPolicy: removalPolicy,
})
```

---

## Meta-Prompt Templates

### Template 1: TDD IaC Issue Meta-Prompt

Use this meta-prompt when creating issues for TDD IaC work:

```markdown
You are implementing infrastructure following strict Test-Driven Development.

**Context:**
- Project: [PROJECT_NAME]
- Framework: AWS CDK in [LANGUAGE]
- Architecture: See ARCHITECTURE.md

**Requirements:**
[SPECIFIC_REQUIREMENTS]

**Strict TDD Workflow:**
1. Write failing test(s) in [TEST_FILE]
2. Commit the failing test
3. Implement minimal code to pass the test
4. Run: [TEST_COMMAND] (must pass)
5. Run: [SYNTH_COMMAND] (must succeed)
6. Update ARCHITECTURE.md if topology changed
7. Commit with conventional commit message

**Success Criteria:**
- [ ] Test(s) written first and committed
- [ ] Implementation makes tests pass
- [ ] [TEST_COMMAND] succeeds
- [ ] [SYNTH_COMMAND] succeeds
- [ ] ARCHITECTURE.md updated if needed
- [ ] Conventional commit used
```

### Template 2: Security Review Meta-Prompt

Use this meta-prompt for security-focused reviews:

```markdown
Review the infrastructure code for security best practices.

**Focus Areas:**
- [ ] IAM policies follow least-privilege
- [ ] Data encryption at rest and in transit
- [ ] Public access blocked where applicable
- [ ] Input validation implemented
- [ ] Secrets not in source code
- [ ] Network security configured properly

**Framework:** AWS CDK
**Context:** [PROJECT_CONTEXT]

Provide specific findings with file locations and remediation steps.
```

### Template 3: Documentation Sync Meta-Prompt

Use this meta-prompt for documentation updates:

```markdown
Update documentation to match infrastructure changes.

**Changed Files:**
[LIST_OF_CHANGED_FILES]

**Required Updates:**
1. ARCHITECTURE.md:
   - [ ] Update prose description
   - [ ] Update Mermaid diagram
   - [ ] Update component inventory
2. README.md:
   - [ ] Update if user-facing changes
3. CONTRIBUTING.md:
   - [ ] Update if workflow changes

**Verification:**
- [ ] All links work
- [ ] Diagrams match deployed topology
- [ ] No broken references
```

### Template 4: Agent Self-Review Meta-Prompt

Use this meta-prompt for agent self-review before finalizing work:

```markdown
Before completing this task, perform self-review:

**TDD Compliance:**
- [ ] Did I write tests before implementation?
- [ ] Do all tests pass?
- [ ] Does CDK synth succeed?

**Code Quality:**
- [ ] Minimal implementation (no speculative code)?
- [ ] L2/L3 constructs used where possible?
- [ ] Conventional commit message format?

**Documentation:**
- [ ] ARCHITECTURE.md updated if topology changed?
- [ ] Documentation matches code?
- [ ] All references and links valid?

**Security:**
- [ ] Least-privilege IAM?
- [ ] Encryption enabled?
- [ ] No secrets in source?

If any item is ✗, address it before completing.
```

---

## Applying These Patterns

### Quick Start Checklist

To apply these meta-prompts to a new project:

1. **Setup:**
   - [ ] Copy AGENT_GUIDELINES.md and customize persona
   - [ ] Create ARCHITECTURE.md with initial design
   - [ ] Setup TDD testing framework
   - [ ] Configure CI with test + synth checks

2. **Process:**
   - [ ] Create issues using issue template
   - [ ] Follow TDD workflow for each issue
   - [ ] Keep documentation in sync
   - [ ] Use conventional commits

3. **Quality:**
   - [ ] Apply security checklist to all resources
   - [ ] Implement observability patterns
   - [ ] Configure multi-environment support
   - [ ] Review using meta-prompt templates

---

## Lessons Learned

### What Worked Well in cdk-sleep-go-copilot

1. **Strict TDD discipline** prevented bugs and ensured comprehensive test coverage
2. **Living architecture documentation** kept the team aligned and facilitated onboarding
3. **Issue-driven development** provided clear scope and progress tracking
4. **AI agent with clear persona** produced consistent, high-quality code
5. **Conventional commits** made git history readable and useful

### Anti-Patterns to Avoid

1. ❌ **Implementing before testing** - leads to incomplete test coverage
2. ❌ **Letting docs diverge from code** - creates confusion and maintenance burden
3. ❌ **Skipping synthesis checks** - allows invalid templates to be committed
4. ❌ **Using L1 constructs unnecessarily** - increases code complexity
5. ❌ **Broad IAM permissions** - violates security best practices

---

## Additional Resources

- [AWS CDK Best Practices](https://docs.aws.amazon.com/cdk/v2/guide/best-practices.html)
- [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/)
- [Conventional Commits Specification](https://www.conventionalcommits.org/)
- [GitHub Copilot Documentation](https://docs.github.com/en/copilot)

---

## Contributing to These Patterns

If you've used these meta-prompts in your project and discovered improvements or new patterns, please contribute back:

1. Fork the repository
2. Add your pattern with clear examples
3. Document the use case and benefits
4. Submit a pull request with conventional commit

**These patterns are living documents** - improve them as you learn!

---

**Generated from:** `cdk-sleep-go-copilot` project  
**License:** MIT  
**Last Updated:** 2026-06-12
