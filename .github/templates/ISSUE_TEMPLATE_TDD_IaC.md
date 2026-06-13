# TDD IaC Issue Template

> **Purpose:** Use this template for issues involving Test-Driven Development (TDD) Infrastructure as Code (IaC) changes.

---

## Issue Title Format

```
[Issue Number] Category: Concise Problem Statement
```

**Examples:**
- `[15] Feature: Add SQS Dead Letter Queue for Failed Processing`
- `[16] Security: Implement VPC Endpoints for S3 Access`
- `[17] Observability: Add Custom CloudWatch Metrics`

---

## 🎯 Goal

[Clear, one-sentence objective describing what this issue accomplishes]

**Example:**  
_"Add a Dead Letter Queue (DLQ) to capture and persist messages that fail processing after maximum retry attempts."_

---

## 📋 Requirements

List specific, testable requirements:

1. [Specific requirement 1]
2. [Specific requirement 2]
3. [Specific requirement 3]
4. ...

**Example:**
1. Create SQS Dead Letter Queue with encryption enabled
2. Configure maximum receive count of 3 on main queue
3. Add CloudWatch Alarm for messages in DLQ
4. Update IAM permissions for queue access
5. Verify DLQ integration with CDK assertions

---

## 🔧 Strict Discipline (must follow)

- [ ] Start by reviewing `ARCHITECTURE.md` to understand current system design
- [ ] Write failing test(s) first before any implementation
- [ ] Implement minimal code to make tests pass
- [ ] Run `go test ./...` (or equivalent) - must pass
- [ ] Run `cdk synth` (or equivalent) - must succeed
- [ ] Update `ARCHITECTURE.md` if infrastructure topology changes
- [ ] Use conventional commit messages (`feat:`, `fix:`, `test:`, etc.)
- [ ] Verify CI passes before marking issue complete

---

## ✅ Success Criteria

- [ ] Test(s) written and committed before implementation
- [ ] All tests pass (`go test ./...`)
- [ ] CDK synthesis succeeds (`cdk synth`)
- [ ] `ARCHITECTURE.md` updated if topology changed (diagram + prose)
- [ ] Conventional commit message used
- [ ] CI workflow passes
- [ ] Code review approved (if required)

**Additional Criteria:**
- [ ] [Specific criterion 1]
- [ ] [Specific criterion 2]

---

## 📚 Context & References

**Related Issues:**
- Depends on: #[issue number]
- Blocks: #[issue number]
- Related to: #[issue number]

**Relevant Documentation:**
- [Link to AWS documentation]
- [Link to CDK construct documentation]
- [Link to internal architecture docs]

**Design Decisions:**
- [Key design decision 1]
- [Key design decision 2]

---

## 🔍 Testing Strategy

**Unit Tests:**
```
Describe tests to verify:
- Resource exists with correct properties
- IAM permissions are correct
- Security settings applied
- Integration points configured
```

**Integration Tests:**
```
Describe tests to verify:
- End-to-end flow works
- Error handling behaves correctly
- Multi-resource interactions function
```

---

## 🚧 Implementation Notes

**Suggested Approach:**
1. [Step-by-step implementation suggestion]
2. [Key considerations or gotchas]

**Security Considerations:**
- [Security requirement 1]
- [Security requirement 2]

**Performance Considerations:**
- [Performance requirement 1]
- [Performance requirement 2]

---

## 📝 Definition of Done

This issue is complete when:
- [ ] All success criteria met
- [ ] Tests written before implementation
- [ ] All tests passing
- [ ] CDK synth succeeds
- [ ] Architecture documentation updated
- [ ] PR reviewed and approved
- [ ] Changes merged to main branch

---

**Issue Created By:** [GitHub Username]  
**Assigned To:** [Agent/Human]  
**Estimated Effort:** [S/M/L/XL]  
**Priority:** [Low/Medium/High/Critical]
