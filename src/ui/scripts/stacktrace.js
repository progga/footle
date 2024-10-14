/**
 * @file
 * Stack trace display related functions.
 */

/**
 * Display trace in a tabular format.
 *
 * @param array callStack
 */
function display (callStack) {
  jQuery('.stacktrace--tabular > .traces').empty()

  if (callStack.length === 0) {
    jQuery('.stacktrace--tabular > .traces').append('<tr><td>Empty stack.</td></tr>')
  }

  for (var stackIndex in callStack) {
    var where = callStack[stackIndex].Where
    var filename = callStack[stackIndex].Filename
    var lineNo = callStack[stackIndex].LineNo

    var traceMarkup = '<tr class="call-detail">' +
                        '<td>' + where + '</td>' +
                        '<td class="filename">' + filename + '</td>' +
                        '<td>' + lineNo + '</td>' +
                      '</tr>'

    jQuery('.stacktrace--tabular > .traces').append(traceMarkup)
  }

  // Initially, the table remains hidden to avoid displaying table headers for
  // an empty table.
  jQuery('.stacktrace--tabular').removeClass('uk-hidden')
}

export { display }
