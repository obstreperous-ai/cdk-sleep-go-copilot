# Templates Directory

This directory contains reusable templates for TDD Infrastructure as Code (IaC) development with AI agents.

## 📋 Available Templates

### 1. [ISSUE_TEMPLATE_TDD_IaC.md](ISSUE_TEMPLATE_TDD_IaC.md)

**Purpose:** Structured issue template for TDD IaC changes.

**Use when:**
- Creating new infrastructure features
- Fixing infrastructure bugs
- Adding observability or security features
- Making any infrastructure change that requires testing

**Key Sections:**
- Clear goal definition
- Testable requirements
- Strict TDD discipline checklist
- Success criteria
- Testing strategy
- Security and performance considerations

### 2. [PULL_REQUEST_TEMPLATE.md](PULL_REQUEST_TEMPLATE.md)

**Purpose:** Comprehensive PR template ensuring TDD compliance and documentation updates.

**Use when:**
- Opening pull requests for infrastructure changes
- Ensuring all TDD workflow steps were followed
- Documenting changes for reviewers

**Key Sections:**
- TDD checklist verification
- Documentation update tracking
- Security review checklist
- Testing information
- Deployment notes and rollback plan

### 3. [AGENT_PROMPT_TEMPLATE.md](AGENT_PROMPT_TEMPLATE.md)

**Purpose:** Structured prompt format for instructing AI agents on TDD IaC work.

**Use when:**
- Delegating infrastructure tasks to AI agents
- Ensuring consistent agent behavior
- Providing clear context and constraints
- Enforcing TDD workflow

**Key Sections:**
- Agent persona definition
- Task context and requirements
- Step-by-step TDD workflow
- Success criteria
- Security requirements
- Self-review checklist

## 🎯 How to Use These Templates

### For GitHub Issues

Copy [ISSUE_TEMPLATE_TDD_IaC.md](ISSUE_TEMPLATE_TDD_IaC.md) content when creating new issues:

```bash
# When creating a new issue in GitHub, paste the template and fill in:
- Goal (one sentence)
- Specific requirements (numbered list)
- Success criteria (checkboxes)
- Testing strategy
```

### For Pull Requests

You can set this as the default PR template by copying it to `.github/PULL_REQUEST_TEMPLATE.md`:

```bash
cp .github/templates/PULL_REQUEST_TEMPLATE.md .github/PULL_REQUEST_TEMPLATE.md
```

Or manually use when creating PRs by pasting the content into the PR description.

### For AI Agent Instructions

When working with GitHub Copilot or other AI tools:

1. **Copy** [AGENT_PROMPT_TEMPLATE.md](AGENT_PROMPT_TEMPLATE.md)
2. **Fill in** the placeholders:
   - `[LANGUAGE]` - Go, TypeScript, Python, etc.
   - `[DOMAIN]` - Event-driven architectures, microservices, etc.
   - `[Project Name]` - Your project name
   - `[test command]` - Your test command
   - `[synth command]` - Your synthesis command
3. **Customize** the persona and requirements for your specific task
4. **Provide** to the AI agent as context

## 🔄 Template Customization

These templates are **starting points** - customize them for your project:

### Recommended Customizations

1. **Language/Framework Specific:**
   - Update test command syntax
   - Adjust file naming conventions
   - Customize testing frameworks

2. **Organization Specific:**
   - Add company security policies
   - Include required approvers
   - Add compliance checklists

3. **Project Specific:**
   - Reference project-specific docs
   - Add custom testing patterns
   - Include deployment procedures

### Example Customizations

**For Python CDK:**
```bash
# Change test command from:
go test ./...

# To:
pytest tests/
```

**For TypeScript CDK:**
```bash
# Change test command from:
go test ./...

# To:
npm test
```

## 📚 Related Documentation

- **[META-PROMPTS.md](../META-PROMPTS.md)** - Comprehensive meta-prompting patterns and guidelines
- **[AGENT_GUIDELINES.md](../AGENT_GUIDELINES.md)** - Agent persona and strict rules for this project
- **[CONTRIBUTING.md](../../CONTRIBUTING.md)** - Full contribution guidelines

## 🎓 Best Practices

### Issue Creation

1. ✅ **Be Specific** - Clear, testable requirements
2. ✅ **Define Success** - Explicit success criteria
3. ✅ **Consider Security** - Include security considerations upfront
4. ✅ **Plan Tests** - Outline testing strategy before implementation

### Pull Requests

1. ✅ **Complete Checklist** - Verify all TDD steps followed
2. ✅ **Update Docs** - Include documentation updates
3. ✅ **Provide Context** - Explain why changes were made
4. ✅ **Plan Rollback** - Document how to undo changes

### Agent Instructions

1. ✅ **Clear Persona** - Define agent role explicitly
2. ✅ **Strict Workflow** - Enforce TDD steps with "must" language
3. ✅ **Concrete Examples** - Show expected test and code structure
4. ✅ **Self-Review** - Include checklist for agent to verify work

## 🤝 Contributing Improvements

If you discover improvements to these templates:

1. Test the improvement in your project
2. Document the benefit it provides
3. Submit a PR with the enhancement
4. Update this README with usage notes

## 📝 Version History

- **v1.0** (2026-06-12) - Initial templates extracted from cdk-sleep-go-copilot
  - Issue template for TDD IaC
  - Pull request template with TDD checklist
  - Agent prompt template

---

**Maintained By:** cdk-sleep-go-copilot project  
**License:** MIT  
**Last Updated:** 2026-06-12
