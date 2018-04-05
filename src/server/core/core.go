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
 * Footle's behavior.  These are acted up on.  Some other commands (e.g.
 * breakpoint_set, breakpoint_remove) need special treatment outside a
 * debugging session to allow breakpoint management at all times.
 */
func ProcessUICmds(CmdsFromUIs, DBGpCmds chan string, DBGpMessages chan message.Message, DBGpConnection *conn.Connection) {

	config := config.Get()

	for cmd := range CmdsFromUIs {

		// First, deal with Footle specific commands.
		if cmd == "on" {
			DBGpConnection.Activate()
			broadcastFakeMsgToUIs("on", "awake", DBGpMessages)
			continue
		} else if cmd == "off" {
			DBGpConnection.Deactivate()
			broadcastFakeMsgToUIs("off", "asleep", DBGpMessages)
			continue
		} else if cmd == "continue" {
			DBGpConnection.Disconnect()
			broadcastFakeMsgToUIs("continue", "stopped", DBGpMessages)
			continue
		}

		// Now the DBGp commands.
		cmdName, cmdArgs, err := command.Break(cmd)

		if nil != err {
			log.Println(err)
			continue
		}

		DBGpCmdName, err := command.Extract(cmd)

		if nil != err {
			log.Println(err)
			continue
		}

		if DBGpCmdName == "breakpoint_set" {
			// Filepaths coming from UIs *could be* relative paths.  These need to be
			// turned into absolute file URIs such as file:///foo/bar/baz.php
			cmdArgs[0] = breakpoint.ToAbsoluteUri(cmdArgs[0], config)
		}

		if DBGpCmdName == "breakpoint_set" && !DBGpConnection.IsOnAir() {
			// Example command from UI: breakpoint_set index.php 18
			filename := cmdArgs[0]
			lineNo := cmdArgs[1]
			breakpoint.Enqueue(breakpoint.Line_type_breakpoint, filename, lineNo)
			breakpoint.BroadcastPending(DBGpMessages)
		} else if DBGpCmdName == "breakpoint_remove" && !DBGpConnection.IsOnAir() {
			// Example command from UI: breakpoint_remove 18
			breakpointId := cmdArgs[0]
			breakpoint.RemovePending(breakpointId)
			breakpoint.BroadcastPending(DBGpMessages)
		} else if fullDBGpCmd, err := command.Prepare(DBGpCmdName, cmdArgs); err != nil {
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
		} else if state == "" && (msg.Properties.Command == "breakpoint_set" || msg.Properties.Command == "breakpoint_remove") {
			requestBreakpointList(DBGpCmds)
		} else if state == "" && msg.Properties.Command == "breakpoint_list" {
			breakpoint.RenewList(msg.Breakpoints)
		}

		BroadcastMsgToUIs(msg, MsgsForCmdLineUI, MsgsForHTTPUI)
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
 * Broadcast message for Footle's internal commands.
 *
 * Knowing about the execution states resulting from the internal commands
 * allows UIs to offer better UX.
 *
 * Example commands: on, off, continue.
 */
func broadcastFakeMsgToUIs(cmd string, state string, DBGpMessages chan message.Message) {

	fakeMsg := message.Message{}
	fakeMsg.MessageType = "response"
	fakeMsg.Properties.Command = cmd
	fakeMsg.State = state

	DBGpMessages <- fakeMsg
}
