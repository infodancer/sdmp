# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SDMP (Secure Domain Messaging Protocol) is the server-to-server federation protocol for the infodancer messaging stack. It handles all inter-domain communication — no client ever speaks this protocol.

This is a new protocol designed from scratch. See the parent `infodancer/infodancer` repo for the full requirements document (`docs/next-gen-messaging-protocol.md`) and the protocol outline (`docs/protocol-outlines.md`).

### Responsibilities

- Message notification (encrypted UDP signals to receiving domains)
- Message fetch (recipient domain retrieves message from sender domain)
- Domain key publication and discovery
- Domain-to-domain authentication
- Reputation queries (federated, opt-in)
- Gateway-to-gateway bulk transfer
- Responsibility transfer (when a gateway accepts notification duty)
- Msgcoin ledger synchronization

### Key Design Principles

- **Sender-stores, recipient-pulls**: sender bears storage cost until recipient explicitly accepts
- **No rejection at notification time**: silence is the only response to notifications
- **Reject early, never bounce**: inherited from the smtpd philosophy
- **Encrypted envelopes**: intermediate servers are blind caches seeing only the destination domain
- **Message IDs as capability tokens**: knowing the ID is sufficient to fetch ciphertext

## Technology

- **Go** for the server implementation
- **Protocol Buffers** for message serialization across all channels
- **gRPC** (mTLS) for authenticated domain-to-domain RPC (key exchange, reputation, responsibility transfer)
- **HTTP** for CDN-cacheable blob fetch (`GET /message/{id}` — message ID as capability token, no auth required)
- **UDP + protobuf** for lightweight message notifications (fire-and-forget, no handshake)
- **DNS SRV** (`_mail._tcp.domain`) for domain discovery

## Development Commands

This module uses [Task](https://taskfile.dev/) as the build tool:

```bash
task build          # Build the binary
task test           # Run tests
task lint           # Run golangci-lint
task vulncheck      # Check for vulnerabilities
task all            # Run all checks (build, lint, vulncheck, test)
task test:coverage  # Run tests with coverage report
task install:deps   # Install golangci-lint and govulncheck
task hooks:install  # Configure git to use .githooks directory
```

## Module Structure

```
/cmd/sdmp/main.go       # Entrypoint only, minimal logic
/internal/sdmp/          # Domain-specific implementation
/proto/                  # Protocol buffer definitions
/errors/                 # Centralized error definitions
```

## Development Workflow

### Branch and Issue Protocol

**This workflow is MANDATORY.** All significant work must follow this process with no exceptions:

1. **Create a GitHub issue first** - Before creating a branch, draft an issue describing the purpose and design based on your understanding of the user's request. Assign the issue to the user who requested it. Ask the user to approve the issue before proceeding.

2. **Create a feature or content branch** - Only after issue approval, create the branch. Use descriptive names that include the issue id like `feature/UUID` or `bug/UUID`.

3. **Reference the issue in all commits** - Every commit message and pull request must include the issue URL.

4. **Stay focused on the issue** - Make only changes directly related to the approved issue. Do not refactor unrelated code, fix unrelated bugs, or make "improvements" outside the scope.

5. **Handle unrelated problems separately** - If you notice bugs, technical debt, or potential issues unrelated to your current work, ask the user to approve creating a separate GitHub issue. Do not address them in the current branch.

### Pull Request Workflow

- All branches merge to main via PR
- PRs should reference the originating issue
- **NEVER ask users to merge or approve a PR** - PR approval and merging must always be manual actions taken by the user
- After creating a PR, checkout the main branch before starting any further work

### Security Best Practices

- Never commit secrets, API keys, credentials, or tokens
- Use `crypto/rand` for random number generation in security contexts
- Validate all external input at system boundaries
- Use parameterized queries to prevent injection
- Regularly audit dependencies with `govulncheck`

Read CONVENTIONS.md for Go coding standards.
