
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

'use strict'

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
  /* We have missed the very first "load" event for the iframe.  So we
     need to explicitely call setupFileList() for setting up the click
     handlers on the file links inside the file browser. */
  setupFileList()

  // iframe is loaded again when a directory is opened.
  jQuery('iframe').on('load', setupFileList)

  setupTabCloser()
  setupContinuationControls()
  setupStateControl()
  setupBreakpointTrigger()
  setupVariableInteraction()

  // Process responses from the server.
  var sse = new EventSource('/message-stream')
  jQuery(sse).on('message', function (event) {
    var msg = JSON.parse(event.originalEvent.data)
    console.log(msg)

    processMsg(msg)
  })
})

/**
 * Update UI based on debugging status.
 *
 * @param object msg
 */
function processMsg (msg) {
  if (msg.MessageType === 'response' && msg.State === 'break' && msg.Properties.Filename) {
    updateBreak(msg.Properties.Filename, msg.Properties.LineNumber)
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'breakpoint_list') {
    refreshBreakpoints(msg.Breakpoints)
  } else if (msg.MessageType === 'response' && msg.State === 'stopped') {
    removePreviousBreak()
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'context_get') {
    updateVarsDisplay(msg.Context.Local)
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'property_get') {
    displaySingleVar(msg.Context.Local)
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'stack_get') {
    displayStackTrace(msg.Stacktrace)
  } else if (msg.MessageType === 'response' && msg.State === 'awake') {
    toggleOnOffbuttons()
    removePreviousBreak()
  } else if (msg.MessageType === 'response' && msg.State === 'asleep') {
    toggleOnOffbuttons()
    removePreviousBreak()
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
