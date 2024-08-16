package await

import (
	"os"
	"testing"
)

func TestIsAwaitEnabled_SetToTrue(t *testing.T) {
	os.Setenv(EnvDebuguerEnabled, "true")
	expected := true
	result := isAwaitEnabled()
	if result != expected {
		t.Errorf("isAwaitEnabled() = %v; want %v", result, expected)
	}
}

func TestIsAwaitEnabled_SetToFalse(t *testing.T) {
	os.Setenv(EnvDebuguerEnabled, "false")
	expected := false
	result := isAwaitEnabled()
	if result != expected {
		t.Errorf("isAwaitEnabled() = %v; want %v", result, expected)
	}
}

func TestIsAwaitEnabled_NotSet(t *testing.T) {
	os.Unsetenv(EnvDebuguerEnabled)
	expected := false
	result := isAwaitEnabled()
	if result != expected {
		t.Errorf("isAwaitEnabled() = %v; want %v", result, expected)
	}
}
