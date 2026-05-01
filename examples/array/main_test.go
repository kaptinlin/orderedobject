package main

import (
	"strings"
	"testing"

	"github.com/kaptinlin/orderedobject/internal/testutil"
)

func TestMainOutput(t *testing.T) {
	// os.Stdout is process-wide.
	got := testutil.CaptureOutput(t, main)
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
