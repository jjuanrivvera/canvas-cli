# Architecture

This document describes the internal architecture of Canvas CLI.

## Project Structure

```
canvas-cli/
├── cmd/canvas/           # Application entry point
├── commands/             # CLI command definitions (Cobra)
├── internal/
│   ├── api/             # Canvas API client and services
│   ├── auth/            # OAuth 2.0 + PKCE authentication
│   ├── config/          # Configuration management
│   ├── cache/           # Response caching
│   ├── batch/           # Concurrent batch operations
│   ├── repl/            # Interactive shell
│   └── output/          # Output formatters
├── pkg/                 # Public packages
├── docs/                # Documentation
└── test/                # Test fixtures
```

## Component Overview

```mermaid
graph LR
    subgraph User
        CLI[CLI Commands]
        REPL[REPL Shell]
    end

    subgraph Core
        API[API Client]
        AUTH[Auth Manager]
        CFG[Config Manager]
    end

    subgraph Features
        CACHE[Cache]
        BATCH[Batch Processor]
        OUT[Output Formatter]
    end

    CLI --> API
    REPL --> CLI
    API --> AUTH
    API --> CACHE
    API --> BATCH
    CLI --> OUT
    AUTH --> CFG
```

## Core Components

### API Client

The API client (`internal/api/`) provides a type-safe interface to the Canvas REST API.

```mermaid
classDiagram
    class Client {
        +BaseURL string
        +HTTPClient *http.Client
        +Get(path, params) Response
        +Post(path, body) Response
        +Put(path, body) Response
        +Delete(path) Response
    }

    class Service {
        +client *Client
    }

    Client <|-- CoursesService
    Client <|-- AssignmentsService
    Client <|-- UsersService
    Client <|-- ModulesService
    Client <|-- PagesService
    Client <|-- DiscussionsService
    Client <|-- CalendarService
    Client <|-- PlannerService
```

**Key features:**
- Automatic pagination handling
- Rate limit awareness
- Exponential backoff retry
- Request/response logging

### Authentication

OAuth 2.0 with PKCE (Proof Key for Code Exchange) for secure authentication.

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Browser
    participant Canvas

    User->>CLI: canvas auth login
    CLI->>CLI: Generate PKCE verifier + challenge
    CLI->>Browser: Open authorization URL
    Browser->>Canvas: User authenticates
    Canvas->>CLI: Authorization code (callback)
    CLI->>Canvas: Exchange code + verifier for token
    Canvas->>CLI: Access token + refresh token
    CLI->>Keyring: Store tokens securely
```

**Token storage priority:**
1. System keyring (macOS Keychain, Windows Credential Manager, Linux Secret Service)
2. Encrypted file fallback (AES-256-GCM)

### Rate Limiting

Adaptive rate limiting respects Canvas API quotas.

```mermaid
graph TD
    REQ[Request] --> CHECK{Check Quota}
    CHECK -->|>50%| NORMAL[5 req/sec]
    CHECK -->|20-50%| WARN[2 req/sec]
    CHECK -->|<20%| CRITICAL[1 req/sec]

    NORMAL --> SEND[Send Request]
    WARN --> SEND
    CRITICAL --> SEND

    SEND --> RESP[Response Headers]
    RESP --> UPDATE[Update Quota State]
    UPDATE --> CHECK
```

### Caching

Smart caching with TTL-based invalidation:

| Resource | TTL |
|----------|-----|
| Courses | 15 minutes |
| Users | 5 minutes |
| Assignments | 10 minutes |
| Modules | 10 minutes |

### Batch Processing

Concurrent processing with configurable parallelism:

```mermaid
graph LR
    INPUT[Items] --> POOL[Worker Pool]
    POOL --> W1[Worker 1]
    POOL --> W2[Worker 2]
    POOL --> W3[Worker 3]
    POOL --> W4[Worker 4]
    POOL --> W5[Worker 5]

    W1 --> COLLECT[Collector]
    W2 --> COLLECT
    W3 --> COLLECT
    W4 --> COLLECT
    W5 --> COLLECT

    COLLECT --> RESULT[Results + Errors]
```

## Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.21+ |
| CLI Framework | [Cobra](https://github.com/spf13/cobra) |
| Configuration | [Viper](https://github.com/spf13/viper) |
| OAuth | golang.org/x/oauth2 |
| Rate Limiting | golang.org/x/time/rate |
| Keyring | zalando/go-keyring |
| Logging | log/slog (stdlib) |

## Design Principles

1. **Security First** - All credentials encrypted, no hardcoded secrets
2. **Graceful Degradation** - Fallbacks for keyring, network, and API issues
3. **User Experience** - Progress indicators, helpful error messages
4. **Testability** - Interface-driven design, mock-friendly
5. **Performance** - Caching, batching, concurrent operations

## Error Handling

Custom error types with actionable suggestions:

```go
type CanvasError struct {
    StatusCode int
    Message    string
    Suggestion string
    DocURL     string
}
```

## Testing Strategy

- Unit tests for all services
- Integration tests with mock HTTP server
- 90% code coverage target
- Race condition detection enabled
