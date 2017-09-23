
/**
 * @file
 * Behaviour for tabbed file browser and viewer.
 *
 * File browser is displayed in an iframe.  This iframe sits inside the first
 * tab.  When a filename is clicked in the file browser, it should be displayed
 * in another tab.  This tab will contain the HTML formatted file content.
 */

'use strict'

/**
 * Onload event handler.
 *
 * Does the following:
 *   - Sets up click handlers on file links in the file browser.
 *   - Sets up click handlers on tab close links.
 *   - Creates a Server-sent-event handler to listen to the data stream from the
 *     Footle server.
 */
jQuery(function () {
  /* We have missed the very first "load" event for the iframe.  So we
     need to explicitely call clickSetter() for setting up the click
     handlers on the file links inside the file browser. */
  clickSetter()

  /* iframe is loaded again when a directory is opened. */
  jQuery('iframe').on('load', clickSetter)

  /* Close a tab when its close link is clicked. */
  jQuery('.tab-nav').on('click', '.tab-closer', function (event) {
    event.preventDefault()
    event.stopPropagation()

    removeTabForFile(this.offsetParent.id)
  })

  var sse = new EventSource('/message-stream')
  jQuery(sse).on('message', function (event) {
    var msg = JSON.parse(event.originalEvent.data)
    console.log(msg)

    processMsg(msg)
  })
})

/**
 * Setup click handler on all *file* links in the file browser.
 *
 * @param object ignoredEvent
 *    Optional, it is okay to call a Javascript function without its arguments.
 */
function clickSetter (ignoredEvent) {
  /* Directory names end in a slash, filenames do not. */
  jQuery('pre :not(a[href$="/"])', window.file_browser.document).on('click', function (event) {
    var relativeFilepath = this.pathname.replace('/files/', '')
    addTab(relativeFilepath)

    event.preventDefault()
  })
}

/**
 * Update UI based on debugging status.
 *
 * @param object msg
 */
function processMsg (msg) {
  if (msg.MessageType === 'response' && msg.State === 'break' && msg.Properties.Filename) {
    updateBreak(msg.Properties.Filename, msg.Properties.LineNumber)
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'breakpoint_set') {
  } else if (msg.MessageType === 'response' && msg.Properties.Command === 'breakpoint_list') {
    refreshBreakpoints(msg.Breakpoints)
  } else if (msg.MessageType === 'response' && msg.State === 'stopped') {
    removePreviousBreak()
  } else if (msg.MessageType === 'init') {
  }
}
