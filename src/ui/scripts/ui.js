
/**
 * @file
 * Behaviour for tabbed file browser and viewer.
 *
 * File browser is displayed in an iframe.  This iframe sits inside the first
 * tab.  When a filename is clicked in the file browser, it should be displayed
 * in another tab.  This tab will contain the HTML formatted file content.
 *
 * *** Js Syntax: ES2015 (AKA ES6) ***
 */

import * as filelist from './file-list.js'
import * as breakpoint from './breakpoints.js'
import * as breaks from './breaks.js'
import * as control from './controls.js'
import * as source from './source.js'
import * as stacktrace from './stacktrace.js'
import * as tab from './tabs.js'
import * as variable from './variables.js'

/**
 * Onload event handler.
 *
 * Does the following:
 *   - Sets up click handlers on file links in the file browser.
 *   - Sets up click handlers on tab close links.
 *   - Adds buttons for Run and Step commands.
 *   - Sets up new breakpoint trigger.
 *   - Creates a Server-sent-event handler to listen to the data stream from the
 *     Footle server.
 */
jQuery(function () {
  /* We may have missed the very first "load" event for the iframe.  So we
     need to explicitely call setupFileList() for setting up the click
     handlers on the file links inside the file browser. */
  filelist.setup()

  // iframe is loaded again when a directory is opened.
  jQuery('iframe').on('load', filelist.setup)

  filelist.setupRecent()
  tab.setupRefresher()
  tab.setupCloser()
  control.setupContinuationControls()
  control.setupStateControl()
  breakpoint.setupTrigger()
  variable.setupInteraction()
  control.disable()
  initBreakpoints()

  // Process responses from the server.
  var sse = new EventSource('/message-stream')
  jQuery(sse).on('message', function (event) {
    var msg = JSON.parse(event.originalEvent.data)
    console.log(msg)

    processMsg(msg)
  })
})

/**
 * Update UI based on debugger response.
 *
 * @param object msg
 */
function processMsg (msg) {
  if (msg.MessageType === 'response' && msg.State === 'break' && msg.Properties.Filename) {
    breaks.update(msg.Properties.Filename, msg.Properties.LineNumber)
    control.enable()
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'breakpoint_list') {
    breakpoint.refresh(msg.Breakpoints)
  } else if (msg.MessageType === 'response' && msg.State === 'stopped') {
    breaks.removePrevious()
    control.disable()
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'context_get') {
    variable.updateDisplay(msg.Context)
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'property_get') {
    variable.displaySingle(msg.Context.Local)
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'stack_get') {
    stacktrace.display(msg.Stacktrace)
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'update_source') {
    source.update(msg.Properties.Filename)
  } else if (msg.MessageType === 'response' && msg.State === 'awake') {
    control.toggleOnOffbuttons()
    breaks.removePrevious()
  } else if (msg.MessageType === 'response' && msg.State === 'asleep') {
    control.toggleOnOffbuttons()
    breaks.removePrevious()
    control.disable()
  } else if (msg.MessageType === 'init') {
  }

  updateExecutionState(msg.State)
}

/**
 * Update programme execution state.
 *
 * This is useful for providing hints in the UI about the state of the debugger
 * and the programme we are debugging.
 *
 * @param string state
 */
function updateExecutionState (state) {
  if (typeof state === 'undefined' || state === '' || state === 'waiting') {
    return
  }

  jQuery('.execution-states').attr('data-state', state)
}

/**
 * Fetch and process existing breakpoints.
 *
 * Fetch existing breakpoints from the server.  Then arrange for the breakpoints
 * to be fully processed and displayed.
 */
function initBreakpoints () {
  jQuery.getJSON('breakpoints', processMsg)
}
