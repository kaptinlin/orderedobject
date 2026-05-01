package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/kaptinlin/orderedobject"
)

func TestMainOutput(t *testing.T) {
	// os.Stdout is process-wide.
	got := captureOutput(t, main)
	want := `=== Error Handling Example ===

Successfully parsed JSON
Port type error: string
Key 'nonexistent' does not exist

Server port: 8080
`
	if got != want {
		t.Fatalf("main() output = %q, want %q", got, want)
	}
}

func TestPrintConfigReportsParseError(t *testing.T) {
	// os.Stdout is process-wide.
	var config *orderedobject.Object[any]
	got := captureOutput(t, func() {
		config = printConfig(`[`)
	})
	if !strings.HasPrefix(got, "\nParse error: ") {
		t.Fatalf("printConfig() output = %q, want parse error", got)
	}
	if got := config.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}
}

func TestPrintPortReportsNumber(t *testing.T) {
	// os.Stdout is process-wide.
	config := orderedobject.NewObject[any]().Set("port", float64(8080))
	got := captureOutput(t, func() {
		printPort(config)
	})
	want := "Port number: 8080\n"
	if got != want {
		t.Fatalf("printPort() output = %q, want %q", got, want)
	}
}

func TestPrintLookupReportsFoundValue(t *testing.T) {
	// os.Stdout is process-wide.
	config := orderedobject.NewObject[any]().Set("answer", 42)
	got := captureOutput(t, func() {
		printLookup(config, "answer")
	})
	want := "Found value: 42\n"
	if got != want {
		t.Fatalf("printLookup() output = %q, want %q", got, want)
	}
}

func TestPrintServerPortReportsMissingServer(t *testing.T) {
	// os.Stdout is process-wide.
	got := captureOutput(t, func() {
		printServerPort(orderedobject.NewObject[any]())
	})
	want := "server configuration not found\n"
	if got != want {
		t.Fatalf("printServerPort() output = %q, want %q", got, want)
	}
}

func TestPrintServerPortReportsWrongType(t *testing.T) {
	// os.Stdout is process-wide.
	config := orderedobject.NewObject[any]().Set("server", "localhost")
	got := captureOutput(t, func() {
		printServerPort(config)
	})
	want := "server is not an object type\n"
	if got != want {
		t.Fatalf("printServerPort() output = %q, want %q", got, want)
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
