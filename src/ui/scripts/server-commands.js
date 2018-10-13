
/**
 * @file
 * Manage commands sent to the Footle server.
 */

/**
 * Send command to the Footle server.
 *
 * @param string command
 *   Example: breakpoint_set, breakpoint_list, step_over, etc.
 * @param array args
 *   [Optional] Any arguments needed by the command above.
 *
 * Example *Footle* command: breakpoint_set index.php 16
 */
function sendCommand (command, args) {
  args = args || []

  var footleCommand = [command].concat(args).join(' ')

  jQuery.post('steering-wheel', {
    'cmd': footleCommand
  }).fail(function (jqXHR, textStatus, errorThrown) {
    console.log(errorThrown)
  })
}

export default sendCommand
