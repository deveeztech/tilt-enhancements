package await

import (
	"log"
	"os"
	"strconv"

	"github.com/deveeztech/tilt-enhancements/internal/await"
)

const (
	TILT_AWAIT_DEBUGGER_ENABLED = "TILT_AWAIT_DEBUGGER_ENABLED"
)

// It checks if the TILT_AWAIT_DEBUGGER_ENABLED environment variable is set to true,
// and if so, it starts the await functionality with the specified logger.
func init() {
	if isAwaitEnabled() {
		await.Start(await.Logger(log.Printf))
	}
}

// isAwaitEnabled is a helper function that checks if the TILT_AWAIT_DEBUGGER_ENABLED
// environment variable is set to true.
// It returns true if the variable is set to true, and false otherwise.
func isAwaitEnabled() bool {
	awaitDebuggerEnabled, exists := os.LookupEnv(TILT_AWAIT_DEBUGGER_ENABLED)
	isAwaitEnabled, _ := strconv.ParseBool(awaitDebuggerEnabled)
	return exists && isAwaitEnabled
}
