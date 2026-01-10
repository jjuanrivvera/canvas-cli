# Building a Canvas LMS CLI tool in Go

**Canvas LMS provides a robust REST API with OAuth 2.0 authentication that works well for CLI applications**, though it requires careful handling of token lifecycle and multi-instance support. The optimal Go stack combines Cobra for CLI structure, Viper for configuration, and OS keychain integration for secure credential storage. This guide covers everything needed to build a production-ready Canvas CLI tool.

## Canvas API authentication fundamentals

Canvas offers two authentication approaches: **manual access tokens** for development/testing and **OAuth 2.0** for production applications. The API strictly requires HTTPS and returns JSON responses.

### API token authentication

Manual tokens are generated in the Canvas UI at `/profile` under "Approved Integrations." Tokens grant the same permissions as the user who created them. The **Authorization header format** is:

```bash
curl -H "Authorization: Bearer <ACCESS-TOKEN>" "https://canvas.instructure.com/api/v1/courses"
```

Query string authentication (`?access_token=<TOKEN>`) works but is discouraged for security reasons. Critically, Canvas's Terms of Service **prohibit** asking users to manually generate tokens for multi-user applications—OAuth 2.0 is mandatory for production tools.

### OAuth 2.0 implementation for CLI tools

Canvas fully supports OAuth 2.0 (RFC-6749) with these key endpoints:

| Endpoint | Purpose |
|----------|---------|
| `/login/oauth2/auth` | Authorization (redirect users here) |
| `/login/oauth2/token` | Token exchange and refresh |
| `DELETE /login/oauth2/token` | Token revocation |

**Redirect URI requirements** are crucial for CLI applications. Canvas supports two approaches:

1. **Out-of-band flow**: Use `urn:ietf:wg:oauth:2.0:oob` as redirect URI—Canvas displays the authorization code in the browser for the user to copy
2. **Local server callback**: Register `http://localhost:<port>/callback` or `http://127.0.0.1:<port>/callback` in the developer key configuration

Developer keys must be created by a Canvas administrator at Admin → Account → Developer Keys. The numeric ID becomes `client_id`, and clicking "Show Key" reveals `client_secret`.

### Token lifecycle and refresh mechanism

Access tokens expire after **1 hour** (3600 seconds). Applications must implement refresh token handling:

```bash
POST /login/oauth2/token
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token&
client_id=XXX&
client_secret=YYY&
refresh_token=<refresh_token>
```

The response provides a new access token but **not a new refresh token**—the same refresh token can be reused indefinitely until revoked. When tokens expire, Canvas returns `401 Unauthorized` with a `WWW-Authenticate` header distinguishing expiration from permission issues.

## Canvas API technical implementation

### URL structure and multi-tenancy

Canvas follows a multi-tenant architecture where each institution has its own domain:
- Cloud-hosted: `https://<institution>.instructure.com/api/v1/<resource>`
- Self-hosted: `https://canvas.<institution>.edu/api/v1/<resource>`

**API tokens are instance-specific**—a token from one Canvas instance won't work on another. CLI tools must store the base URL alongside credentials and support multiple configured instances.

### Rate limiting architecture

Canvas uses a **leaky bucket algorithm** with these characteristics:

- Each request costs units deducted from a ~**700-unit quota**
- Quota replenishes automatically faster than real-time consumption
- Rate limits apply **per API token**, not per user or institution
- Parallel requests incur a 50-unit pre-flight penalty (credited back after completion)

Monitor these response headers:
- `X-Request-Cost`: Floating-point cost of current request
- `X-Rate-Limit-Remaining`: Remaining quota (when throttling applies)

When rate-limited, Canvas returns `429 Forbidden`. Sequential requests rarely hit limits, but parallel operations should be limited to **~5 concurrent threads** with exponential backoff on 429 responses.

### Pagination via Link headers

Canvas paginates using HTTP `Link` headers (not query parameters):

```
Link: <https://canvas.example.com/api/v1/courses?page=2&per_page=10>; rel="next",
      <https://canvas.example.com/api/v1/courses?page=5&per_page=10>; rel="last"
```

**Implementation requirements:**
- Default page size is 10 items; use `?per_page=100` for maximum efficiency
- Treat pagination URLs as **opaque**—don't construct them manually
- Access tokens are NOT included in Link URLs; re-append if using query string auth
- Check for `rel="next"` to determine if more pages exist

### Essential API endpoints for CLI tools

**Courses:**
```
GET  /api/v1/courses                     # List user's courses
GET  /api/v1/courses/:id                 # Get course details
GET  /api/v1/courses/:course_id/users    # List course users
```

**Assignments and submissions:**
```
GET  /api/v1/courses/:course_id/assignments                              # List assignments
GET  /api/v1/courses/:course_id/assignments/:id/submissions              # List submissions
PUT  /api/v1/courses/:course_id/assignments/:id/submissions/:user_id     # Grade submission
```

**Users:**
```
GET  /api/v1/users/self          # Current authenticated user
GET  /api/v1/users/:id/profile   # User profile
```

**File operations** follow a three-step upload process: POST metadata to get upload URL, POST file to that URL, then follow redirect to confirm.

