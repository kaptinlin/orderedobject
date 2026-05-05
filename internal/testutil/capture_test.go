package testutil

import (
	"fmt"
	"testing"
)

func TestCaptureOutputReturnsStdout(t *testing.T) {
	// os.Stdout is process-wide.
	got := CaptureOutput(t, func() {
		fmt.Println("hello")
	})
	want := "hello\n"
	if got != want {
		t.Fatalf("CaptureOutput() = %q, want %q", got, want)
	}
}
