/**
 * Manage list of established and pending breakpoints.
 */

package breakpoint

import (
	"log"
	"server/dbgp/message"
	"strconv"
	"strings"
)

/**
 * The DBGp protocol uses "enabled" to indicate that a breakpoint is currently
 * in use.
 */
const BreakpointEnabledState = "enabled"

/**
 * Renew breakpoint list.
 *
 * Update our list of existing breakpoints maintained by the DBGp engine.
 */
func RenewList(breakpoints map[int]message.Breakpoint) {

	list.Empty()

	for _, v := range breakpoints {
		add(v.Type, v.Filename, v.LineNo, v.State)
	}
}

/**
 * Delete the given breakpoint record from *our list*.
 */
func Delete(breakpointId int) {

	delete(list, breakpointId)
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
func add(breakpointType, filename string, lineNo int, state string) {

	breakpointState := (state == BreakpointEnabledState)

	if breakpointType == Line_type_breakpoint {
		list.AddLine(filename, lineNo, breakpointState)
	}
}

/**
 * Add a pending breakpoint record for a source code line.
 */
func enqueueLine(filename, lineNoArg string) {

	lineNoArg = strings.Trim(lineNoArg, "\x00")
	lineNo, err := strconv.Atoi(lineNoArg)

	if filename == "" || err != nil {
		log.Println(err)
		return
	}

	b := breakpoint{
		Type:     Line_type_breakpoint,
		LineNo:   lineNo,
		Filename: filename,
		DBGpId:   -1,
		State:    true,
	}

	pending.push(b)
}
