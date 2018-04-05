/**
 * Class for preparing a fake response Message.
 *
 * This fake response message resembles those passed down to UIs based on
 * responses from the DBGp engine.
 *
 * This message carries breakpoint related info.
 */

package breakpoint

import (
	"server/dbgp/message"
)

const fake_debugger_state_waiting = "waiting"

type FakeMessage struct {
	msg message.Message
}

/**
 * Initialized message.
 *
 * The debugger state is set to "waiting" which is an indication that we are
 * acting outside a debugging session.  This "waiting" state is not part of the
 * DBGp protocol.
 */
func (m *FakeMessage) init(cmd string) {

	m.msg.Breakpoints = make(map[int]message.Breakpoint)

	m.msg.MessageType = "response"
	m.msg.Properties = message.Properties{Command: cmd}
	m.msg.State = fake_debugger_state_waiting
}

/**
 * Getter for the message record.
 */
func (m *FakeMessage) GetMsg() (msg message.Message) {

	return m.msg
}

/**
 * Add pending breakpoints to the response message.
 */
func (m *FakeMessage) AddPendingBreakpoints(pending Queue) {

	for k, breakpointRecord := range pending {
		m.msg.Breakpoints[k] = message.Breakpoint{
			Filename: breakpointRecord.Filename,
			LineNo:   breakpointRecord.LineNo,
			Type:     Line_type_breakpoint,
			Id:       breakpointRecord.DBGpId,
		}
	}
}

/**
 * Add existing breakpoints to the response message.
 */
func (m *FakeMessage) AddExistingBreakpoints(existingList breakpointList) {

	for k, breakpointRecord := range existingList {
		m.msg.Breakpoints[k] = message.Breakpoint{
			Filename: breakpointRecord.Filename,
			LineNo:   breakpointRecord.LineNo,
			Type:     Line_type_breakpoint,
			Id:       breakpointRecord.DBGpId,
		}
	}
}
