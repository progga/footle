/**
 * @file
 * Package for post-processing debugger commands and responses.
 */

package core

import (
	"log"
	"os"
	"server/config"
	"server/core/breakpoint"
	footlecmd "server/core/cmd"
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

	for cmd := range CmdsFromUIs {
		cmdAlias, cmdArgs, err := command.Break(cmd)
		if nil != err {
			log.Println(err)
			continue
		}

		if footlecmd.Is(cmdAlias) {
			processFootleCmds(cmdAlias, cmdArgs, DBGpMessages, DBGpConnection)
		} else if DBGpCmdName, err := command.Extract(cmd); err == nil {
			processDBGpCmds(DBGpCmdName, cmdArgs, DBGpCmds, DBGpMessages, DBGpConnection)
		} else {
			log.Println(err)
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

		broadcastMsgToUIs(msg, MsgsForCmdLineUI, MsgsForHTTPUI)
	}
}

/**
 * Post-processing of DBGp commands.
 *
 * Commands coming from the UIs may need some massaging before they are ready
 * for the consumption of the DBGp engine.  Examples:
 *   - We need to provide absolute file URIs for creating breakpoints.  UIs
 *     usually will send us relative filepaths.
 *   - Some DBGp commands are queued when the DBGp engine is unavailable.  These
 *     queued commands are later sent when the engine makes contact.
 */
func processDBGpCmds(cmdName string, cmdArgs []string, DBGpCmds chan string, DBGpMessages chan message.Message, DBGpConnection *conn.Connection) {

	if cmdName == "breakpoint_set" {
		// Filepaths coming from UIs *could be* relative paths.  These need to be
		// turned into absolute file URIs such as file:///foo/bar/baz.php
		config := config.Get()
		cmdArgs[0] = toAbsoluteUri(cmdArgs[0], config)
	}

	if cmdName == "breakpoint_set" && !DBGpConnection.IsOnAir() {
		// Example command from UI: breakpoint_set index.php 18
		filename := cmdArgs[0]
		lineNo := cmdArgs[1]

		breakpoint.Enqueue(breakpoint.Line_type_breakpoint, filename, lineNo)
		breakpoint.BroadcastPending(DBGpMessages)
	} else if cmdName == "breakpoint_remove" && !DBGpConnection.IsOnAir() {
		// Example command from UI: breakpoint_remove 18
		breakpointId := cmdArgs[0]
		breakpoint.RemovePending(breakpointId)
		breakpoint.BroadcastPending(DBGpMessages)
	} else if err := command.Validate(cmdName, cmdArgs); err == nil {
		fullDBGpCmd, _ := command.Prepare(cmdName, cmdArgs)
		DBGpCmds <- fullDBGpCmd
	}
}

/**
 * Processing of Footle's internal commands.
 */
func processFootleCmds(cmdAlias string, cmdArgs []string, DBGpMessages chan message.Message, DBGpConnection *conn.Connection) {

	if cmdAlias == "on" {
		DBGpConnection.Activate()

		fakeCmd := message.Properties{Command: "on"}
		broadcastFakeMsg(fakeCmd, "awake", DBGpMessages)
	} else if cmdAlias == "off" {
		DBGpConnection.Deactivate()

		fakeCmd := message.Properties{Command: "off"}
		broadcastFakeMsg(fakeCmd, "asleep", DBGpMessages)
	} else if cmdAlias == "continue" {
		DBGpConnection.Disconnect()

		fakeCmd := message.Properties{Command: "continue"}
		broadcastFakeMsg(fakeCmd, "stopped", DBGpMessages)
	} else if cmdAlias == "update_source" && len(cmdArgs) == 1 {
		filename := cmdArgs[0]

		config := config.Get()
		absoluteFilename := toAbsolutePath(filename, config)

		if _, err := os.Stat(absoluteFilename); os.IsNotExist(err) {
			log.Printf("File doesn't exist: %s", filename)
		}

		fakeCmd := message.Properties{Command: cmdAlias, Filename: filename}
		broadcastFakeMsg(fakeCmd, "", DBGpMessages)
	}
}

/**
 * Pass on a DBGP message to all the user interfaces.
 *
 * User interfaces include the command line interface and the HTTP interface.
 */
func broadcastMsgToUIs(msg message.Message, toCmdLine, toHTTP chan<- message.Message) {

	if nil != toCmdLine {
		toCmdLine <- msg
	}

	if nil != toHTTP {
		toHTTP <- msg
	}
}

/**
 * Broadcast response for Footle's internal commands.
 *
 * Knowing about the execution states resulting from the internal commands
 * allows UIs to offer better UX.
 *
 * Example commands: on, off, continue, update_source.
 */
func broadcastFakeMsg(prop message.Properties, state string, DBGpMessages chan message.Message) {

	fakeMsg := message.Message{}
	fakeMsg.MessageType = "response"
	fakeMsg.Properties.Command = prop.Command
	fakeMsg.Properties.Filename = prop.Filename
	fakeMsg.State = state

	DBGpMessages <- fakeMsg
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
