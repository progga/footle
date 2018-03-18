
/**
 * @file
 * Actions for continuation and state display commands.
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
    'on': '[name="button--on"]',
    'off': '[name="button--off"]'
  }

  setupCommandNControl(commandsNSelectors)

  /**
   * Display either the "On" or the "Off" button.
   *
   * @see Markup for the "On" button.
   */
  jQuery('[name="button--on"], [name="button--off"]').click(function () {
    jQuery('[name="button--on"], [name="button--off"]').toggleClass('uk-hidden')
  })
}

/**
 * Prepare handlers for state update buttons.
 *
 * Setup control buttons for fetching variables and call stack.
 */
function setupStateControl () {
  var commandsNSelectors = {
    'context_get': '[name="button--variable__local"]',
    'context_get global': '[name="button--variable__global"]',
    'stack_get': '[name="button--stacktrace"]'
  }

  setupCommandNControl(commandsNSelectors)
}

/**
 * Create association between debugging commands and their corresponding buttons.
 *
 * @return object
 *  Key: Command; Value: CSS selector.
 */
function setupCommandNControl (commandsNSelectors) {
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
