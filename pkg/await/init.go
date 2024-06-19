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

func init() {
	if isAwaitEnabled() {
		await.Start(await.Logger(log.Printf))
	}
}

func isAwaitEnabled() bool {
	awaitDebuggerEnabled, exists := os.LookupEnv(TILT_AWAIT_DEBUGGER_ENABLED)
	isAwaitEnabled, _ := strconv.ParseBool(awaitDebuggerEnabled)
	return exists && isAwaitEnabled
}
