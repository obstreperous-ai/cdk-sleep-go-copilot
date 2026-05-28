# Contributing to cdk-sleep-go-copilot

Thank you for contributing! This project follows a strict **Test-Driven Development (TDD)** workflow. Please read and follow the guidelines below before opening a pull request.

---

## Code of Conduct

Be respectful, constructive, and inclusive. See [GitHub's Community Guidelines](https://docs.github.com/en/site-policy/github-terms/github-community-guidelines) for details.

---

## TDD Workflow (mandatory)

1. **Write a failing test first.** Before any implementation, add a test in `*_test.go` that describes the desired behavior and confirm it fails (`go test ./...`).
2. **Write the minimal code to make the test pass.** No extra logic, no premature abstractions.
3. **Refactor** with tests still green.
4. **Run `cdk synth`** to confirm the CloudFormation template is valid.
5. **Update `ARCHITECTURE.md`** if the infrastructure topology changes.

Never open a pull request where `go test ./...` or `cdk synth` fail.

---

## Branching & Commit Conventions

- Branch names: `feat/<short-description>`, `fix/<short-description>`, `chore/<short-description>`.
- Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/):
  - `feat:` new feature
  - `fix:` bug fix
  - `chore:` tooling, CI, dependency updates
  - `docs:` documentation only
  - `test:` test-only changes
  - `refactor:` code restructuring without behavior change

---

## Pull Request Checklist

- [ ] Tests written **before** implementation (TDD).
- [ ] `go test ./...` passes locally.
- [ ] `cdk synth` passes locally.
- [ ] `ARCHITECTURE.md` updated (if infra changed).
- [ ] Conventional commit message used.
- [ ] PR description references the relevant GitHub issue.

---

## CDK Construct Preferences

- Prefer **L2 constructs** over L1 (CloudFormation resources) where available.
- Use **L3 patterns** (e.g., `awss3deployments`) for common multi-resource combinations.
- Follow [AWS Well-Architected Framework](https://aws.amazon.com/architecture/well-architected/) principles.

---

## Running Tests Locally

```bash
# Unit tests
go test ./...

# CDK synthesis (requires Node.js + aws-cdk installed)
npm install -g aws-cdk
go mod download
cdk synth
```

---

## Questions?

Open a GitHub Discussion or reference the relevant issue in your PR. The agent persona and detailed rules are in [.github/AGENT_GUIDELINES.md](.github/AGENT_GUIDELINES.md).
