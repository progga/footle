
/**
 * @file
 * Manage commands for the backend.
 */

/**
 * Send command to the backend.
 *
 * @param string command
 *   Example: breakpoint_set, breakpoint_list, step_over, etc.
 * @param array args
 *   [Optional] Any arguments needed by the command above.
 */
function sendCommand (command, args) {
  args = args || []

  jQuery.post('steering-wheel', {
    'msg': command + args.join(' ')
  }).fail(function (jqXHR, textStatus, errorThrown) {
    console.log(errorThrown)
  })
}