### Error handling patterns

Canvas returns structured JSON errors:

```json
{
  "errors": [{"message": "The specified resource does not exist"}],
  "error_report_id": 12345
}
```

Key status codes: **400** (invalid parameters), **401** (auth failure—check `WWW-Authenticate` header), **403** (forbidden or rate-limited), **422** (validation errors with field-specific messages), **429** (rate limit exceeded).

## Go CLI architecture and libraries

### Recommended technology stack

| Component | Library | Rationale |
|-----------|---------|-----------|
| CLI framework | **spf13/cobra** | Industry standard (Kubernetes, Docker, GitHub CLI); nested subcommands, shell completion |
| Configuration | **spf13/viper** | Seamless Cobra integration, multiple config sources, environment variable support |
| HTTP client | **net/http** with custom Transport | Full control over timeouts, connection pooling, TLS settings |
| OAuth 2.0 | **golang.org/x/oauth2** | Official Go OAuth library, handles token refresh |
| Credential storage | **99designs/keyring** or **zalando/go-keyring** | Cross-platform keychain access |
| Rate limiting | **golang.org/x/time/rate** | Token bucket implementation for client-side limiting |

### Project structure

```
canvas-cli/
├── cmd/canvas/
│   ├── main.go                 # Entry point
│   └── root.go                 # Root command, global flags
├── internal/
│   ├── api/
│   │   ├── client.go           # HTTP client wrapper with auth, rate limiting
│   │   ├── courses.go          # Course API methods
│   │   └── pagination.go       # Link header parsing
│   ├── auth/
│   │   ├── oauth.go            # OAuth flow implementation
│   │   ├── token.go            # Token storage/refresh
│   │   └── keyring.go          # Secure credential storage
│   └── config/
│       └── config.go           # Viper configuration
├── commands/                    # Cobra command implementations
│   ├── courses.go
│   ├── assignments.go
│   └── auth.go
├── go.mod
└── Makefile
```

### HTTP client with rate limiting

```go
type CanvasClient struct {
    http     *http.Client
    baseURL  string
    token    string
    limiter  *rate.Limiter
}

func NewClient(baseURL, token string) *CanvasClient {
    transport := &http.Transport{
        TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS12},
        MaxIdleConnsPerHost: 10,  // Increase from default 2
        IdleConnTimeout:     90 * time.Second,
    }
    return &CanvasClient{
        http:    &http.Client{Transport: transport, Timeout: 30 * time.Second},
        baseURL: strings.TrimSuffix(baseURL, "/"),
        token:   token,
        limiter: rate.NewLimiter(10, 1),  // 10 requests/second
    }
}

func (c *CanvasClient) Do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
    if err := c.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    
    req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+c.token)
    req.Header.Set("Content-Type", "application/json")
    
    return c.http.Do(req)
}
```

### Pagination handling

```go
func ParseLinkHeader(header string) map[string]string {
    links := make(map[string]string)
    for _, part := range strings.Split(header, ",") {
        sections := strings.Split(strings.TrimSpace(part), ";")
        if len(sections) < 2 {
            continue
        }
        url := strings.Trim(sections[0], "<>")
        for _, section := range sections[1:] {
            if strings.Contains(section, `rel="`) {
                rel := strings.Trim(strings.Split(section, `"`)[1], `"`)
                links[rel] = url
            }
        }
    }
    return links
}

func (c *CanvasClient) GetAllPages(ctx context.Context, path string, result interface{}) error {
    var allItems []json.RawMessage
    nextURL := c.baseURL + path + "?per_page=100"
    
    for nextURL != "" {
        resp, err := c.doRequest(ctx, "GET", nextURL)
        if err != nil {
            return err
        }
        
        var pageItems []json.RawMessage
        if err := json.NewDecoder(resp.Body).Decode(&pageItems); err != nil {
            return err
        }
        allItems = append(allItems, pageItems...)
        
        links := ParseLinkHeader(resp.Header.Get("Link"))
        nextURL = links["next"]
    }
    
    // Unmarshal accumulated items into result
    return json.Unmarshal(mustMarshal(allItems), result)
}
```

## OAuth implementation for CLI applications

### Complete OAuth flow with PKCE

PKCE (Proof Key for Code Exchange) adds security for public clients like CLI tools:

```go
func GeneratePKCE() (verifier, challenge string, err error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", "", err
    }
    verifier = base64.RawURLEncoding.EncodeToString(b)
    
    h := sha256.Sum256([]byte(verifier))
    challenge = base64.RawURLEncoding.EncodeToString(h[:])
    return verifier, challenge, nil
}
```

### Local callback server implementation

```go
func StartOAuthFlow(ctx context.Context, cfg OAuthConfig) (*oauth2.Token, error) {
    verifier, challenge, _ := GeneratePKCE()
    state, _ := GenerateState()
    
    codeChan := make(chan string, 1)
    errChan := make(chan error, 1)
    
    // Find available port and start server on loopback only
    listener, _ := net.Listen("tcp", "127.0.0.1:0")
    port := listener.Addr().(*net.TCPAddr).Port
    listener.Close()
    
    server := &http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port)}
    redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)
    
    http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Query().Get("state") != state {
            errChan <- fmt.Errorf("state mismatch: possible CSRF attack")
            return
        }
        codeChan <- r.URL.Query().Get("code")
        fmt.Fprint(w, "<html><body><h1>Success!</h1><p>You may close this window.</p></body></html>")
    })
    
    go server.ListenAndServe()
    defer server.Shutdown(context.Background())
    
    // Build authorization URL
    authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&state=%s&code_challenge=%s&code_challenge_method=S256",
        cfg.AuthURL, cfg.ClientID, url.QueryEscape(redirectURI), state, challenge)
    
    browser.OpenURL(authURL)  // github.com/pkg/browser
    
    // Wait with timeout
    ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
    defer cancel()
    
    select {
    case code := <-codeChan:
        return exchangeCode(cfg, code, verifier, redirectURI)
    case err := <-errChan:
        return nil, err
    case <-ctx.Done():
        return nil, fmt.Errorf("authentication timed out")
    }
}
```

## Secure credential storage

### Cross-platform keychain integration

```go
import "github.com/zalando/go-keyring"

