# CLAUDE.md — wabaman-contrib

Contributor reference for AI assistants working on `github.com/pedidopago/wabaman-contrib/v2`.

## What This Repository Is

A multi-language helper library for integrating with **WABAMan** — Pedido Pago's WhatsApp Business Account management platform. It provides:

- **Go packages** (primary): typed clients, API structs, utility functions
- **TypeScript/JavaScript package** (secondary): companion client published as `@pedidopago/wabaman-contrib`

---

## Repository Layout

```
wabaman-contrib/
├── event/           # Event structs (template approval/rejection lifecycle)
├── fbgraph/         # Facebook Graph API client (send messages, upload media)
├── msgdriver/       # Message driver interfaces (document message types)
├── rest/            # WABAMan REST API types and client
│   └── client/      # REST client implementation
├── shared-types/    # Shared type definitions (e.g. Referral)
├── util/            # Utility functions (template variable counting/validation)
├── whapi/           # WhatsApp API common types and graph objects
├── wsapi/           # WebSocket API types (voice/video call lifecycle)
├── js/              # TypeScript/JavaScript package
│   ├── lib/         # Source TypeScript files
│   ├── test/        # Mocha + Chai tests
│   ├── dist/        # Compiled output (not committed)
│   ├── package.json
│   └── tsconfig.json
├── docs/
│   └── calls.md     # WebRTC call flow documentation (Portuguese)
├── .github/workflows/test.yml
├── go.mod
└── whatsapp_apis.md
```

---

## Go Development

### Module

```
module github.com/pedidopago/wabaman-contrib/v2
go 1.26.3
```

The `/v2` suffix is part of the import path: every import is
`github.com/pedidopago/wabaman-contrib/v2/...`. It was taken when the `fbgraph`
client started requiring a `context.Context`, so v1 consumers keep resolving and
building until they choose to migrate.

### Key Dependencies

| Package | Purpose |
|---|---|
| `github.com/rs/zerolog` | Structured logging |
| `github.com/google/go-querystring` | URL query string encoding |
| `github.com/pedidopago/go-common` | Internal shared utilities |

> `go-common` is a **private** pedidopago repository. CI accesses it via `GIT_USER`/`GIT_TOKEN` secrets configured in GitHub Actions. Local development requires the same git credential setup.

### Running Tests

```bash
go test ./...
```

### Linting / Static Analysis

```bash
go vet ./...
go test -race ./...
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck -go 1.26 ./...
```

CI (`.github/workflows/test.yml`) runs `go test -race ./...` and
`golangci-lint`, which has `staticcheck` enabled, so the standalone run above is
for local use. It triggers on pushes to `main` and pull requests targeting
`main` -- not on every branch. Run `-race` locally too: this package hands a
request to a goroutine that writes a variable the caller reads back.

---

## TypeScript/JavaScript Development

The JS package lives in `js/` and is published independently as `@pedidopago/wabaman-contrib` (MIT, no runtime dependencies).

### Build

```bash
cd js
npm run build   # tsc → dist/
```

Output: `dist/index.js` + `dist/index.d.ts` (declaration files).

### Tests

```bash
npm test          # env-cmd + mocha (compiled JS)
npm run test-ts   # ts-node + mocha (TypeScript directly)
```

Tests live in `js/test/` with naming convention `*.test.ts`.

---

## CI/CD

**GitHub Actions** (`.github/workflows/test.yml`):
- Triggers on every branch push
- Sets up Go ≥1.23 and configures git token auth for private pedidopago repos
- Runs: `go test ./...` → `go vet ./...` → `staticcheck`
- No deployment step; library consumers pull it as a module dependency

---

## Package-by-Package Guide

### `fbgraph` — Facebook Graph API Client

Direct client for the Meta/Facebook Graph API (WhatsApp Business).

**Entry point:** `fbgraph.NewClient(accessToken string) *Client`

Key methods on `*Client`:
- `SendMessage(ctx, phoneID, msg)` — POST to `/messages`
- `SendMarketingMessage(ctx, phoneID, msg)` — POST to `/marketing_messages`
- `UploadMedia(ctx, phoneID, mimeType, r, fsize, filename)` — multipart upload
- `GetMedia(ctx, mediaID)` / `DownloadMedia(ctx, mr, out)` — retrieve media
- `NewUploadSession(ctx, fbAppID, params)` / `UploadHeaderHandle(ctx, sessionID, r)` — resumable uploads

Every method that issues a request takes a `context.Context` first, and every
request is built with `NewRequestWithContext`. `NewRequest` is deprecated and
must not be used: a request without a context ignores cancellation, so the
caller waits out the 120s `DefaultHTTPClient` timeout on work already doomed.
`TestNoContextFreeRequests` enforces both halves — including that the context
passed is the caller's, not a fresh `context.Background()`.

`SendMessage` and `SendMarketingMessage` are not idempotent: `context.Canceled`
from them does not mean the message was not delivered, so do not resend on it.

`UploadMedia`'s `fsize` went from ignored to enforced in v2. v1 accepted the
argument and never looked at it; v2 fails the upload when fewer bytes are copied
than `fsize` promises, because a reader that simply ends reports `io.EOF` rather
than an error and a truncated upload is otherwise indistinguishable from a whole
one. A call site that passes a size from somewhere other than the bytes being
sent -- a database column, a value computed earlier -- compiles unchanged and
starts failing. Pass the real length, or 0 if it is genuinely unknown.

