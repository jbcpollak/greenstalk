package internal

import (
	"fmt"
	"time"
)

// Wait until the provided signal is returned or the timeout is reached
// TODO: replace duration with a context closing
func WaitForSignalOrTimeout(sigChan <-chan bool, d time.Duration) (bool, error) {

	select {
	case c := <-sigChan:
		Logger.Info("loop is finished", "signal", c)

		return c, nil
	case <-time.After(d):
		return false, fmt.Errorf("timeout after delaying %v", d)
	}
}
