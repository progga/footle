/**
 * Footle's current state.
 *
 * Current state = execution state + breakpoints.
 *
 * Possible execution states: awake, asleep, break, stopped.  These states are
 * entered into due to messages from the DBGP engine and commands from the UIs.
 *
 * Please note that Footle's current state is different from the DBGp engine's
 * state.  Footle may be sleeping while the DBGp engine could still be active.
 *
 * The current state is needed during UI initialization.
 */

package currentstate

import (
	"server/core/breakpoint"
	"server/dbgp/message"
)

/**
 * Last message that changed execution state of Footle.
 */
var lastMsg message.Message

/**
 * Fetch the last execution state and existing breakpoints.
 *
 * These are represented in the form of messages for UIs.
 */
func Get() (stateMessages []message.Message) {

	stateMessages = []message.Message{}

	if isRelevant(lastMsg) {
		stateMessages = append(stateMessages, lastMsg)
	}

	breakpointListingMsg := breakpoint.PrepareFakeMsg()
	if len(breakpointListingMsg.Breakpoints) > 0 {
		stateMessages = append(stateMessages, breakpointListingMsg)
	}

	return stateMessages
}

/**
 * Save the last message that changed execution state of Footle.
 */
func SaveLastMsg(msg message.Message) {

	if isRelevant(msg) {
		lastMsg = msg
	}
}

/**
 * Not all *messages* change the execution state.
 *
 * Some messages may have an empty "State" property.  Some states may be just
 * placeholders (e.g. waiting).  These do not change the execution state.
 * But the rest does and we filter them in here.
 */
func isRelevant(msg message.Message) bool {

	switch msg.State {
	case "init", "break", "stopped", "awake", "asleep":
		return true
	}

	return false
}
