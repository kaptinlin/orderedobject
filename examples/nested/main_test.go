package main

import (
	"testing"

	"github.com/kaptinlin/orderedobject/internal/testutil"
)

func TestMainOutput(t *testing.T) {
	// os.Stdout is process-wide.
	got := testutil.CaptureOutput(t, main)
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
