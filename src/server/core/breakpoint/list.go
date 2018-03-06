/**
 * A data structure for storing established breakpoints.
 */
package breakpoint

const Line_type_breakpoint = "line"
const Code_type_breakpoint = "code"

type breakpoint struct {
	Type     string
	State    bool
	LineNo   int
	Filename string
	Code     string
	DBGpId   int
}

type breakpointList map[int]*breakpoint

/**
 * Have we got any?
 */
func (b breakpointList) Exists() bool {

	return len(b) > 0
}

/**
 * Add a breakpoint record of type "line".
 */
func (b breakpointList) AddLine(filename string, lineNo, id int, state bool) {

	b[id] = &breakpoint{
		Type:     Line_type_breakpoint,
		State:    state,
		LineNo:   lineNo,
		Filename: filename,
		DBGpId:   id,
	}
}

/**
 * Activate the given breakpoint.
 */
func (b breakpointList) Activate(id int) (exists bool) {

	_, exists = b[id]

	if !exists {
		return false
	}

	b[id].State = true
	return true
}

/**
 * Deactivate the given breakpoint.
 */
func (b breakpointList) Deactivate(id int) (exists bool) {

	_, exists = b[id]

	if !exists {
		return false
	}

	b[id].State = false
	return true
}

/**
 * Drop all the breakpoint records.
 */
func (b breakpointList) Empty() {

	for id := range b {
		delete(b, id)
	}
}
