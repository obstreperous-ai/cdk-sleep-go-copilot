# cdk-sleep-go-copilot

**cdk-sleep-go-copilot** is a fully serverless, event-driven sleep audio pipeline built with the AWS CDK in Go. Raw audio files uploaded to S3 trigger an EventBridge rule that invokes a Lambda processor. The processor transcribes the audio via Amazon Transcribe, stores enriched session metadata in DynamoDB, persists the transcript JSON to a results S3 bucket, and notifies downstream subscribers through SNS — all without any always-on compute. The project follows a strict **Test-Driven Development (TDD)** discipline: every infrastructure change begins with a failing `go test`, and no code is committed until both `go test ./...` and `cdk synth` pass locally. See [ARCHITECTURE.md](ARCHITECTURE.md) for the full pipeline description and Mermaid diagram, and [CONTRIBUTING.md](CONTRIBUTING.md) for the contribution workflow.

## Strict TDD Rules

1. **Write a failing test first** — always commit the test before the implementation.
2. **Write the minimum code** to make the test pass — no speculative logic.
3. **`go test ./...` must pass** before any Go source commit.
4. **`cdk synth` must succeed** before any CDK stack commit.
5. **Update `ARCHITECTURE.md`** whenever the infrastructure topology changes.
6. **Conventional commits only** (`feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:`).

## Useful commands

 * `go test ./...`   run unit tests
 * `cdk synth`       emits the synthesized CloudFormation template
 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
