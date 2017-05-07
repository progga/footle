
/**
 * Provides the HTTP user interface for the debugger.
 *
 * *** Unimplemented. ***
 */

package http

import (
  "../dbgp/message"
)

/**
 *
 */
func RunUI(out chan<- string) {
}

/**
 *
 */
func UpdateUIStatus(in <-chan message.Message) {

  for msg := range in {
    _ = msg
  }
}
