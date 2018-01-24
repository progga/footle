/**
 * Functions for sending DBGp commands to DBGp engine.
 */

package core

import (
	"log"
	"net"
	"server/config"
	conn "server/core/connection"
)

/**
 * Send DBGp command to DBGp engine (i.e. Xdebug).
 */
func SendCmdsToDBGpEngine(DBGpConnection *conn.Connection, in <-chan string) {

	config := config.Get()

	for DBGpCmd := range in {
		connection := DBGpConnection.Get()

		if isActiveConnection(*connection) {
			if config.IsVerbose() {
				log.Println(DBGpCmd)
			}

			_, err := (*connection).Write([]byte(DBGpCmd))

			if nil != err {
				log.Fatal(err)
			}
		} else {
			log.Println("Inactive connection.")
		}
	}
}

/**
 * Has the given network connection been initialized?
 *
 * Initialization happens when a DBGp engine connects to Footle.
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
