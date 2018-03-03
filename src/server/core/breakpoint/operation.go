/**
 * Interact with DBGp engine.
 */

package breakpoint

import (
	"server/dbgp/command"
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
