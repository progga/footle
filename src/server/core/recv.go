/**
 * Functions for receiving messages from DBGp engine.
 */

package core

import (
	"../dbgp"
	"../dbgp/message"
	"log"
	"net"
)

/**
 * Receive message from DBGp engine and pass it to user interfaces.
 */
func RecvMsgsFromDBGpEngine(sock net.Listener, activeDBGpConnection *net.Conn, MsgsForCmdLineUI, MsgsForHTTPUI chan<- message.Message) {

	for {
		*activeDBGpConnection = StartTalkingToDBGpEngine(sock)

		for {
			msg, err := dbgp.Read(*activeDBGpConnection)
			if len(msg) == 0 || nil != err {
				break
			}

			if parsedMsg, err := message.Decode(msg); nil == err {
				BroadcastMsgToUIs(parsedMsg, MsgsForCmdLineUI, MsgsForHTTPUI)
			}
		}

		(*activeDBGpConnection).Close()
	}
}

/**
 * Establish connection with a DBGp engine.
 */
func StartTalkingToDBGpEngine(sock net.Listener) (connection net.Conn) {

	connection, err := sock.Accept()
	if nil != err {
		log.Fatal(err)
	}

	return connection
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
