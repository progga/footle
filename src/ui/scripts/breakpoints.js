
/**
 * @file
 * Breakpoint management.
 *
 * Keep track of filenames and their breakpoints.  Update the UI to reflect
 * the status of the breakpoints.
 *
 * Only line number based breakpoints are supported at the moment.
 */

"use strict";

/**
 * List of breakpoints.
 *
 * Key: filename
 * Value: List of line numbers.
 */
var fileBreakpointMapping = {}

/**
 * Update breakpoints.
 *
 * When a file containing a new breakpoint has no tab, create a new tab for
 * that file.
 *
 * @param object msg
 */
function processBreakpoint(breakpoint) {

  var filename = breakpoint.Filename;
  var lineNo   = breakpoint.LineNumber;

  addBreakpointMapping(filename, lineNo);
  addTab(filename);
}

/**
 * Add breakpoint record.
 *
 * Update the fileBreakpointMapping global variable.
 *
 * @param string filename
 * @param int lineNumber
 * @return array
 *    Updated copy of fileBreakpointMapping
 */
function addBreakpointMapping(filename, lineNumber) {

  if (!fileBreakpointMapping.hasOwnProperty(filename)) {
    fileBreakpointMapping[filename] = [];
  }

  fileBreakpointMapping[filename].push(lineNumber);

  return fileBreakpointMapping;
}

/**
 * Update breakpoint highlighting.
 */
function refreshBreakpoints() {

  for (var filename in fileBreakpointMapping) {
    for (var lineNumberIndex in fileBreakpointMapping[filename]) {
      var lineNumber = fileBreakpointMapping[filename][lineNumberIndex];
      var tabNavElement = null;

      if (tabNavElement = hasFileTabMapping(filename)) {
        var tabIndex = tabNavElement.index();
        var tabContent = jQuery(".tab-content").get(tabIndex);

        var lineNumberClass = ".line__" + lineNumber;       
        jQuery(lineNumberClass, tabContent).addClass("breakpoint");
      }
    }
  }
}
