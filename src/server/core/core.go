
/**
 * Package for talking to the DBGp engine.
 */

package core

import (
  "log"
  "net"
)

/**
 * Start listening on the standard DBGp port of 9000.
 */
func ListenForDBGpEngine() (sock net.Listener) {

  sock, err := net.Listen("tcp", "127.0.0.1:9000");
  if  nil != err {
    log.Fatal(err)
  }

  return sock
}
