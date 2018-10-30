
/**
 * @file
 * Manage commands sent to the Footle server.
 */

import * as feedback from './feedback.js'

/**
 * Send command to the Footle server.
 *
 * @param string command
 *   Example: breakpoint_set, breakpoint_list, step_over, etc.
 * @param array args
 *   [Optional] Any arguments needed by the command above.
 *
 * Example *Footle* command: breakpoint_set index.php 16
 *
 * All commands have a fixed response when successful: "Got it."
 */
function sendCommand (command, args) {
  args = args || []

  var footleCommand = [command].concat(args).join(' ')

  jQuery.post('steering-wheel', {
    'cmd': footleCommand
  }).done(function (data, textStatus, jqXHR) {
    const cmdHasSucceeded = (data !== 'Got it.')
    if (cmdHasSucceeded) {
      feedback.show(`The "${footleCommand}" command failed.  More in console log.`)
      console.log(jqXHR)
    }
  })
    .fail(function (jqXHR, textStatus, errorThrown) {
      feedback.show(`The "${footleCommand}" command failed.  More in console log.`)
      console.log(jqXHR)
    })
}

export {sendCommand}
