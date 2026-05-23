package loader_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/cronaudit/internal/loader"
)

func writeTempCrontab(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "crontab")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestFromFile_ReadsLines(t *testing.T) {
	content := "# comment\n5 * * * * /usr/bin/backup\n0 2 * * 0 /usr/bin/cleanup\n"
	path := writeTempCrontab(t, content)

	src, err := loader.FromFile(path)
	if err != nil {
		t.Fatalf("FromFile error: %v", err)
	}
	defer loader.Close(src)

	lines, err := loader.Lines(src)
	if err != nil {
		t.Fatalf("Lines error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestFromFile_MissingFile(t *testing.T) {
	_, err := loader.FromFile("/nonexistent/path/crontab")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestFromStdin_ReturnsSource(t *testing.T) {
	src := loader.FromStdin()
	if src == nil {
		t.Fatal("expected non-nil source from FromStdin")
	}
	if src.Name != "<stdin>" {
		t.Errorf("expected name <stdin>, got %q", src.Name)
	}
}

func TestLines_EmptyFile(t *testing.T) {
	path := writeTempCrontab(t, "")
	src, _ := loader.FromFile(path)
	defer loader.Close(src)

	lines, err := loader.Lines(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(lines))
	}
}

func TestLines_NilSource(t *testing.T) {
	_, err := loader.Lines(nil)
	if err == nil {
		t.Error("expected error for nil source")
	}
	if !strings.Contains(err.Error(), "nil source") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestClose_NonCloser(t *testing.T) {
	src := loader.FromStdin()
	// os.Stdin is an io.Closer, but we just ensure no panic
	if err := loader.Close(src); err != nil {
		t.Logf("close returned: %v (acceptable)", err)
	}
}

func TestClose_NilSource(t *testing.T) {
	if err := loader.Close(nil); err != nil {
		t.Errorf("expected nil error closing nil source, got %v", err)
	}
}
