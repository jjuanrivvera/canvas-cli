package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// buildTestStateManager creates a StateManager backed by a temp file so tests
// don't touch ~/.canvas-cli.
func buildTestStateManager(t *testing.T) *StateManager {
	t.Helper()
	return newStateManagerWithPath(filepath.Join(t.TempDir(), "state.json"))
}

// buildAutoUpdaterWithState constructs an AutoUpdater wired to a specific
// StateManager and Updater so network calls can be mocked.
func buildAutoUpdaterWithState(version string, enabled bool, sm *StateManager, u *Updater) *AutoUpdater {
	return &AutoUpdater{
		CurrentVersion: version,
		CheckInterval:  time.Hour,
		Enabled:        enabled,
		updater:        u,
		stateManager:   sm,
	}
}

// --- NewAutoUpdater ---

func TestNewAutoUpdater_DefaultInterval(t *testing.T) {
	au, err := NewAutoUpdater("1.0.0", true, 0)
	if err != nil {
		t.Fatalf("NewAutoUpdater failed: %v", err)
	}
	if au == nil {
		t.Fatal("expected non-nil AutoUpdater")
	}
	if au.CheckInterval != DefaultCheckInterval {
		t.Errorf("expected default interval %v, got %v", DefaultCheckInterval, au.CheckInterval)
	}
}

func TestNewAutoUpdater_CustomInterval(t *testing.T) {
	au, err := NewAutoUpdater("1.0.0", false, 30*time.Minute)
	if err != nil {
		t.Fatalf("NewAutoUpdater failed: %v", err)
	}
	if au.CheckInterval != 30*time.Minute {
		t.Errorf("expected 30m interval, got %v", au.CheckInterval)
	}
	if au.Enabled {
		t.Error("expected Enabled=false")
	}
}

// --- RunUpdateCheck branches ---

func TestAutoUpdater_RunUpdateCheck_Disabled(t *testing.T) {
	sm := buildTestStateManager(t)
	u := NewUpdater("1.0.0")
	au := buildAutoUpdaterWithState("1.0.0", false, sm, u)

	au.RunUpdateCheck(context.Background())
}

func TestAutoUpdater_RunUpdateCheck_DevVersion(t *testing.T) {
	sm := buildTestStateManager(t)
	u := NewUpdater("dev")
	au := buildAutoUpdaterWithState("dev", true, sm, u)

	au.RunUpdateCheck(context.Background())
}

func TestAutoUpdater_RunUpdateCheck_EmptyVersion(t *testing.T) {
	sm := buildTestStateManager(t)
	u := NewUpdater("")
	au := buildAutoUpdaterWithState("", true, sm, u)

	au.RunUpdateCheck(context.Background())
}

func TestAutoUpdater_RunUpdateCheck_IntervalNotElapsed(t *testing.T) {
	sm := buildTestStateManager(t)
	// Save a very recent check so ShouldCheck returns false.
	if err := sm.Save(&State{LastCheckTime: time.Now()}); err != nil {
		t.Fatal(err)
	}

	u := NewUpdater("1.0.0")
	au := buildAutoUpdaterWithState("1.0.0", true, sm, u)
	au.CheckInterval = time.Hour

	au.RunUpdateCheck(context.Background())
}

