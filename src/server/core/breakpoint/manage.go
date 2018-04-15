/**
 * Manage list of established and pending breakpoints.
 */

package breakpoint

import (
	"log"
	"server/dbgp/message"
	"strconv"
)

/**
 * The DBGp protocol uses "enabled" to indicate that a breakpoint is currently
 * in use.
 */
const BreakpointEnabledState = "enabled"

/**
 * @see getNewId()
 */
var lastPendingBreakpointId int = 0

/**
 * List both pending and established breakpoints.
 */
func ListAllBreakpoints() (breakpoints map[int]message.Breakpoint) {

	breakpoints = make(map[int]message.Breakpoint)

	for k, breakpointRecord := range pending {
		breakpoints[k] = message.Breakpoint{
			Filename: breakpointRecord.Filename,
			LineNo:   breakpointRecord.LineNo,
			Type:     Line_type_breakpoint,
			Id:       breakpointRecord.DBGpId,
		}
	}

	for k, breakpointRecord := range established {
		breakpoints[k] = message.Breakpoint{
			Filename: breakpointRecord.Filename,
			LineNo:   breakpointRecord.LineNo,
			Type:     Line_type_breakpoint,
			Id:       breakpointRecord.DBGpId,
		}
	}

	return breakpoints
}

/**
 * Renew breakpoint list.
 *
 * Update our list of existing breakpoints maintained by the DBGp engine.
 */
func RenewList(breakpoints map[int]message.Breakpoint) {

	established.Empty()

	for _, v := range breakpoints {
		add(v.Type, v.Filename, v.LineNo, v.Id, v.State)
	}
}

/**
 * Delete the given breakpoint record from *our list*.
 */
func Delete(breakpointId int) {

	delete(established, breakpointId)
}

/**
 * Add a *pending* breakpoint record.
 *
 * Only deals with line breakpoints at the moment.
 */
func Enqueue(breakpointType, arg0, arg1 string) {

	if breakpointType == Line_type_breakpoint {
		enqueueLine(arg0, arg1)
	}
}

/**
 * Create a new breakpoint record in *our list*.
 *
 * This record is for an existing breakpoint.  Only deals with line breakpoints
 * at the moment.
 */
func add(breakpointType, filename string, lineNo, id int, state string) {

	breakpointState := (state == BreakpointEnabledState)

	if breakpointType == Line_type_breakpoint {
		established.AddLine(filename, lineNo, id, breakpointState)
	}
}

/**
 * Add a pending breakpoint record for a source code line.
 */
func enqueueLine(filename, lineNoArg string) {

	lineNo, err := strconv.Atoi(lineNoArg)

	if filename == "" || err != nil {
		log.Println(err)
		return
	}

	pendingBreakpointId := getNewId()

	b := breakpoint{
		Type:     Line_type_breakpoint,
		LineNo:   lineNo,
		Filename: filename,
		DBGpId:   pendingBreakpointId,
		State:    true,
	}

	pending.push(b)
}

/**
 * Produce a new ID number for breakpoint records.
 *
 * This Breakpoint ID is different from the numbers assigned by the DBGp engine.
 * These are for Footle's internal use as IDs for pending breakpoints.  Numbers
 * start at -1 and keeps going *down*: -2, -3,...
 */
func getNewId() int {

	lastPendingBreakpointId -= 1

	return lastPendingBreakpointId
}
