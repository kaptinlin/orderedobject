package main

import (
	"testing"

	"github.com/kaptinlin/orderedobject/internal/testutil"
)

func TestMainOutput(t *testing.T) {
	// os.Stdout is process-wide.
	got := testutil.CaptureOutput(t, main)
	want := `=== Basic Operations Example ===
Version: 1.0.0

Configuration:
  app_name: MyApp
  version: 1.0.1
  max_connections: 100

All entries:
  app_name: MyApp
  version: 1.0.1
  max_connections: 100

Development config:
  app_name: MyApp
  version: 1.0.1
  max_connections: 100
  debug: true
  environment: development
`
	if got != want {
		t.Fatalf("main() output = %q, want %q", got, want)
	}
}
