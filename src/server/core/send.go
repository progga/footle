
/**
 * Functions for sending DBGp commands to DBGp engine.
 */

package core

import (
  "net"
  "log"
  "fmt"
)

/**
 * Send DBGp command to DBGp engine (i.e. Xdebug).
 */
func SendCmdsToDBGpEngine(conn *net.Conn, in <-chan string) {

  for DBGpCmd := range in {
    connection := *conn

    if isActiveConnection(connection) {
      writeCount, err := connection.Write([]byte(DBGpCmd))
      _ = writeCount

      if nil != err {
        log.Fatal(err)
      }
    } else {
      fmt.Println("Inactive connection.")
    }
  }
}

/**
 * Has the given network connection been initialized?
 *
 * Initialization happens when a DBGp engine connects to the debugger.
 */
func isActiveConnection(connection net.Conn) bool {

  ignore := []byte{}

  if nil == connection {
    return false
  }

  if readCount, err := connection.Write(ignore); nil != err {
    _ = readCount
    return false
  }

  return true
}
