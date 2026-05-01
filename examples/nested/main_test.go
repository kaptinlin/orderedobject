package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestMainOutput(t *testing.T) {
	// os.Stdout is process-wide.
	got := captureOutput(t, main)
	want := `=== Nested Structures Example ===

App name: MyApp

Full configuration:
app: 
  name: MyApp
  version: 1.0.0
  debug: true
server: 
  host: localhost
  port: 8080
  ssl: 
    enabled: true
    cert: /path/to/cert.pem
database: 
  driver: postgres
  host: db.example.com
  port: 5432
  credentials: 
    username: admin
    password: secret
`
	if got != want {
		t.Fatalf("main() output = %q, want %q", got, want)
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
