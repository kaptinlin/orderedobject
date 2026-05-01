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

	required := []string{
		"=== Map Operations Example ===\n",
		"\nFrom map (order preserved):\n",
		"\nBack to map (order not preserved):\n",
		"  theme: dark\n",
		"  font_size: 14\n",
		"  notifications: true\n",
		"  language: en\n",
	}
	for _, want := range required {
		if !strings.Contains(got, want) {
			t.Fatalf("main() output = %q, want to contain %q", got, want)
		}
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