func TestAutoUpdater_RunUpdateCheck_CurrentVersionSame(t *testing.T) {
	release := Release{
		TagName: "v1.0.0",
		Assets:  []Asset{},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	sm := buildTestStateManager(t)
	u := NewUpdater("1.0.0")
	u.HTTPClient = &http.Client{
		Transport: &urlRewriteTransport{targetURL: server.URL},
	}

	au := buildAutoUpdaterWithState("1.0.0", true, sm, u)

	au.RunUpdateCheck(context.Background())

	state, err := sm.Load()
	if err != nil {
		t.Fatal(err)
	}
	if state.LastCheckTime.IsZero() {
		t.Error("expected LastCheckTime to be set after check")
	}
}

// --- RunUpdateCheckAsync / WaitForCompletion ---

func TestAutoUpdater_RunUpdateCheckAsync_Disabled(t *testing.T) {
	sm := buildTestStateManager(t)
	u := NewUpdater("1.0.0")
	au := buildAutoUpdaterWithState("1.0.0", false, sm, u)

	au.RunUpdateCheckAsync(context.Background())
	completed := au.WaitForCompletion(5 * time.Second)
	if !completed {
		t.Error("expected async check to complete within 5s")
	}
}

func TestAutoUpdater_WaitForCompletion_NoDoneChannel(t *testing.T) {
	sm := buildTestStateManager(t)
	u := NewUpdater("1.0.0")
	au := buildAutoUpdaterWithState("1.0.0", false, sm, u)
	// done is nil; WaitForCompletion should return true immediately.
	completed := au.WaitForCompletion(time.Second)
	if !completed {
		t.Error("expected true when no async operation was started")
	}
}

func TestAutoUpdater_WaitForCompletion_Timeout(t *testing.T) {
	sm := buildTestStateManager(t)
	u := NewUpdater("1.0.0")
	au := buildAutoUpdaterWithState("1.0.0", false, sm, u)

	// Manually create a done channel that never closes.
	au.done = make(chan struct{})

	completed := au.WaitForCompletion(10 * time.Millisecond)
	if completed {
		t.Error("expected false (timeout) when done channel never closes")
	}
}

func TestAutoUpdater_RunUpdateCheckAsync_CompletesNaturally(t *testing.T) {
	sm := buildTestStateManager(t)
	u := NewUpdater("dev") // dev version skips immediately
	au := buildAutoUpdaterWithState("dev", true, sm, u)

	au.RunUpdateCheckAsync(context.Background())
	completed := au.WaitForCompletion(5 * time.Second)
	if !completed {
		t.Error("expected async check to complete quickly for dev version")
	}
}

// --- PrintNotifications ---

func TestAutoUpdater_PrintNotifications_NoPending(t *testing.T) {
	sm := buildTestStateManager(t)
	u := NewUpdater("1.0.0")
	au := buildAutoUpdaterWithState("1.0.0", true, sm, u)

	au.PrintNotifications()
}

func TestAutoUpdater_PrintNotifications_RecentUpdate(t *testing.T) {
	sm := buildTestStateManager(t)
	if err := sm.Save(&State{
		LastUpdateTime:   time.Now(),
		LastVersion:      "1.0.0",
		UpdatedToVersion: "1.1.0",
	}); err != nil {
		t.Fatal(err)
	}

	u := NewUpdater("1.1.0")
	au := buildAutoUpdaterWithState("1.1.0", true, sm, u)

	au.PrintNotifications()
}

func TestAutoUpdater_PrintNotifications_RecentError(t *testing.T) {
	sm := buildTestStateManager(t)
	if err := sm.Save(&State{
		LastError:     "rate limit exceeded",
		LastErrorTime: time.Now(),
	}); err != nil {
		t.Fatal(err)
	}

	u := NewUpdater("1.0.0")
	au := buildAutoUpdaterWithState("1.0.0", true, sm, u)

	au.PrintNotifications()
}

// --- CheckNow ---

func TestAutoUpdater_CheckNow_DevVersion(t *testing.T) {
	sm := buildTestStateManager(t)
	u := NewUpdater("dev")
	au := buildAutoUpdaterWithState("dev", true, sm, u)

	result := au.CheckNow(context.Background())
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Updated {
		t.Error("dev version should never update")
	}
}

func TestAutoUpdater_CheckNow_NetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	sm := buildTestStateManager(t)
	u := NewUpdater("1.0.0")
	u.HTTPClient = &http.Client{
		Transport: &urlRewriteTransport{targetURL: server.URL},
	}

	au := buildAutoUpdaterWithState("1.0.0", true, sm, u)

	result := au.CheckNow(context.Background())
	if result.Error == nil {
		t.Error("expected error when server returns 500")
	}
}

// --- GetState ---

func TestAutoUpdater_GetState_NoFile(t *testing.T) {
	sm := buildTestStateManager(t)
	u := NewUpdater("1.0.0")
	au := buildAutoUpdaterWithState("1.0.0", true, sm, u)

	state, err := au.GetState()
	if err != nil {
		t.Fatalf("GetState failed: %v", err)
	}
	if state == nil {
		t.Fatal("expected non-nil state")
	}
}

