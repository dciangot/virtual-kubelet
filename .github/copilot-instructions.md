# Copilot Instructions for virtual-kubelet

Purpose

This file provides concise, repository-specific instructions to help Copilot-style assistants and contributors quickly understand how to build, test, lint, and navigate the Virtual Kubelet codebase.

Build, test, and lint commands

- Build the main binary (from repo root):
  - make build
  - or: CGO_ENABLED=0 go build -a --tags "netgo osusergo $(VK_BUILD_TAGS)" -ldflags '-extldflags "-static"' -o bin/virtual-kubelet github.com/virtual-kubelet/virtual-kubelet
- Run full test suite:
  - make test
  - or: go test ./... (the Makefile runs go vet and GODEBUG=cgocheck=2 in non-CI)
- Run a single package's tests or a single test:
  - go test ./providers/azure -run '^TestName$' -v
  - go test ./vkubelet -run '^TestSomething$'
  - Use -run with a regex to target tests; add -v for verbose output.
- Coverage and reports:
  - make cover (generates merged coverage and HTML in CI/local as configured)
- Formatting & linting:
  - make format  (uses goimports)
  - go vet ./... (Makefile runs this in test)
  - The repo uses goimports for formatting; ensure goimports is installed (make setup installs tools).
- Docker / reproducible build:
  - make docker
  - make safebuild (builds inside Docker)

High-level architecture (big picture)

- Language & layout: Go project using modules (go.mod). The runnable binary is under cmd/ ("virtual-kubelet"), with the main runtime under packages such as vkubelet/ and providers/.
- Providers: providers/ contains pluggable backends (azure, aws, alicloud, etc.). Providers implement the Provider interface (see vkubelet/provider.go) with methods like CreatePod, UpdatePod, DeletePod, GetPod, GetPodStatus, GetPods, Capacity, NodeConditions, OperatingSystem.
- Core runtime: vkubelet implements the kubelet-facing logic (node status, pod lifecycle orchestration) and delegates provider-specific operations to implementations under providers/.
- Shared code: internal/ and errdefs/ contain internal utilities and standardized error definitions used across providers and core code.
- Testing: unit tests are distributed across packages; some provider tests (e.g., Azure) require external credentials or environment variables.

Key conventions and repository-specific patterns

- Use context.Context in public/exported APIs. Public APIs should accept context.Context for tracing, cancellation, and logging.
- Error handling: use the repo's errdefs package for defining and checking error categories rather than relying on concrete error types.
- Provider interface: providers must be pluggable and not assume direct access to the Kubernetes API server; use callbacks for secrets/configmaps where applicable.
- Build tags and flags: builds commonly use tags "netgo osusergo" and may include VK_BUILD_TAGS. The Makefile encodes these conventions.
- Tooling: Makefile targets wrap common workflows (build, test, format, cover, docker). make setup bootstraps developer tools (goimports, gocovmerge, goreleaser, dep). Keep CI vs local differences in mind: make test behaves differently under CI.
- Tests with external deps: some provider tests require credentials (Azure expects credentials.json at repo root or AZURE_AUTH_LOCATION). Document any required env vars in provider READMEs.
- Formatting: goimports is the canonical formatter; run make format before committing.

Files and docs to reference

- README.md — usage, provider guides, and high-level project description.
- CONTRIBUTING.md — contribution process, CLA, maintainers, and test guidance.
- providers/<provider>/README.md — provider-specific setup and test notes (e.g., azure credentials instructions).

AI assistant / other assistant configs

- No repository-specific Copilot / Claude / Cursor / Aider configs were detected in top-level or .github (if an assistant-specific config exists, include its contents here).

Notes for Copilot sessions

- Start by running `make test` to learn common failures and the test surface for the branch you check out.
- When asked to check out or prepare PR branches: `gh pr checkout <number>` will create a local branch tracking the PR head (e.g., `gh pr checkout 1295`).
- To run a single failing test locally, prefer `go test <pkg> -run '^TestName$' -v` to iterate quickly.

Relevant repository conventions copied from docs

- Contributors must sign the CNCF CLA before PRs can be accepted (see CONTRIBUTING.md).
- Public APIs should take context.Context; errors should use errdefs; avoid breaking API changes unless part of a planned major release.

If anything in this file should be expanded (e.g., per-provider test setup, CI matrix notes, or additional helper commands), say which area to expand and a small example to include.

MCP Servers: Kubernetes e2e runner

- Recommended MCP server: Kubernetes e2e runner to execute e2e targets and verify cluster integration.
- How Copilot can use it:
  - Use the repo's Makefile.e2e targets. Example commands the runner should support:
    - make e2e KUBECONFIG=/path/to/kubeconfig (builds bin/e2e/virtual-kubelet, deploys via skaffold, runs e2e tests under internal/test/e2e)
    - make skaffold MODE=run (deploy using hack/skaffold/virtual-kubelet/skaffold.yml)
    - make e2e.clean (cleanup deployed resources)
  - Required runner capabilities:
    - A Kubernetes cluster (kind/minikube/docker-desktop) available to the runner with kubectl configured
    - Docker build support to create the e2e binary image or cross-build artifacts
    - Ability to run skaffold and kubectl

- Notes:
  - The e2e Makefile expects the current kubectl context to be compatible (kind or minikube/docker-desktop) and will fail otherwise.
  - The e2e target sets VK_BUILD_TAGS+=mock_provider so the runner must permit build tags during Go builds.

If you want, add other MCP servers to run unit tests or coverage collectors; specify which ones to add and their expected commands.
