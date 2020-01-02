
/**
 * @file
 * Functions common to all elements.
 */

/**
 * Escape CSS selector.
 *
 * @param string selector
 * @return string
 *
 * @see CSS.escape()
 * @see jQuery.escapeSelector()
 * @see https://learn.jquery.com/using-jquery-core/faq/how-do-i-select-an-element-by-an-id-that-has-characters-used-in-css-notation/
 *
 * @todo Add unit test.
 */
function escapeSelector (selector) {
  var escapedSelector = ''

  // The CSS.escape() function is in Draft status.  So unavailable in some
  // browsers.
  if ((typeof CSS === 'function' || typeof CSS === 'object') && typeof CSS.escape === 'function') {
    escapedSelector = CSS.escape(selector)
  } else {
    escapedSelector = selector.replace(/(:|\.|\[|\]|,|=|@|%)/g, '\\$1')
  }

  return escapedSelector
}

export { escapeSelector }
