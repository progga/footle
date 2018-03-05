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
	"path/filepath"
	"server/config"
	"server/dbgp/message"
)

type FakeMessage struct {
	msg    message.Message
	config config.Config
}

/**
 * Initialized message.
 *
 * The debugger state is set to "waiting" which is an indication that we are
 * acting outside a debugging session.  This "waiting" state is not part of the
 * DBGp protocol.
 */
func (m *FakeMessage) init(c config.Config, cmd string) {

	m.config = c
	m.msg.Breakpoints = make(map[int]message.Breakpoint)

	m.msg.MessageType = "response"
	m.msg.Properties = message.Properties{Command: cmd}
	m.msg.State = "waiting"
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
		filepath := m.toAbsolutePath(breakpointRecord.Filename)

		m.msg.Breakpoints[k] = message.Breakpoint{
			Filename: filepath,
			LineNo:   breakpointRecord.LineNo,
			Type:     Line_type_breakpoint,
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
		}
	}
}

/**
 * Convert a relative filepath into an absolute file URI.
 *
 * Our fake responses initially contain relative filepaths which we here turn
 * into absolute paths.  This is needed because the DBGp engine always returns
 * absolute file URIs such as "file:///home/foo/bar/baz.php".  The HTTP UI,
 * for example, expects such absolute file URIs as part of responses.
 */
func (m *FakeMessage) toAbsolutePath(relativePath string) (fullpathUri string) {

	codeDir := m.config.DetermineCodeDir()

	fullpath := filepath.Join(codeDir, relativePath)
	fullpathUri = "file://" + fullpath

	return fullpathUri
}
