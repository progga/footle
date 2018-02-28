/**
 * Package for talking to the DBGp engine.
 */

package core

import (
	"log"
	conn "server/core/connection"
	"server/dbgp/command"
	"server/dbgp/message"
)

/**
 * Process commands coming from UIs.
 *
 * Some commands (e.g. run) are meant for the DBGp engine.  These are forwarded
 * to the appropriate channel.  Other commands (e.g. on) are meant to control
 * Footle's behavior.  These are acted up on.
 */
func ProcessUICmds(CmdsFromUIs, DBGpCmds chan string, DBGpConnection *conn.Connection) {

	for fullDBGpCmd := range CmdsFromUIs {
		cmdName, _, err := command.Break(fullDBGpCmd)

		if nil != err {
			log.Println(err)
			continue
		}

		if cmdName == "on" {
			DBGpConnection.Activate()
		} else if cmdName == "off" {
			DBGpConnection.Deactivate()
		} else if cmdName == "continue" {
			DBGpConnection.Disconnect()
		} else {
			DBGpCmds <- fullDBGpCmd
		}
	}
}

/**
 * Act on messages coming from the DBGp engine.
 *
 * Some messages need automated responses.  These are acted up on.  Some other
 * messages are going to affect the state of the UIs.  These are broadcast to
 * the UIs.
 */
func ProcessDBGpMessages(DBGpCmds chan string, DBGpMessages, MsgsForCmdLineUI, MsgsForHTTPUI chan message.Message) {

	for msg := range DBGpMessages {
		state := msg.State

		if state == "stopping" {
			endSession(DBGpCmds)
		} else if state == "starting" {
			sendPendingBreakpoints(DBGpCmds)
		}

		BroadcastMsgToUIs(msg, MsgsForCmdLineUI, MsgsForHTTPUI)
	}
}

/**
 * Respond to the "stopping" state.
 *
 * End the debugging session by issuing the the DBGp "stop" command.
 */
func endSession(DBGpCmds chan string) {

	stopCmd, err := command.Prepare("stop", []string{})

	if err != nil {
		return
	}

	DBGpCmds <- stopCmd
}

/**
 * Setup breakpoints at the beginning of the debugging session.
 *
 * This function is a placeholder.  It does *not* do anything at the moment.
 */
func sendPendingBreakpoints(DBGpCmds chan string) {
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
