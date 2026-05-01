package main

import (
	"strings"
	"testing"

	"github.com/kaptinlin/orderedobject/internal/testutil"
)

func TestMainOutput(t *testing.T) {
	// os.Stdout is process-wide.
	got := testutil.CaptureOutput(t, main)

	required := []string{
		"=== Map Operations Example ===\n",
		"\nFrom map (order follows Go map iteration):\n",
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