func TestAutoUpdater_GetState_WithData(t *testing.T) {
	sm := buildTestStateManager(t)
	if err := sm.Save(&State{LastVersion: "1.2.3"}); err != nil {
		t.Fatal(err)
	}

	u := NewUpdater("1.2.3")
	au := buildAutoUpdaterWithState("1.2.3", true, sm, u)

	state, err := au.GetState()
	if err != nil {
		t.Fatalf("GetState failed: %v", err)
	}
	if state.LastVersion != "1.2.3" {
		t.Errorf("expected LastVersion 1.2.3, got %s", state.LastVersion)
	}
}

// --- Updater.CheckAndUpdate full-path tests ---

func TestUpdater_CheckAndUpdate_NewerVersionAvailable(t *testing.T) {
	archName := runtime.GOARCH
	if archName == "amd64" {
		archName = "x86_64"
	}
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	archiveName := fmt.Sprintf("canvas-cli_%s_%s%s", runtime.GOOS, archName, ext)

	var archive []byte
	if runtime.GOOS == "windows" {
		archive, _ = buildZipArchive(t)
	} else {
		archive = buildTarGzArchive(t, []byte("new binary"))
	}

	release := Release{
		TagName: "v9.9.9",
		Assets: []Asset{
			{Name: archiveName, BrowserDownloadURL: "/download/binary"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/download/binary":
			w.WriteHeader(http.StatusOK)
			w.Write(archive)
		default:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(release)
		}
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	fakeExec := filepath.Join(tmpDir, "canvas")
	if err := os.WriteFile(fakeExec, []byte("old binary"), 0755); err != nil {
		t.Fatal(err)
	}

	u := NewUpdater("1.0.0")
	u.ExecutablePath = fakeExec
	u.HTTPClient = &http.Client{
		Transport: &urlRewriteTransport{targetURL: server.URL},
	}

	result := u.CheckAndUpdate(context.Background())
	if result.Error != nil {
		t.Fatalf("CheckAndUpdate failed: %v", result.Error)
	}
	if !result.Updated {
		t.Error("expected Updated=true")
	}
	if result.ToVersion != "9.9.9" {
		t.Errorf("expected ToVersion 9.9.9, got %s", result.ToVersion)
	}
}

func TestUpdater_CheckAndUpdate_NoCompatibleAsset(t *testing.T) {
	release := Release{
		TagName: "v9.9.9",
		Assets:  []Asset{},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	u := NewUpdater("1.0.0")
	u.HTTPClient = &http.Client{
		Transport: &urlRewriteTransport{targetURL: server.URL},
	}

	result := u.CheckAndUpdate(context.Background())
	if result.Error == nil {
		t.Error("expected error when no compatible asset found")
	}
}

func TestUpdater_CheckAndUpdate_DownloadFails(t *testing.T) {
	archName := runtime.GOARCH
	if archName == "amd64" {
		archName = "x86_64"
	}
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	archiveName := fmt.Sprintf("canvas-cli_%s_%s%s", runtime.GOOS, archName, ext)

	release := Release{
		TagName: "v9.9.9",
		Assets:  []Asset{{Name: archiveName, BrowserDownloadURL: "/fail"}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/fail" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	u := NewUpdater("1.0.0")
	u.HTTPClient = &http.Client{
		Transport: &urlRewriteTransport{targetURL: server.URL},
	}

	result := u.CheckAndUpdate(context.Background())
	if result.Error == nil {
		t.Error("expected error when download fails")
	}
}

func TestUpdater_CheckAndUpdate_ChecksumVerificationFails(t *testing.T) {
	archName := runtime.GOARCH
	if archName == "amd64" {
		archName = "x86_64"
	}
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	archiveName := fmt.Sprintf("canvas-cli_%s_%s%s", runtime.GOOS, archName, ext)

	archive := buildTarGzArchive(t, []byte("new binary"))
	// Wrong checksum ensures verification fails.
	wrongChecksums := fmt.Sprintf("deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef  %s\n", archiveName)

	release := Release{
		TagName: "v9.9.9",
		Assets: []Asset{
			{Name: archiveName, BrowserDownloadURL: "/binary"},
			{Name: "checksums.txt", BrowserDownloadURL: "/checksums"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/binary":
			w.Write(archive)
		case "/checksums":
			w.Write([]byte(wrongChecksums))
		default:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(release)
		}
	}))
	defer server.Close()

	u := NewUpdater("1.0.0")
	u.HTTPClient = &http.Client{
		Transport: &urlRewriteTransport{targetURL: server.URL},
	}

	result := u.CheckAndUpdate(context.Background())
	if result.Error == nil {
		t.Error("expected error when checksum verification fails")
	}
}

func TestUpdater_CheckAndUpdate_ChecksumDownloadFails(t *testing.T) {
	archName := runtime.GOARCH
	if archName == "amd64" {
		archName = "x86_64"
	}
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	archiveName := fmt.Sprintf("canvas-cli_%s_%s%s", runtime.GOOS, archName, ext)

	archive := buildTarGzArchive(t, []byte("new binary"))

	release := Release{
		TagName: "v9.9.9",
		Assets: []Asset{
			{Name: archiveName, BrowserDownloadURL: "/binary"},
			{Name: "checksums.txt", BrowserDownloadURL: "/checksums-fail"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/binary":
			w.Write(archive)
		case "/checksums-fail":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(release)
		}
	}))
	defer server.Close()

	u := NewUpdater("1.0.0")
	u.HTTPClient = &http.Client{
		Transport: &urlRewriteTransport{targetURL: server.URL},
	}

	result := u.CheckAndUpdate(context.Background())
	if result.Error == nil {
		t.Error("expected error when checksums download fails")
	}
}

// --- downloadAsset ---

func TestUpdater_DownloadAsset_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("file content"))
	}))
	defer server.Close()

	u := NewUpdater("1.0.0")
	u.HTTPClient = server.Client()

	asset := &Asset{BrowserDownloadURL: server.URL + "/file"}
	data, err := u.downloadAsset(context.Background(), asset)
	if err != nil {
		t.Fatalf("downloadAsset failed: %v", err)
	}
	if string(data) != "file content" {
		t.Errorf("unexpected content: %q", data)
	}
}

