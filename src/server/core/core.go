/**
 * Package for talking to the DBGp engine.
 */

package core

import (
	"log"
	"server/config"
	"server/core/breakpoint"
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
func ProcessUICmds(CmdsFromUIs, DBGpCmds chan string, DBGpMessages chan message.Message, DBGpConnection *conn.Connection) {

	config := config.Get()

	for fullDBGpCmd := range CmdsFromUIs {
		cmdName, cmdArgs, err := command.Break(fullDBGpCmd)

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
		} else if cmdName == "breakpoint_set" && !DBGpConnection.IsOnAir() {
			// Example of cmd: breakpoint_set -i 5 -t line -f index.php -n 18\x00
			breakpoint.Enqueue(breakpoint.Line_type_breakpoint, cmdArgs[5], cmdArgs[7])
			breakpoint.BroadcastPending(DBGpMessages, config)
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
			breakpoint.SendPending(DBGpCmds)
			proceedWithSession(DBGpCmds)
		} else if state == "" && msg.Properties.Command == "breakpoint_set" {
			requestBreakpointList(DBGpCmds)
		} else if state == "" && msg.Properties.Command == "breakpoint_list" {
			breakpoint.RenewList(msg.Breakpoints)
		}

		BroadcastMsgToUIs(msg, MsgsForCmdLineUI, MsgsForHTTPUI)
	}
}

/**
 * Respond to the "stopping" state.
 *
 * End the debugging session by issuing the DBGp "stop" command.
 */
func endSession(DBGpCmds chan string) {

	stopCmd, err := command.Prepare("stop", []string{})

	if err != nil {
		return
	}

	DBGpCmds <- stopCmd
}

/**
 * Respond to the "starting" state.
 *
 * Carry on with the debugging session by issuing the DBGp "run" command.
 */
func proceedWithSession(DBGpCmds chan string) {

	runCmd, err := command.Prepare("run", []string{})

	if err != nil {
		return
	}

	DBGpCmds <- runCmd
}

/**
 * Ask the DBGp engine for its breakpoint list.
 *
 * Respond to "breakpoint_set" command by requesting the complete breakpoint
 * list.
 */
func requestBreakpointList(DBGpCmds chan string) {

	runCmd, err := command.Prepare("breakpoint_list", []string{})

	if err != nil {
		return
	}

	DBGpCmds <- runCmd
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
