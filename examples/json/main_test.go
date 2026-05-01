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
