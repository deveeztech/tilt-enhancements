package await

import (
	"os"
	"testing"
)

func TestIsAwaitEnabled_SetToTrue(t *testing.T) {
	// Test when TILT_AWAIT_DEBUGGER_ENABLED is set to true
	os.Setenv(TILT_AWAIT_DEBUGGER_ENABLED, "true")
	expected := true
	result := isAwaitEnabled()
	if result != expected {
		t.Errorf("isAwaitEnabled() = %v; want %v", result, expected)
	}
}

func TestIsAwaitEnabled_SetToFalse(t *testing.T) {
	// Test when TILT_AWAIT_DEBUGGER_ENABLED is set to false
	os.Setenv(TILT_AWAIT_DEBUGGER_ENABLED, "false")
	expected := false
	result := isAwaitEnabled()
	if result != expected {
		t.Errorf("isAwaitEnabled() = %v; want %v", result, expected)
	}
}

func TestIsAwaitEnabled_NotSet(t *testing.T) {
	// Test when TILT_AWAIT_DEBUGGER_ENABLED is not set
	os.Unsetenv(TILT_AWAIT_DEBUGGER_ENABLED)
	expected := false
	result := isAwaitEnabled()
	if result != expected {
		t.Errorf("isAwaitEnabled() = %v; want %v", result, expected)
	}
}