func TestUpdater_DownloadAsset_Non200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	u := NewUpdater("1.0.0")
	u.HTTPClient = server.Client()

	asset := &Asset{BrowserDownloadURL: server.URL + "/missing"}
	_, err := u.downloadAsset(context.Background(), asset)
	if err == nil {
		t.Error("expected error for 404 response")
	}
}

// --- downloadChecksums ---

func TestUpdater_DownloadChecksums(t *testing.T) {
	checksumContent := "abc123  file1.tar.gz\ndef456  file2.tar.gz\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(checksumContent))
	}))
	defer server.Close()

	u := NewUpdater("1.0.0")
	u.HTTPClient = server.Client()

	asset := &Asset{BrowserDownloadURL: server.URL + "/checksums.txt"}
	checksums, err := u.downloadChecksums(context.Background(), asset)
	if err != nil {
		t.Fatalf("downloadChecksums failed: %v", err)
	}
	if checksums["file1.tar.gz"] != "abc123" {
		t.Errorf("expected abc123 for file1.tar.gz, got %q", checksums["file1.tar.gz"])
	}
	if checksums["file2.tar.gz"] != "def456" {
		t.Errorf("expected def456 for file2.tar.gz, got %q", checksums["file2.tar.gz"])
	}
}

func TestUpdater_DownloadChecksums_EmptyLines(t *testing.T) {
	checksumContent := "abc123  file1.tar.gz\n\n\ndef456  file2.tar.gz\n"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(checksumContent))
	}))
	defer server.Close()

	u := NewUpdater("1.0.0")
	u.HTTPClient = server.Client()

	asset := &Asset{BrowserDownloadURL: server.URL + "/checksums.txt"}
	checksums, err := u.downloadChecksums(context.Background(), asset)
	if err != nil {
		t.Fatalf("downloadChecksums failed: %v", err)
	}
	if len(checksums) != 2 {
		t.Errorf("expected 2 entries, got %d", len(checksums))
	}
}

// --- extractBinary routing ---

