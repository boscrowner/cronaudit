package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempCrontab(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "crontab-*.txt")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRun_MutuallyExclusiveFlags(t *testing.T) {
	os.Args = []string{"cronaudit", "-invalid", "-valid"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if err := run(); err == nil {
		t.Error("expected error for mutually exclusive flags, got nil")
	}
}

func TestRun_MissingFile(t *testing.T) {
	os.Args = []string{"cronaudit", "/nonexistent/path/crontab"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if err := run(); err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRun_ValidFile_TextFormat(t *testing.T) {
	content := "# daily backup\n0 2 * * * /usr/bin/backup\n"
	path := writeTempCrontab(t, content)

	os.Args = []string{"cronaudit", "-format", "text", path}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if err := run(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_ValidFile_JSONFormat(t *testing.T) {
	content := "*/5 * * * * /usr/bin/check\n"
	path := writeTempCrontab(t, content)

	os.Args = []string{"cronaudit", "-format", "json", path}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if err := run(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRun_InvalidOnlyFlag(t *testing.T) {
	content := "bad entry\n0 2 * * * /usr/bin/ok\n"
	path := writeTempCrontab(t, content)

	os.Args = []string{"cronaudit", "-invalid", path}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if err := run(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestBinaryExists(t *testing.T) {
	binPath := filepath.Join("..", "..", "bin", "cronaudit")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Skip("binary not built, skipping integration test")
	}
	out, err := exec.Command(binPath, "-format", "text", "/dev/null").CombinedOutput()
	if err != nil {
		t.Errorf("binary returned error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(string(out), "Source") {
		t.Errorf("expected output to contain 'Source', got: %s", out)
	}
}
