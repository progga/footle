/**
 * Functions for receiving messages from DBGp engine.
 */

package core

import (
	"log"
	"server/config"
	conn "server/core/connection"
	"server/dbgp"
	"server/dbgp/message"
)

/**
 * Receive messages from DBGp engine and send it for further processing.
 */
func RecvMsgsFromDBGpEngine(DBGpConnection *conn.Connection, DBGpMessages chan<- message.Message) {

	config := config.Get()

	for {
		DBGpConnection.WaitUntilActive()

		activeDBGpConnection := DBGpConnection.Connect()

		if *activeDBGpConnection == nil {
			continue
		}

		for {
			msg, err := dbgp.Read(*activeDBGpConnection)
			if len(msg) == 0 || nil != err {
				break
			}

			if parsedMsg, err := message.Decode(msg); nil == err {
				DBGpMessages <- parsedMsg
			}

			if config.IsVerbose() {
				log.Println(msg)
			}
		}

		DBGpConnection.Disconnect()
	}
}
