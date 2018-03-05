/**
 * Interact with DBGp engine.
 */

package breakpoint

import (
	"server/config"
	"server/dbgp/command"
	"server/dbgp/message"
	"strconv"
)

var list breakpointList = make(breakpointList)
var pending Queue

/**
 * Send breakpoint creation commands for queued breakpoints.
 *
 * When Footle is not connected to the DBGp engine, new breakpoints coming from
 * the UI are queued.  These are sent to the DBGp engine when the next debugging
 * session starts.
 */
func SendPending(DBGpCmds chan string) {

	// As well as pending breakpoints, breakpoints from the previous session have
	// to be set again.
	for _, v := range list {
		pending.push(*v)
	}

	for len(pending) > 0 {
		breakpointRecord := pending.pop()

		lineNoArg := strconv.Itoa(breakpointRecord.LineNo)
		cmdArgs := []string{breakpointRecord.Filename, lineNoArg}
		cmd, err := command.Prepare("breakpoint_set", cmdArgs)

		if err != nil {
			continue
		}

		DBGpCmds <- cmd
	}
}

/**
 * Broadcast the list of existing and pending breakpoints.
 *
 * Existing breakpoints are the ones that have been set during the previous
 * debugging session.  Pending breakpoints have been added through the UI, but
 * have not been sent to the debugging engine yet.
 */
func BroadcastPending(DBGpMessages chan message.Message, config config.Config) {

	fakeMsg := FakeMessage{}
	fakeMsg.init(config, "breakpoint_list")
	fakeMsg.AddExistingBreakpoints(list)
	fakeMsg.AddPendingBreakpoints(pending)

	msg := fakeMsg.GetMsg()
	DBGpMessages <- msg
}