const ServiceName = "canvas-cli"

type TokenStore struct {
    Host string
}

func (s *TokenStore) Save(token *oauth2.Token) error {
    data, _ := json.Marshal(token)
    return keyring.Set(ServiceName, s.Host, string(data))
}

func (s *TokenStore) Load() (*oauth2.Token, error) {
    data, err := keyring.Get(ServiceName, s.Host)
    if err != nil {
        return nil, err
    }
    var token oauth2.Token
    return &token, json.Unmarshal([]byte(data), &token)
}
```

For systems without keychain support, fall back to encrypted files with **0600 permissions**:

```go
func saveToFile(path string, token *oauth2.Token) error {
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        return err
    }
    defer f.Close()
    return json.NewEncoder(f).Encode(token)
}
```

### Configuration management with Viper

```go
func InitConfig() {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("$HOME/.canvas-cli")
    viper.AddConfigPath(".")
    
    viper.SetEnvPrefix("CANVAS")
    viper.AutomaticEnv()
    
    viper.SetDefault("output.format", "table")
    viper.SetDefault("api.timeout", "30s")
    
    viper.ReadInConfig()  // Ignore error if config doesn't exist
}
```

Config file structure (`~/.canvas-cli/config.yaml`):
```yaml
instances:
  default: myschool
  myschool:
    url: "https://myschool.instructure.com"
    use_keyring: true
output:
  format: table  # json, yaml, table
```

## Testing strategies

### VCR-style recording with go-vcr

```go
import "github.com/dnaeon/go-vcr/v2/recorder"

func TestGetCourses(t *testing.T) {
    r, _ := recorder.New("testdata/fixtures/get_courses")
    defer r.Stop()
    
    // Redact sensitive data
    r.AddFilter(func(i *cassette.Interaction) error {
        delete(i.Request.Headers, "Authorization")
        return nil
    })
    
    client := &http.Client{Transport: r}
    canvasClient := NewClientWithHTTP("https://test.instructure.com", "fake-token", client)
    
    courses, err := canvasClient.GetCourses(context.Background())
    require.NoError(t, err)
    assert.Len(t, courses, 3)
}
```

### Mock server for integration tests

```go
func TestSubmitGrade(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, "/api/v1/courses/123/assignments/456/submissions/789", r.URL.Path)
        assert.Equal(t, "PUT", r.Method)
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "id": 789,
            "score": 85.0,
            "grade": "B",
        })
    }))
    defer server.Close()
    
    client := NewClient(server.URL, "test-token")
    submission, err := client.GradeSubmission(ctx, 123, 456, 789, 85.0)
    
    require.NoError(t, err)
    assert.Equal(t, 85.0, submission.Score)
}
```

## Cross-platform build and distribution

```makefile
VERSION := $(shell git describe --tags --always)
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/canvas-darwin-amd64 ./cmd/canvas
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/canvas-darwin-arm64 ./cmd/canvas
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/canvas-linux-amd64 ./cmd/canvas
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/canvas-windows-amd64.exe ./cmd/canvas
```

## Key implementation considerations

**Multi-instance support** is essential since users may access multiple Canvas instances (work, school, test). Store credentials keyed by hostname and allow switching active instances via `canvas config use <instance>`.

**Token refresh** should happen transparently. The `golang.org/x/oauth2` package's `TokenSource` interface handles this automatically when you use `oauth2.Config.Client()`.

**Concurrent operations** require careful rate limit management. Use a semaphore pattern limiting to 5 concurrent requests, and respect the `X-Rate-Limit-Remaining` header.

**Self-hosted Canvas instances** use identical APIs to cloud-hosted, but file upload URLs may point to local storage instead of S3. Always follow redirects completely during file operations.

The Canvas API documentation at `canvas.instructure.com/doc/api/` and the new developer portal at `developerdocs.instructure.com` provide comprehensive endpoint references. The GitHub repository `instructure/canvas-lms` contains the source code for additional implementation details.
