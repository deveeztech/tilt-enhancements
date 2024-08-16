package await

import (
	"log"
	"os"
	"strconv"

	"github.com/deveeztech/tilt-enhancements/internal/await"
)

const (
	// EnvDebuguerEnabled is the environment variable to enable the await debugger
	EnvDebuguerEnabled = "TILT_AWAIT_DEBUGGER_ENABLED"
)

// It checks if the EnvDebuguerEnabled environment variable is set to true,
// and if so, it starts the await functionality with the specified logger.
func init() {
	if isAwaitEnabled() {
		await.Start(await.Logger(log.Printf))
	}
}

// isAwaitEnabled is a helper function that checks if the EnvDebuguerEnabled
// environment variable is set to true.
// It returns true if the variable is set to true, and false otherwise.
func isAwaitEnabled() bool {
	awaitDebuggerEnabled, exists := os.LookupEnv(EnvDebuguerEnabled)
	isAwaitEnabled, _ := strconv.ParseBool(awaitDebuggerEnabled)
	return exists && isAwaitEnabled
}
