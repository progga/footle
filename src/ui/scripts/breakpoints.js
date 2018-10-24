
/**
 * @file
 * Breakpoint management.
 *
 * Keep track of breakpoint Ids and their associated filepaths and lineNos.
 * Update the UI to reflect the current status of the breakpoints.
 *
 * Only line number based breakpoints are supported at the moment.
 */

import * as tab from './tabs.js'
import * as server from './server-commands.js'

/**
 * List of breakpoints.
 *
 * This helps to maintain association between breakpoint Ids, filepaths,
 * and line numbers.
 *
 * Key: breakpointId
 * Value: Map; key: filepath, lineNo
 */
var existingBreakpointList = new Map()

/**
 * Click handler for creating/removing breakpoints.
 *
 * Clicking a line number toggles its breakpoint.  Using event delegation, we
 * setup one click handler per tab to process clicks on any line of the file
 * shown in that tab.
 *
 * Note that when a breakpoint is present, the "breakpoint" class is assigned to
 * the *parent* of the ".line__number" element.  The breakpoint Id is also
 * stored as a data attribute of this parent using an attribute name of
 * "breakpoint-id".
 */
function setupTrigger () {
  jQuery('.tab').on('click', '.tab-content', function (event) {
    var hasClickedLineNoWOBreakpoint = event.target.classList.contains('line__number') && !jQuery(event.target).parent('.line.breakpoint').length
    var hasClickedLineNoWBreakpoint = event.target.classList.contains('line__number') && jQuery(event.target).parent('.line.breakpoint').length
    var hasNothingToDoWBreakpoint = !(hasClickedLineNoWOBreakpoint || hasClickedLineNoWBreakpoint)

    if (hasNothingToDoWBreakpoint) {
      return
    }

    // event.currentTarget is the tab element displaying a file's content.
    var filepath = jQuery(event.currentTarget).data('filepath')
    // event.target is the line number element which received the click.
    var lineNo = event.target.innerText
    var breakpointId = jQuery(event.target).parent('.line.breakpoint').data('breakpoint-id')

    if (hasClickedLineNoWOBreakpoint) {
      server.sendCommand('breakpoint_set', [filepath, lineNo])
    } else if (hasClickedLineNoWBreakpoint && breakpointId) {
      server.sendCommand('breakpoint_remove', [breakpointId])
    }
  })
}

/**
 * Highlight new ones, remove deleted ones.
 *
 * @param array newBreakpointList
 *    List of breakpoint objects containing filepath and lineNo.
 */
function refresh (newBreakpointList) {
  removeDeletedBreakpoints(newBreakpointList)
  addNewBreakpoints(newBreakpointList)
}

/**
 * Add newly created breakpoints.
 *
 * @param array newBreakpointList
 *    List of breakpoint objects containing filepath and lineNo.
 *
 * @see existingBreakpointList
 */
function addNewBreakpoints (newBreakpointList) {
  for (const breakpointIndex in newBreakpointList) {
    var breakpoint = newBreakpointList[breakpointIndex]
    var breakpointId = breakpoint.Id
    var isNewBreakpoint = !existingBreakpointList.has(breakpointId)

    if (isNewBreakpoint) {
      addBreakpoint(breakpoint)
    }
  }
}

/**
 * Remove newly deleted breakpoints.
 *
 * @param array newBreakpointList
 *    List of breakpoint objects containing filepath and lineNo.
 *
 * @see existingBreakpointList
 */
function removeDeletedBreakpoints (newBreakpointList) {
  for (const [existingBreakpointId, breakpointDetails] of existingBreakpointList) {
    var isRemoved = isRemovedBreakpoint(existingBreakpointId, newBreakpointList)

    if (isRemoved) {
      var filepath = breakpointDetails.filepath
      var lineNo = breakpointDetails.lineNo

      removeBreakpoint(filepath, lineNo, existingBreakpointId)
    }
  }
}

/**
 * Is the given breakpoint Id absent from the updated list?
 *
 * If a breakpoint Id is absent from the fresh list then that means it has been
 * removed.
 *
 * @param int existingBreakpointId
 * @param array newBreakpointList
 *    List of breakpoint objects containing filepath and lineNo.
 * @return bool
 */
function isRemovedBreakpoint (existingBreakpointId, newBreakpointList) {
  for (var breakpointIndex in newBreakpointList) {
    var breakpoint = newBreakpointList[breakpointIndex]
    var alsoExistsInNewList = breakpoint.Id === existingBreakpointId

    if (alsoExistsInNewList) {
      return false
    }
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
  var filepath = breakpoint.Filename
  var lineNo = breakpoint.LineNo
  var breakpointId = breakpoint.Id

  addBreakpointMapping(filepath, lineNo, breakpointId)
  tab.add(filepath, () => highlightBreakpoint(filepath, lineNo, breakpointId))
}

/**
 * Remove a breakpoint record and its highlighting.
 *
 * @param string filepath
 * @param int lineNo
 * @param int breakpointId
 */
function removeBreakpoint (filepath, lineNo, breakpointId) {
  removeMapping(breakpointId)
  removeHighlighting(filepath, lineNo)
}

/**
 * Add breakpoint record.
 *
 * Update the existingBreakpointList global variable.
 *
 * @param string filepath
 * @param int lineNo
 * @param int breakpointId
 */
function addBreakpointMapping (filepath, lineNo, breakpointId) {
  if (existingBreakpointList.has(breakpointId)) {
    return
  }

  existingBreakpointList.set(breakpointId, {
    'filepath': filepath,
    'lineNo': lineNo
  })
}

/**
 * Remove breakpoint record.
 *
 * @param int breakpointId
 *
 * @see existingBreakpointList
 */
function removeMapping (breakpointId) {
  existingBreakpointList.delete(breakpointId)
}

/**
 * Update breakpoint highlighting for the given file.
 *
 * @param string targetFilepath
 *    This is the file where we want to redraw the breakpoints.
 */
function highlightFile (targetFilepath) {
  for (const [breakpointId, breakpointDetails] of existingBreakpointList) {
    var filepath = breakpointDetails.filepath
    var lineNo = breakpointDetails.lineNo

    if (filepath === targetFilepath) {
      highlightBreakpoint(filepath, lineNo, breakpointId)
    }
  }
}

/**
 * Highlight a breakpoint.
 *
 * Also, save the breakpoint Id as a data attribute of the highlighted element.
 *
 * @param string filepath
 * @param int lineNo
 * @param int breakpointId
 */
function highlightBreakpoint (filepath, lineNo, breakpointId) {
  var tabNavElement = tab.hasFileMapping(filepath)

  if (!tabNavElement) {
    return
  }

  var tabContent = tab.getContentElement(tabNavElement)

  var lineNoClass = '.line__' + lineNo
  jQuery(lineNoClass, tabContent).addClass('breakpoint').data('breakpoint-id', breakpointId)
}

/**
 * Remove highlight for a breakpoint.
 *
 * Also remove the breakpoint-id data attribute from the highlighted element.
 *
 * @param string filepath
 * @param int lineNo
 */
function removeHighlighting (filepath, lineNo) {
  var tabNavElement = tab.hasFileMapping(filepath)

  if (!tabNavElement) {
    return
  }
  var tabContent = tab.getContentElement(tabNavElement)

  var lineNoClass = '.line__' + lineNo
  jQuery(lineNoClass, tabContent).removeClass('breakpoint').removeData('breakpoint-id')
}

export {highlightFile, setupTrigger, refresh}
