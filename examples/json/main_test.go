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
		`=== JSON Operations Example ===

1. Using ToJSON:
{"id":1001,"name":"John Doe","email":"john@example.com","active":true}

2. Using json.Marshal:
{"id":1001,"name":"John Doe","email":"john@example.com","active":true}

3. Via map (order not preserved):
`,
		`"id":1001`,
		`"name":"John Doe"`,
		`"email":"john@example.com"`,
		`"active":true`,
		`
Parsed JSON:
  name: Alice
  age: 30
  skills: [Go Python]
`,
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
