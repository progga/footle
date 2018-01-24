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
 * Receive message from DBGp engine and pass it to user interfaces.
 */
func RecvMsgsFromDBGpEngine(DBGpConnection *conn.Connection, MsgsForCmdLineUI, MsgsForHTTPUI chan<- message.Message) {

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
				BroadcastMsgToUIs(parsedMsg, MsgsForCmdLineUI, MsgsForHTTPUI)
			}

			if config.IsVerbose() {
				log.Println(msg)
			}
		}

		DBGpConnection.Disconnect()
	}
}

/**
 * Pass on a DBGP message to all the user interfaces.
 *
 * User interfaces include the command line interface and the HTTP interface.
 */
func BroadcastMsgToUIs(msg message.Message, toCmdLine, toHTTP chan<- message.Message) {

	if nil != toCmdLine {
		toCmdLine <- msg
	}

	if nil != toHTTP {
		toHTTP <- msg
	}
}
