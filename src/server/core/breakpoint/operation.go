/**
 * Interact with DBGp engine.
 */

package breakpoint

import (
	"server/dbgp/command"
	"server/dbgp/message"
	"strconv"
)

var established breakpointList = make(breakpointList)
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
	for _, v := range established {
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
 * Remove given breakpoint even if it is pending.
 */
func RemovePending(breakpointId string) (err error) {

	breakpointIdNum, err := strconv.Atoi(breakpointId)

	if err != nil {
		return err
	}

	// Because pending breakpoints are always assigned a negative Id.
	// @see getNewId()
	isPending := breakpointIdNum < 0

	if isPending {
		for breakpointIndex, breakpoint := range pending {
			if breakpoint.DBGpId == breakpointIdNum {
				pending.delete(breakpointIndex)
			}
		}
	} else if _, exists := established[breakpointIdNum]; exists {
		delete(established, breakpointIdNum)
	}

	return err
}

/**
 * Broadcast the list of existing and pending breakpoints.
 *
 * Existing breakpoints are the ones that have been set during the previous
 * debugging session.  Pending breakpoints have been added through the UI, but
 * have not been sent to the debugger engine yet.
 */
func BroadcastPending(DBGpMessages chan message.Message) {

	fakeMsg := FakeMessage{}
	fakeMsg.init("breakpoint_list")
	fakeMsg.AddExistingBreakpoints(established)
	fakeMsg.AddPendingBreakpoints(pending)

	msg := fakeMsg.GetMsg()
	DBGpMessages <- msg
}
