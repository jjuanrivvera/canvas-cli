package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

// proxyChildEnv marks the re-executed test binary that performs the actual
// request. net/http reads HTTP_PROXY/HTTPS_PROXY once per process (sync.Once
// in ProxyFromEnvironment), so any earlier request in this package would pin
// an empty proxy config for the rest of the test run. Running the client side
// in a fresh process is the only way to observe the environment being honored.
const proxyChildEnv = "CANVAS_CLI_PROXY_TEST_CHILD"

// runProxyChild re-executes the current test binary running only testName,
// with the given environment appended. It returns the child's combined output
// and whether it exited 0.
func runProxyChild(t *testing.T, testName string, env ...string) (string, bool) {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run=^"+testName+"$", "-test.v") // #nosec G204 -- re-executing the test binary itself
	cmd.Env = append(os.Environ(), env...)
	cmd.Env = append(cmd.Env, proxyChildEnv+"=1")
	out, err := cmd.CombinedOutput()
	return string(out), err == nil
}

// recordingProxy is an httptest server that records what reaches it.
type recordingProxy struct {
	mu       sync.Mutex
	methods  []string
	hosts    []string
	absolute []bool
}

func (p *recordingProxy) record(r *http.Request) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.methods = append(p.methods, r.Method)
	p.hosts = append(p.hosts, r.Host)
	p.absolute = append(p.absolute, r.URL.IsAbs())
}

func newRecordingProxy(t *testing.T) (*httptest.Server, *recordingProxy) {
	t.Helper()
	rec := &recordingProxy{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.record(r)
		if r.Method == http.MethodConnect {
			// Refuse the tunnel: the assertion is only that CONNECT arrived here.
			http.Error(w, "tunnel refused by test proxy", http.StatusBadGateway)
			return
		}
		// Plain-HTTP proxying: the client sends an absolute-URI request line.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id": 1, "name": "via proxy"}`))
	}))
	return server, rec
}

// TestProxyFromEnvironment_HTTPProxy proves an http:// Canvas URL is routed
// through HTTP_PROXY: the proxy sees an absolute-URI request for the target
// host, and the client gets the proxy's response.
func TestProxyFromEnvironment_HTTPProxy(t *testing.T) {
	if os.Getenv(proxyChildEnv) == "1" {
		proxyChildRequest(t, "http://canvas.proxy-test.invalid")
		return
	}

	proxy, rec := newRecordingProxy(t)
	defer proxy.Close()

	out, ok := runProxyChild(t, "TestProxyFromEnvironment_HTTPProxy",
		"HTTP_PROXY="+proxy.URL, "HTTPS_PROXY="+proxy.URL, "NO_PROXY=")
	if !ok {
		t.Fatalf("child process failed:\n%s", out)
	}

	rec.mu.Lock()
	defer rec.mu.Unlock()
	if len(rec.hosts) == 0 {
		t.Fatalf("no request reached the proxy; child output:\n%s", out)
	}
	for i, host := range rec.hosts {
		if !strings.HasPrefix(host, "canvas.proxy-test.invalid") {
			t.Errorf("request %d: proxy saw Host %q, want target host", i, host)
		}
		if !rec.absolute[i] {
			t.Errorf("request %d: proxy saw a relative URL, want absolute-URI proxy form", i)
		}
	}
}

// TestProxyFromEnvironment_HTTPSProxy proves an https:// Canvas URL is routed
// through HTTPS_PROXY: the proxy receives a CONNECT for the target host. The
// test proxy refuses the tunnel, so the child only asserts that the failure
// came from the proxy rather than from a direct connection attempt.
func TestProxyFromEnvironment_HTTPSProxy(t *testing.T) {
	if os.Getenv(proxyChildEnv) == "1" {
		proxyChildRequestExpectingProxyRefusal(t, "https://canvas.proxy-test.invalid")
		return
	}

	proxy, rec := newRecordingProxy(t)
	defer proxy.Close()

	out, ok := runProxyChild(t, "TestProxyFromEnvironment_HTTPSProxy",
		"HTTP_PROXY="+proxy.URL, "HTTPS_PROXY="+proxy.URL, "NO_PROXY=")
	if !ok {
		t.Fatalf("child process failed:\n%s", out)
	}

	rec.mu.Lock()
	defer rec.mu.Unlock()
	if len(rec.methods) == 0 {
		t.Fatalf("no request reached the proxy; child output:\n%s", out)
	}
	for i, method := range rec.methods {
		if method != http.MethodConnect {
			t.Errorf("request %d: proxy saw method %s, want CONNECT", i, method)
		}
		if rec.hosts[i] != "canvas.proxy-test.invalid:443" {
			t.Errorf("request %d: CONNECT host = %q, want canvas.proxy-test.invalid:443", i, rec.hosts[i])
		}
	}
}

// proxyChildRequest runs in the child process: build a real client against an
// unresolvable host and expect the request to succeed (only possible via the proxy).
func proxyChildRequest(t *testing.T, baseURL string) {
	t.Helper()
	client, err := NewClient(ClientConfig{
		BaseURL:             baseURL,
		Token:               "test-token",
		RequestsPerSec:      100,
		RetryInitialBackoff: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	var course Course
	if err := client.GetJSON(context.Background(), "/api/v1/courses/1", &course); err != nil {
		t.Fatalf("GetJSON through proxy failed: %v", err)
	}
	if course.Name != "via proxy" {
		t.Fatalf("response did not come from the proxy: %+v", course)
	}
}

// proxyChildRequestExpectingProxyRefusal runs in the child process: the
// request must fail with the proxy's 502 on CONNECT, not with a DNS error from
// a direct attempt against the .invalid host.
func proxyChildRequestExpectingProxyRefusal(t *testing.T, baseURL string) {
	t.Helper()
	client, err := NewClient(ClientConfig{
		BaseURL:        baseURL,
		Token:          "test-token",
		RequestsPerSec: 100,
		// The proxy's 502 is retryable; keep the backoff tiny so the child
		// does not sleep through the default 1s/2s/4s schedule.
		RetryInitialBackoff: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	var course Course
	err = client.GetJSON(context.Background(), "/api/v1/courses/1", &course)
	if err == nil {
		t.Fatal("expected the refused CONNECT to surface as an error")
	}
	if !strings.Contains(err.Error(), "502") && !strings.Contains(err.Error(), "Bad Gateway") {
		t.Fatalf("error did not come from the proxy's CONNECT refusal: %v", err)
	}
}
