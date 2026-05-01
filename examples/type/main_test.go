package main

import (
	"testing"

	"github.com/kaptinlin/orderedobject/internal/testutil"
)

func TestMainOutput(t *testing.T) {
	// os.Stdout is process-wide.
	got := testutil.CaptureOutput(t, main)
	want := `=== Type Safety Example ===

User: Alice, Age: 30

All settings:
  default: Theme=light, Active=true
  custom: Theme=dark, Active=false
`
	if got != want {
		t.Fatalf("main() output = %q, want %q", got, want)
	}
}
