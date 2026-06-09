package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- fixExamplesInMarkdown (pure function) ---

func TestFixExamplesInMarkdown_NoSynopsis(t *testing.T) {
	// Content without a Synopsis section should pass through unchanged.
	input := "# canvas\n\nSome description\n\n### Flags\n"
	got := fixExamplesInMarkdown(input)
	if got != input {
		t.Errorf("expected unchanged content, got diff:\nwant: %q\ngot:  %q", input, got)
	}
}

func TestFixExamplesInMarkdown_SynopsisNoExamples(t *testing.T) {
	input := "### Synopsis\n\nA description without examples.\n\n### Flags\n"
	got := fixExamplesInMarkdown(input)
	// Still contains Synopsis header and Flags header.
	if !strings.Contains(got, "### Synopsis") {
		t.Error("expected '### Synopsis' to remain")
	}
	if !strings.Contains(got, "### Flags") {
		t.Error("expected '### Flags' to remain")
	}
}

func TestFixExamplesInMarkdown_SynopsisWithExamples(t *testing.T) {
	input := strings.Join([]string{
		"### Synopsis",
		"",
		"Description of the command.",
		"",
		"Examples:",
		"",
		"  canvas courses list",
		"  # list all courses",
		"",
		"### Flags",
		"",
		"Some flags.",
		"",
	}, "\n")

	got := fixExamplesInMarkdown(input)

	// The output should wrap example lines in a code block.
	if !strings.Contains(got, "```bash") {
		t.Error("expected output to contain '```bash'")
	}
	if !strings.Contains(got, "canvas courses list") {
		t.Error("expected output to contain 'canvas courses list'")
	}
	if !strings.Contains(got, "# list all courses") {
		t.Error("expected output to contain '# list all courses'")
	}
	if !strings.Contains(got, "```") {
		t.Error("expected output to contain closing code fence")
	}
	// The Flags section should still be present.
	if !strings.Contains(got, "### Flags") {
		t.Error("expected '### Flags' to remain after Synopsis")
	}
}

func TestFixExamplesInMarkdown_ExamplesFlushAtEndOfSection(t *testing.T) {
	// Examples at the end of the Synopsis section (no following ### header)
	// should still be flushed.
	input := strings.Join([]string{
		"### Synopsis",
		"",
		"Description.",
		"",
		"Examples:",
		"",
		"  canvas users list",
		"",
	}, "\n")

	got := fixExamplesInMarkdown(input)

	if !strings.Contains(got, "```bash") {
		t.Errorf("expected code fence in:\n%s", got)
	}
	if !strings.Contains(got, "canvas users list") {
		t.Errorf("expected 'canvas users list' in:\n%s", got)
	}
}

func TestFixExamplesInMarkdown_MultipleExampleGroups(t *testing.T) {
	// Two example lines separated by a blank line where the next line is
	// still an example — they should remain in a single code block.
	input := strings.Join([]string{
		"### Synopsis",
		"",
		"Description.",
		"",
		"Examples:",
		"",
		"  canvas assignments list --course-id 1",
		"",
		"  canvas assignments get 2 --course-id 1",
		"",
		"### Options",
		"",
	}, "\n")

	got := fixExamplesInMarkdown(input)

	if !strings.Contains(got, "```bash") {
		t.Errorf("expected code fence in:\n%s", got)
	}
	if !strings.Contains(got, "canvas assignments list") {
		t.Error("expected first example command")
	}
	if !strings.Contains(got, "canvas assignments get") {
		t.Error("expected second example command")
	}
}

func TestFixExamplesInMarkdown_NonExampleLineEndsBlock(t *testing.T) {
	// A non-example line inside the Synopsis after Examples should end the block.
	input := strings.Join([]string{
		"### Synopsis",
		"",
		"Description.",
		"",
		"Examples:",
		"",
		"  canvas courses list",
		"This is a normal text line.",
		"",
		"### Flags",
		"",
	}, "\n")

	got := fixExamplesInMarkdown(input)

	if !strings.Contains(got, "```bash") {
		t.Errorf("expected code fence in:\n%s", got)
	}
	// The non-example line should still appear.
	if !strings.Contains(got, "This is a normal text line.") {
		t.Errorf("expected normal text line in:\n%s", got)
	}
}

func TestFixExamplesInMarkdown_Idempotent_EmptyContent(t *testing.T) {
	got := fixExamplesInMarkdown("")
	if got != "" {
		t.Errorf("expected empty string for empty input, got %q", got)
	}
}

// --- postProcessMarkdownFiles ---

func TestPostProcessMarkdownFiles_ProcessesMarkdownFiles(t *testing.T) {
	dir := t.TempDir()

	// Write a .md file that has an Examples section.
	content := strings.Join([]string{
		"### Synopsis",
		"",
		"Short description.",
		"",
		"Examples:",
		"",
		"  canvas courses list",
		"",
	}, "\n")

	mdPath := filepath.Join(dir, "canvas_courses.md")
	if err := os.WriteFile(mdPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	if err := postProcessMarkdownFiles(dir); err != nil {
		t.Fatalf("postProcessMarkdownFiles returned error: %v", err)
	}

	result, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("failed to read processed file: %v", err)
	}

	got := string(result)
	if !strings.Contains(got, "```bash") {
		t.Errorf("expected code fence after processing, got:\n%s", got)
	}
}

func TestPostProcessMarkdownFiles_SkipsIndexMd(t *testing.T) {
	dir := t.TempDir()

	// index.md should be skipped — write a sentinel that would be modified
	// by fixExamplesInMarkdown so we can detect if it was touched.
	content := strings.Join([]string{
		"### Synopsis",
		"",
		"Examples:",
		"",
		"  canvas courses list",
		"",
	}, "\n")

	indexPath := filepath.Join(dir, "index.md")
	if err := os.WriteFile(indexPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write index.md: %v", err)
	}

	if err := postProcessMarkdownFiles(dir); err != nil {
		t.Fatalf("postProcessMarkdownFiles returned error: %v", err)
	}

	result, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index.md: %v", err)
	}

	// index.md must NOT have been transformed — still has plain "canvas courses list".
	if strings.Contains(string(result), "```bash") {
		t.Error("index.md should not have been processed")
	}
}

func TestPostProcessMarkdownFiles_SkipsNonMarkdownFiles(t *testing.T) {
	dir := t.TempDir()

	// A .txt file in the same directory should be silently skipped.
	txtPath := filepath.Join(dir, "readme.txt")
	if err := os.WriteFile(txtPath, []byte("hello"), 0644); err != nil {
		t.Fatalf("failed to write readme.txt: %v", err)
	}

	if err := postProcessMarkdownFiles(dir); err != nil {
		t.Fatalf("postProcessMarkdownFiles returned unexpected error: %v", err)
	}
}

func TestPostProcessMarkdownFiles_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := postProcessMarkdownFiles(dir); err != nil {
		t.Fatalf("postProcessMarkdownFiles on empty dir returned error: %v", err)
	}
}

func TestPostProcessMarkdownFiles_ErrorOnUnreadableFile(t *testing.T) {
	// Simulate an error by passing a path that does not exist.
	err := postProcessMarkdownFiles("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("expected error for non-existent directory, got nil")
	}
}
