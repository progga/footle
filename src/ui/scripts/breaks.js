
/**
 * @file
 * Break management.
 */

import {addTab, getTabContentElementForFile, hasFileTabMapping} from './tabs'

var filenameOfLastBreak = ''
var lineNoOfLastBreak = -1

/**
 * Apply a new break.
 *
 * Remove the old break if any and then display the new break.
 *
 * @param string filename
 * @param int lineNo
 */
function updateBreak (filename, lineNo) {
  removePreviousBreak()

  displayNewBreak(filename, lineNo)
  displayFileWithNewBreak(filename)

  recordBreak(filename, lineNo)
}

/**
 * Save the latest break.
 *
 * Save it so that it can be later removed from display when another break is
 * encountered.
 *
 * @param string filename
 * @param int lineNo
 */
function recordBreak (filename, lineNo) {
  if (filename) {
    filenameOfLastBreak = filename
  }

  if (lineNo > -1) {
    lineNoOfLastBreak = lineNo
  }
}

/**
 * Display the break.
 *
 * If the file with the break is not currently open in a tab, then open
 * it now.
 *
 * @param string filename
 */
function displayFileWithNewBreak (filename) {
  addTab(filename, redrawCurrentBreak)
}

/**
 * Remove previous break...
 *
 * ...so that the latest one can be drawn.  At any point, there can be only one
 * break.  The display should reflect that.
 */
function removePreviousBreak () {
  removeBreak(filenameOfLastBreak, lineNoOfLastBreak)
}

/**
 * Redraw the break.
 *
 * This is needed just after opening a new file with a break.
 */
function redrawCurrentBreak () {
  displayNewBreak(filenameOfLastBreak, lineNoOfLastBreak)
}

/**
 * Remove break from display.
 *
 * @param string filename
 * @param int lineNo
 */
function removeBreak (filename, lineNo) {
  var tabContentForFile = getTabContentElementForFile(filename)
  var tabContentIsAbsent = (tabContentForFile === undefined)

  if (tabContentIsAbsent) {
    return
  }

  var breakSelector = '.break.line__' + lineNo
  jQuery(breakSelector, tabContentForFile).removeClass('break')
}

/**
 * Highlight the new break.
 *
 * @param string filename
 * @param int lineNo
 */
function displayNewBreak (filename, lineNo) {
  var tabContentForFile = getTabContentElementForFile(filename)
  var tabContentIsAbsent = (tabContentForFile === null)

  if (tabContentIsAbsent) {
    return
  }

  // Open the tab when it is not already open.
  var tabNavElement = hasFileTabMapping(filename)
  var isActiveTab = tabNavElement.hasClass('uk-active')
  if (!isActiveTab) {
    tabNavElement.click()
  }

  var lineNoClass = '.line__' + lineNo
  var lineElement = jQuery(lineNoClass, tabContentForFile)
  lineElement.addClass('break')

  // When the line *number* is outside the viewport, bring the line within.
  var lineNoElement = jQuery('.line__number', lineElement)
  if (!isInViewport(lineNoElement)) {
    scrollLineIntoView(lineElement)
  }
}

/**
 * Bring line into view.
 *
 * If possible, display the given line in the middle of the viewable area of
 * the browser.
 *
 * @param object lineElement
 *    jQuery object.
 */
function scrollLineIntoView (lineElement) {
  var lineY = lineElement.get(0).offsetTop
  var windowHeight = window.innerHeight
  var centreLine = lineY - windowHeight / 2

  if (centreLine > 0) {
    // element.scrollIntoView({block: 'center', behavior: 'smooth'}) does
    // just what we want.  But it is yet to be implemented in Firefox :(
    window.scrollTo(0, centreLine)
  }
}

/**
 * Is the given element currently present in the viewport?
 *
 * @param object element
 *    jQuery element
 * @return bool
 *
 * @see https://stackoverflow.com/questions/123999/how-to-tell-if-a-dom-element-is-visible-in-the-current-viewport/7557433#7557433
 */
function isInViewport (element) {
  if (element.length === 0) {
    return false
  }

  var rect = element[0].getBoundingClientRect()

  if (rect.height === 0) {
    return false
  }

  var windowWidth = window.innerWidth
  var windowHeight = window.innerHeight

  var isInViewportHorizontally = (rect.left >= 0) && (rect.right <= windowWidth)
  var isInViewportVertically = (rect.top >= 0) && (rect.bottom <= windowHeight)
  var isInViewport = isInViewportHorizontally && isInViewportVertically

  return isInViewport
}

export {updateBreak, removePreviousBreak}
