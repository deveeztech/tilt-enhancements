package await

import (
	"os"
	"testing"
)

func TestIsAwaitEnabledSetToTrue(t *testing.T) {
	os.Setenv(EnvDebuguerEnabled, "true")
	expected := true
	result := isAwaitEnabled()
	if result != expected {
		t.Errorf("isAwaitEnabled() = %v; want %v", result, expected)
	}
}

func TestIsAwaitEnabledSetToFalse(t *testing.T) {
	os.Setenv(EnvDebuguerEnabled, "false")
	expected := false
	result := isAwaitEnabled()
	if result != expected {
		t.Errorf("isAwaitEnabled() = %v; want %v", result, expected)
	}
}

func TestIsAwaitEnabledNotSet(t *testing.T) {
	os.Unsetenv(EnvDebuguerEnabled)
	expected := false
	result := isAwaitEnabled()
	if result != expected {
		t.Errorf("isAwaitEnabled() = %v; want %v", result, expected)
	}
}
