
/**
 * @file
 * Breakpoint management.
 *
 * Keep track of filenames and their breakpoints.  Update the UI to reflect
 * the status of the breakpoints.
 *
 * Only line number based breakpoints are supported at the moment.
 */

'use strict'

/**
 * List of breakpoints.
 *
 * Key: filename
 * Value: List of line numbers.
 */
var fileBreakpointMapping = {}

/**
 * Click handler for creating new breakpoints.
 *
 * When a line number is clicked, a breakpoint is added for that line.
 */
function setupBreakpointTrigger () {
  jQuery('.tab').on('click', '.tab-content', function (event) {
    // event.currentTarget is the tab element for a file.
    var filepath = jQuery(event.currentTarget).data('filepath')

    if (event.target.classList.contains('line__number')) {
      // event.target is the line number element which received the click.
      var lineNumber = event.target.innerText

      sendCommand('breakpoint_set', [filepath, lineNumber])
    }
  })
}

/**
 * Highlight new ones, remove deleted ones.
 *
 * @todo Remove deleted breakpoints.
 *
 * @param array breakpointList
 *    List of breakpoint objects.
 */
function refreshBreakpoints (breakpointList) {
  for (var breakpointId in breakpointList) {
    var breakpoint = breakpointList[breakpointId]

    if (isNewBreakpoint(breakpoint)) {
      addBreakpoint(breakpoint)
    }
  }

  highlightBreakpoints()
}

/**
 * Is this breakpoint already in our list?
 *
 * @param object breakpoint
 *    Properties: Filename, LineNo.
 * @return bool
 */
function isNewBreakpoint (breakpoint) {
  var filename = breakpoint.Filename
  var lineNo = breakpoint.LineNo

  var unknownFilename = !fileBreakpointMapping.hasOwnProperty(filename)
  if (unknownFilename) {
    return true
  }

  var lineNoIndex = fileBreakpointMapping[filename].indexOf(lineNo)
  var unknownLineNo = lineNoIndex > -1

  if (unknownLineNo) {
    return false
  }

  return true
}

/**
 * Update breakpoints.
 *
 * When a file containing a new breakpoint has no tab, create a new tab for
 * that file.
 *
 * @param object breakpoint
 */
function addBreakpoint (breakpoint) {
  var filename = breakpoint.Filename
  var lineNo = breakpoint.LineNo

  addBreakpointMapping(filename, lineNo)
  addTab(filename, highlightBreakpoints)
}

/**
 * Remove a breakpoint record and its highlighting.
 *
 * @param object breakpoint
 */
function removeBreakpoint (breakpoint) {
  var filename = breakpoint.Filename
  var lineNo = breakpoint.LineNo

  removeBreakpointMapping(filename, lineNo)
  removeBreakpointHighlighting(filename, lineNo)
}

/**
 * Remove breakpoint record.
 *
 * @param string filename
 * @param int lineNo
 */
function removeBreakpointMapping (filename, lineNo) {
  if (!fileBreakpointMapping.hasOwnProperty(filename)) {
    return
  }

  var lineNoIndex = fileBreakpointMapping[filename].indexOf(lineNo)

  if (lineNoIndex > -1) {
    fileBreakpointMapping[filename].splice(lineNoIndex, 1)
  }
}

/**
 * Add breakpoint record.
 *
 * Update the fileBreakpointMapping global variable.
 *
 * @param string filename
 * @param int lineNo
 * @return array
 *    Updated copy of fileBreakpointMapping
 */
function addBreakpointMapping (filename, lineNo) {
  if (!fileBreakpointMapping.hasOwnProperty(filename)) {
    fileBreakpointMapping[filename] = []
  }

  fileBreakpointMapping[filename].push(lineNo)

  return fileBreakpointMapping
}

/**
 * Update breakpoint highlighting.
 */
function highlightBreakpoints () {
  for (var filename in fileBreakpointMapping) {
    for (var lineNoIndex in fileBreakpointMapping[filename]) {
      var lineNo = fileBreakpointMapping[filename][lineNoIndex]

      highlightABreakpoint(filename, lineNo)
    }
  }
}

/**
 * Highlight a breakpoint.
 *
 * @param string filename
 * @param int lineNo
 */
function highlightABreakpoint (filename, lineNo) {
  var tabNavElement = hasFileTabMapping(filename)

  if (!tabNavElement) {
    return
  }

  var tabContent = getTabContentElement(tabNavElement)

  var lineNoClass = '.line__' + lineNo
  jQuery(lineNoClass, tabContent).addClass('breakpoint')
}

/**
 * Remove highlight for a breakpoint.
 *
 * @param string filename
 * @param int lineNo
 */
function removeBreakpointHighlighting (filename, lineNo) {
  var tabNavElement = hasFileTabMapping(filename)

  if (!tabNavElement) {
    return
  }
  var tabContent = getTabContentElement(tabNavElement)

  var lineNoClass = '.line__' + lineNo
  jQuery(lineNoClass, tabContent).removeClass('breakpoint')
}
