/**
 * Functions for sending DBGp commands to DBGp engine.
 */

package core

import (
	"log"
	"server/config"
	conn "server/core/connection"
)

/**
 * Send DBGp command to DBGp engine (e.g. Xdebug).
 */
func SendCmdsToDBGpEngine(DBGpConnection *conn.Connection, in <-chan string) {

	config := config.Get()

	for DBGpCmd := range in {
		connection := DBGpConnection.Get()

		if DBGpConnection.IsOnAir() {
			if config.IsVerbose() {
				log.Println(DBGpCmd)
			}

			_, err := (*connection).Write([]byte(DBGpCmd))

			if nil != err {
				log.Fatal(err)
			}
		} else {
			log.Println("Cannot speak to an inactive connection.")
		}
	}
}
