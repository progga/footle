
/**
 * Feedback message display module.
 */

/**
 * Display a feedback message.
 *
 * @param string msg
 */
function show (msg) {
  const msgMarkup = `<div class="uk-alert feedback-message" data-uk-alert><a class="uk-alert-close uk-close"></a><p>${msg}</p></div>`

  const tmpElement = document.createElement('span')
  tmpElement.innerHTML = msgMarkup
  const msgElement = tmpElement.firstChild

  document.getElementById('messages').appendChild(msgElement)

  recordAreaHeight()
}

/**
 * Initialize the message display area.
 *
 * Whenever a message is closed, note down the new height of #messages.  This
 * height is used to push down all following elements.
 */
function init () {
  jQuery('#messages').on('closed.uk.alert', recordAreaHeight)
}

/**
 * Update the --messages-height CSS variable.
 *
 * The new value is the current height of #messages.  This is used to calculate
 * the position of tab contents and control buttons.
 */
function recordAreaHeight () {
  const messagesTotalHeight = document.getElementById('messages').offsetHeight
  document.querySelector(':root').style.setProperty('--messages-height', `${messagesTotalHeight}px`)
}

export {init, show}
