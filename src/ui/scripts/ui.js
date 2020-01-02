
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
import * as feedback from './feedback.js'
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
  tab.setupScrollRestoration()
  control.setupContinuationControls()
  control.setupStateControl()
  breakpoint.setupTrigger()
  variable.setupInteraction()
  control.disable()
  feedback.init()
  applyInitialState()
  initServerMessageProcessing()
})

/**
 * Process responses from the Footle server.
 */
function initServerMessageProcessing () {
  const sse = new EventSource('/message-stream')
  let hasAttemptedReconnection = false

  jQuery(sse).on('message', function (event) {
    try {
      var msg = JSON.parse(event.originalEvent.data)
    } catch (e) {
      feedback.show('Trouble parsing JSON formatted response from Footle server.  More in console log.')
      console.log(e)
      console.log('Response was: ' + event.originalEvent.data)

      return
    }

    console.log(msg)

    processMsg(msg)
  })

  // When a connection is lost, attempt reconnection only once.
  jQuery(sse).on('error', function (event) {
    console.log(event)

    if (event.target.readyState === EventSource.CLOSED) {
      feedback.show('Footle server has gone away. Try reloading this page.')

      updateExecutionState('asleep')
      hasAttemptedReconnection = false
    } else if (event.target.readyState === EventSource.CONNECTING && hasAttemptedReconnection) {
      sse.close()
      feedback.show('Footle server has gone away. Try reloading this page.')

      hasAttemptedReconnection = false
      updateExecutionState('asleep')
    } else if (event.target.readyState === EventSource.CONNECTING && !hasAttemptedReconnection) {
      hasAttemptedReconnection = true
    }
  })
}

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
 * Fetch and process existing state and breakpoints.
 *
 * The state is represented in the form of a message.  This is the last message
 * that changed Footle's state to one of awake, asleep, or break.  Depending on
 * this message, display a break if needed.  Then activate or deactivate control
 * buttons.
 *
 * Also grab the list of existing breakpoints and display them.
 */
function applyInitialState () {
  jQuery.getJSON('current-state', messages => messages.forEach(processMsg))
    .fail(function (jqXHR, textStatus, errorThrown) {
      feedback.show('Failed to grab current state of Footle.  More in console log.')
      console.log(jqXHR)
    })
}