func TestUpdater_ExtractBinary_TarGz(t *testing.T) {
	archive := buildTarGzArchive(t, []byte("binary content"))

	u := NewUpdater("1.0.0")
	data, err := u.extractBinary(archive, "canvas-cli_linux_x86_64.tar.gz")
	if err != nil {
		t.Fatalf("extractBinary failed: %v", err)
	}
	if string(data) != "binary content" {
		t.Errorf("unexpected content: %q", data)
	}
}

func TestUpdater_ExtractBinary_ZipDispatch(t *testing.T) {
	archive, content := buildZipArchive(t)

	u := NewUpdater("1.0.0")
	data, err := u.extractBinary(archive, "canvas-cli_windows_x86_64.zip")
	if err != nil {
		t.Fatalf("extractBinary (.zip) failed: %v", err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("content mismatch: got %q, want %q", data, content)
	}
}

func TestUpdater_ExtractFromZip_Success(t *testing.T) {
	archive, content := buildZipArchive(t)

	u := NewUpdater("1.0.0")
	data, err := u.extractFromZip(archive)
	if err != nil {
		t.Fatalf("extractFromZip failed: %v", err)
	}
	if !bytes.Equal(data, content) {
		t.Errorf("content mismatch: got %q, want %q", data, content)
	}
}

func TestUpdater_ExtractFromZip_MissingBinary(t *testing.T) {
	// Build a zip that contains an unrelated file only.
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, err := zw.Create("not-the-binary.txt")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("unrelated"))
	zw.Close()

	u := NewUpdater("1.0.0")
	_, err = u.extractFromZip(buf.Bytes())
	if err == nil {
		t.Error("expected error when binary not found in zip")
	}
}

func TestUpdater_ExtractFromZip_InvalidData(t *testing.T) {
	u := NewUpdater("1.0.0")
	_, err := u.extractFromZip([]byte("not a zip"))
	if err == nil {
		t.Error("expected error for invalid zip data")
	}
}

func TestUpdater_ExtractFromTarGz_MissingBinary(t *testing.T) {
	// Build a tar.gz that does NOT contain the expected binary name.
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	content := []byte("irrelevant")
	hdr := &tar.Header{
		Name:     "other-file.txt",
		Mode:     0644,
		Size:     int64(len(content)),
		Typeflag: tar.TypeReg,
	}
	tw.WriteHeader(hdr)
	tw.Write(content)
	tw.Close()
	gz.Close()

	u := NewUpdater("1.0.0")
	_, err := u.extractFromTarGz(buf.Bytes())
	if err == nil {
		t.Error("expected error when binary not found in tar.gz")
	}
}

func TestUpdater_ExtractFromTarGz_InvalidGzip(t *testing.T) {
	u := NewUpdater("1.0.0")
	_, err := u.extractFromTarGz([]byte("this is not a gzip file"))
	if err == nil {
		t.Error("expected error for invalid gzip data")
	}
}

// --- GetLatestRelease error paths ---

func TestUpdater_GetLatestRelease_Non200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	u := NewUpdater("1.0.0")
	u.HTTPClient = &http.Client{
		Transport: &urlRewriteTransport{targetURL: server.URL},
	}

	_, err := u.GetLatestRelease(context.Background())
	if err == nil {
		t.Error("expected error for 403 response")
	}
}

// --- helpers ---

// buildTarGzArchive creates a valid tar.gz containing the platform binary.
func buildTarGzArchive(t *testing.T, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	binaryName := BinaryName
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	hdr := &tar.Header{
		Name:     binaryName,
		Mode:     0755,
		Size:     int64(len(content)),
		Typeflag: tar.TypeReg,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}
	tw.Close()
	gz.Close()
	return buf.Bytes()
}

// buildZipArchive creates a zip archive containing the platform binary.
// Returns (archiveBytes, binaryContent).
func buildZipArchive(t *testing.T) ([]byte, []byte) {
	t.Helper()
	binaryName := BinaryName
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	content := []byte("zip binary content")
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, err := zw.Create(binaryName)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(content); err != nil {
		t.Fatal(err)
	}
	zw.Close()
	return buf.Bytes(), content
}