`LastGraphError()` / `LastErrorRawBody()` are per-call as of v2: every request
method clears them on entry, so they always describe the call that just
returned. In v1 only a *failing* call refreshed them, so an error could be read
long afterwards; now a successful call in between wipes it. Read them
immediately after the failure you care about. `TestErrorAccessorsResetOnEntry`
enforces this across the package, following calls transitively -- a method that
reports a Graph error through a helper is held to the same rule.

Package-level globals:
- `DefaultGraphAPIVersion = "v23.0"` — override per-client via `Client.GraphAPIVersion`
- `DefaultHTTPClient` — 120s timeout; replace for testing
- `DebugTrace bool` — prints request URLs and bodies to stdout when true

Error handling: non-200 responses are parsed into `*GraphError`. Use `AsGraphError(err)` to type-assert. `ErrApplicationRateLimitReached` is a sentinel for Graph error code 4.

### `rest` — WABAMan REST API Types

Type definitions for the WABAMan service REST API. No client logic here — just request/response structs.

Key types: `NewMessageRequest`, `NewMessageResponse`, `PreviewMessageOutcomeRequest`, `NewContactRequest`, `UpdateContactRequest`, `GetContactsV2Request`, `GetMessagesRequest`, `AnyMessage` (interface), `SentMessage`, `ReceivedMessage`, `ErrorResponse`.

`ErrorResponse` implements `error` and carries `StatusCode` + raw body for non-JSON responses.

### `rest/client` — WABAMan REST Client

**Entry point:** `client.Client{JWT, BaseURL, Origin}`

Default base URL: `https://api.first.v2.pedidopago.com.br/wabaman`

All methods accept `context.Context` as first argument.

Key methods:
- `NewMessage` / `PreviewMessageOutcome` / `NewMessageReaction`
- `NewContact` / `UpdateContact` / `GetContactByID` / `GetContactByBranchIDAndWABAContactID` / `GetContactsV2`
- `GetMessages` / `GetMessageByID` / `UpdateMessages`
- `NewNote` / `RegisterClientMessage`
- `NewTemplate` / `TemplateExists`
- `CheckIntegration` / `CheckIntegrationV2`
- `GetBusinesses` / `GetPhones`
- `BusinessContactBroadcast` / `PhoneContactBroadcast`

The `X-Origin` header is set when `Client.Origin` is non-empty. JWT auth is Bearer token.

`GetMessageByID` deserializes dynamically: it inspects `message.object_type` (`"host"` → `SentMessage`, otherwise `ReceivedMessage`) before full decode.

### `wsapi` — WebSocket API Types

Structs for real-time WebSocket messages used in voice/video call flows:
- `IncomingCallFromClient`, `SetupCallFromBrowser`, `TerminateCall`, `CallConsumed`
- `ICECandidate`, `CallOnAnswerSDP`, `CallStartTimer`, `SendBrowserCandidate`

See `docs/calls.md` for the full call state-machine (Portuguese).

### `event` — Template Lifecycle Events

`TemplateEvent` with `EventKind` constants: `APPROVED`, `REJECTED`, `ARCHIVED`, `DELETED`, `DISABLED`, `FLAGGED`, `PAUSED`, `REINSTATED`, `PENDING_DELETION`.

### `msgdriver` — Message Driver Interfaces

Generic interfaces and document message type constants. Used for implementing pluggable message drivers.

### `whapi` — WhatsApp API Common Types

Shared WhatsApp primitives (`MediaObject`, template component types, error codes) referenced across packages.

### `shared-types` — Cross-Package Shared Types

Currently contains `Referral` struct for WhatsApp referral tracking.

### `util` — Utility Functions

- `CountAndValidateTemplateVariables(text string) (count int, valid bool)` — validates `{{1}}`, `{{2}}` … `{{N}}` placeholders are sequential and counts them. Regex: `\{\{([1-9]*[0-9])\}\}`.
- `IsBSUID(s string) bool` — checks if a string is a WhatsApp Business-Scoped User ID.

---

## Code Conventions

### Go

- Each package has a `doc.go` with package-level documentation.
- Constructor pattern: `NewClient(…) *Type` returns a pointer.
- Errors wrap with `fmt.Errorf("context: %w", err)`.
- Logging uses `zerolog` (`github.com/rs/zerolog/log`).
- Type aliases used for semantic safety (e.g. `type MessageType string`).
- Constants grouped with `const (…)` blocks, not scattered.
- Tests live alongside source (`foo_test.go` next to `foo.go`), using `package foo` (white-box) or `package foo_test` (black-box).

### TypeScript

- Enums for all discriminated union values (`WSMessageType`, `WSErrorCode`, `FBGraphTemplateComponentType`, etc.).
- Interfaces for all data shapes; classes only when behavior is needed (`ParsedTemplate`, `ParsedTemplateHeader`).
- Template variable syntax matches Go: `{{1}}`, `{{2}}` etc.
- Compiled output goes to `dist/` — never commit `dist/`.
- Default fallback language: `pt_BR`.

### Template Variables

Both Go and TypeScript use `{{N}}` placeholders (1-indexed, sequential). The Go `util.CountAndValidateTemplateVariables` enforces sequential ordering — gaps are invalid.

---

## Private Repository Access

`go-common` is private. To develop locally, configure git credentials:

```bash
git config --global url."https://<USER>:<TOKEN>@github.com/pedidopago".insteadOf "https://github.com/pedidopago"
```

In CI this is handled automatically via the `GIT_USER` and `GIT_TOKEN` repository secrets.

---

## Key External APIs

- **Facebook Graph API** — `https://graph.facebook.com/{version}/{phoneID}/messages`
- **WABAMan REST API** — `https://api.first.v2.pedidopago.com.br/wabaman`
- **WhatsApp Business API docs** — see `whatsapp_apis.md` for official links
