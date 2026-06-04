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

func TestPrintArrayConfigFallsBackWhenSettingsIsMissing(t *testing.T) {
	// os.Stdout is process-wide.
	config := orderedobject.New[any]().
		Set("tags", []string{"go"}).
		Set("numbers", []int{1, 2})

	got := testutil.CaptureOutput(t, func() {
		printArrayConfig(config)
	})
	want := `
Tags: [go]

All values:
  tags: [go]
  numbers: [1 2]
`
	if got != want {
		t.Fatalf("printArrayConfig() output = %q, want %q", got, want)
	}
}

func TestPrintArrayConfigSkipsInvalidSettings(t *testing.T) {
	// os.Stdout is process-wide.
	config := orderedobject.New[any]().
		Set("settings", []any{
			"not an object",
			orderedobject.New[any]().Set("value", 100),
			orderedobject.New[any]().Set("name", "missing value"),
			orderedobject.New[any]().Set("name", "valid").Set("value", 200),
		})

	got := testutil.CaptureOutput(t, func() {
		printArrayConfig(config)
	})
	wantPrefix := `
Settings:
  4. valid = 200

All values:
  settings: [not an object `
	if !strings.HasPrefix(got, wantPrefix) {
		t.Fatalf("printArrayConfig() output = %q, want prefix %q", got, wantPrefix)
	}
}
