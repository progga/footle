
/**
 * @file
 * Actions for continuation commands
 */

/**
 * Prepare handlers for continuation buttons.
 *
 * Setup click handlers for continuation buttons of the following
 * commands: step_over, step_into, step_out, and run.
 */
function setupContinuationControls () {
  var commandsNSelectors = {
    'step_over': '[name="button--step-over"]',
    'step_into': '[name="button--step-in"]',
    'step_out': '[name="button--step-out"]',
    'run': '[name="button--run"]',
    'continue': '[name="button--continue"]',
    'context_get': '[name="button--variable__local"]',
    'context_get global': '[name="button--variable__global"]',
    'stack_get': '[name="button--stacktrace"]'
  }

  var command = ''
  var selector = ''

  for (command in commandsNSelectors) {
    selector = commandsNSelectors[command]

    jQuery(selector).click(command, function (event) {
      event.preventDefault()

      sendCommand(event.data)
    })
  }
}
