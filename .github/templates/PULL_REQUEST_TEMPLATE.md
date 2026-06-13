# Pull Request Template

> **Purpose:** Ensure all PRs follow TDD discipline and include necessary documentation updates.

---

## 📝 Description

### What does this PR do?

[Clear description of changes made in this PR]

### Which issue does this close?

Closes #[issue number]

---

## ✅ TDD Checklist

Confirm all TDD workflow steps were followed:

- [ ] **Test First:** Failing test(s) written before implementation
- [ ] **Test Committed:** Test(s) committed separately before implementation
- [ ] **Minimal Implementation:** Only code necessary to pass tests was added
- [ ] **Tests Pass:** `go test ./...` (or equivalent) passes locally
- [ ] **Synthesis Succeeds:** `cdk synth` (or equivalent) succeeds locally
- [ ] **Conventional Commit:** Used conventional commit format (`feat:`, `fix:`, etc.)

---

## 📚 Documentation Updates

- [ ] **ARCHITECTURE.md** updated (if infrastructure topology changed)
  - [ ] Prose description updated
  - [ ] Mermaid diagram updated (if applicable)
  - [ ] Component inventory updated
- [ ] **README.md** updated (if user-facing changes)
- [ ] **CONTRIBUTING.md** updated (if workflow changed)
- [ ] Code comments added for complex logic

---

## 🔐 Security Review

- [ ] IAM policies follow least-privilege principle
- [ ] No wildcard (`*`) permissions added
- [ ] Encryption enabled for data at rest
- [ ] Encryption enabled for data in transit
- [ ] Public access blocked where applicable
- [ ] Input validation implemented
- [ ] No secrets or credentials in code
- [ ] Security best practices followed

---

## 🧪 Testing

### Tests Added/Modified

List new or modified tests:
- `Test[Name]` - [Brief description]
- `Test[Name]` - [Brief description]

### Test Coverage

- **Total Tests:** [number]
- **Tests Passing:** [number]
- **Coverage:** [percentage or description]

### How to Test Locally

```bash
# Commands to run tests locally
go test ./...
cdk synth
```

---

## 📊 Changes Made

### Files Changed

- `[filename]` - [Brief description of changes]
- `[filename]` - [Brief description of changes]

### Resources Added/Modified

List AWS resources affected:
- **Added:** [Resource type] - [Resource name/description]
- **Modified:** [Resource type] - [Resource name/description]
- **Removed:** [Resource type] - [Resource name/description]

---

## 🚀 Deployment Notes

### Breaking Changes

- [ ] **No breaking changes**
- [ ] **Contains breaking changes** (describe below)

[Describe any breaking changes and migration steps]

### Environment-Specific Considerations

**Development:**
- [Any dev-specific notes]

**Staging:**
- [Any stage-specific notes]

**Production:**
- [Any prod-specific notes]

### Rollback Plan

[Describe how to rollback these changes if needed]

---

## 📸 Screenshots / Diagrams

[If applicable, add screenshots of:
- AWS Console showing deployed resources
- CloudWatch dashboards
- Architecture diagram changes
- Test results]

---

## 🔗 Related Links

- Issue: #[issue number]
- Related PRs: #[pr number]
- AWS Documentation: [link]
- Design Doc: [link]

---

## 👀 Reviewer Checklist

Reviewers should verify:

- [ ] PR description clearly explains changes
- [ ] Tests were written before implementation
- [ ] All tests pass locally and in CI
- [ ] CDK synthesis succeeds
- [ ] Architecture documentation updated if needed
- [ ] Security best practices followed
- [ ] Code is readable and well-commented
- [ ] Conventional commit message used
- [ ] No unnecessary changes included

---

## 📝 Additional Notes

[Any additional context, considerations, or discussion points for reviewers]

---

**PR Author:** @[username]  
**Reviewers:** @[username], @[username]  
**Labels:** [feat/fix/docs/chore]
