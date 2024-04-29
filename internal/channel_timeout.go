package internal

import (
	"fmt"
	"time"
)

func WaitForSignalOrTimeout(sigChan <-chan bool, d time.Duration) error {
LOOP:
	for {
		select {
		case c := <-sigChan:
			Logger.Info("loop is finished", "signal", c)

			break LOOP
		case <-time.After(d):
			return fmt.Errorf("Timeout after delaying %v", d)
		}
	}

	return nil
}
