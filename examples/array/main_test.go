package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestMainOutput(t *testing.T) {
	// os.Stdout is process-wide.
	got := captureOutput(t, main)
	wantPrefix := `=== Array Operations Example ===

Tags: [go json ordered]

Settings:
  1. setting1 = 100
  2. setting2 = 200

All values:
  tags: [go json ordered]
  numbers: [1 2 3 4 5]
  settings: [`
	if !strings.HasPrefix(got, wantPrefix) {
		t.Fatalf("main() output = %q, want prefix %q", got, wantPrefix)
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = write

	fn()

	if err := write.Close(); err != nil {
		t.Fatalf("stdout pipe close error = %v", err)
	}
	os.Stdout = original

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, read); err != nil {
		t.Fatalf("stdout pipe read error = %v", err)
	}
	if err := read.Close(); err != nil {
		t.Fatalf("stdout pipe read close error = %v", err)
	}
	return buf.String()
}
