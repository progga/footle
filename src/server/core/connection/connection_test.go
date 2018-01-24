/**
 * Tests for network connection management.
 */

package connection

import (
	"testing"
	"time"
)

/**
 * Tests for Connection.WaitUntilActive().
 *
 * WaitUntilActive() should wait until an activation has been signalled.
 */
func TestWaitUntilActive(t *testing.T) {

	// Signal activation.  This should end wait.
	var conn Connection
	conn.wait = make(chan bool)

	activated := false

	go func() {
		time.Sleep(time.Millisecond)
		conn.signalActivation()
	}()

	go func() {
		conn.WaitUntilActive()
		activated = true
	}()

	time.Sleep(3 * time.Millisecond)
	if activated == false {
		t.Error("Activation failed.")
	}

	// Do not signal activation.  Goroutine should keep waiting.
	var conn2 Connection
	conn2.wait = make(chan bool)

	activated = false

	go func() {
		conn2.WaitUntilActive()
		activated = true
	}()

	time.Sleep(3 * time.Millisecond)
	if activated == true {
		t.Error("Failed to wait despite no activation.")
	}
}
