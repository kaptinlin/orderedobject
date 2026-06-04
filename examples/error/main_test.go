package main

import (
	"strings"
	"testing"

	"github.com/kaptinlin/orderedobject"
	"github.com/kaptinlin/orderedobject/internal/testutil"
)

func TestMainOutput(t *testing.T) {
	// os.Stdout is process-wide.
	got := testutil.CaptureOutput(t, main)
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
	got := testutil.CaptureOutput(t, func() {
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
	config := orderedobject.New[any]().Set("port", float64(8080))
	got := testutil.CaptureOutput(t, func() {
		printPort(config)
	})
	want := "Port number: 8080\n"
	if got != want {
		t.Fatalf("printPort() output = %q, want %q", got, want)
	}
}

func TestPrintPortSkipsMissingPort(t *testing.T) {
	// os.Stdout is process-wide.
	got := testutil.CaptureOutput(t, func() {
		printPort(orderedobject.New[any]())
	})
	if got != "" {
		t.Fatalf("printPort() output = %q, want empty", got)
	}
}

func TestPrintLookupReportsFoundValue(t *testing.T) {
	// os.Stdout is process-wide.
	config := orderedobject.New[any]().Set("answer", 42)
	got := testutil.CaptureOutput(t, func() {
		printLookup(config, "answer")
	})
	want := "Found value: 42\n"
	if got != want {
		t.Fatalf("printLookup() output = %q, want %q", got, want)
	}
}

func TestPrintServerPortReportsMissingServer(t *testing.T) {
	// os.Stdout is process-wide.
	got := testutil.CaptureOutput(t, func() {
		printServerPort(orderedobject.New[any]())
	})
	want := "server configuration not found\n"
	if got != want {
		t.Fatalf("printServerPort() output = %q, want %q", got, want)
	}
}

func TestPrintServerPortReportsWrongType(t *testing.T) {
	// os.Stdout is process-wide.
	config := orderedobject.New[any]().Set("server", "localhost")
	got := testutil.CaptureOutput(t, func() {
		printServerPort(config)
	})
	want := "server is not an object type\n"
	if got != want {
		t.Fatalf("printServerPort() output = %q, want %q", got, want)
	}
}

func TestPrintServerPortSkipsMissingPort(t *testing.T) {
	// os.Stdout is process-wide.
	config := orderedobject.New[any]().Set("server", orderedobject.New[any]())
	got := testutil.CaptureOutput(t, func() {
		printServerPort(config)
	})
	if got != "" {
		t.Fatalf("printServerPort() output = %q, want empty", got)
	}
}
