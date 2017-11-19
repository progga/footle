
/**
 * @file
 * Variable display.
 */

'use strict'

/**
 * Click handler for collapsing/uncollapsing the variable tree.
 *
 * When a line number is clicked, a breakpoint is added for that line.
 */
function setupVariableDisplay () {
  // When a variable with children is clicked, collapse it.
  jQuery('.variables').on('click', '.variable[data-is-composite="true"]', function (event) {
    jQuery(event.target).toggleClass('expanded')
  })
}

/**
 * Display variables.
 *
 * @param object varDetailList
 */
function updateVarsDisplay (varDetailList) {
  var varListMarkup = listBasicVars(varDetailList)

  jQuery('.variables').html(varListMarkup)
}

/**
 * Prepare list markup for given variables.
 *
 * @param object varDetailList
 * @return string
 */
function listBasicVars (varDetailList) {
  var markup = ''
  var childrenMarkup = ''

  for (var varFullname in varDetailList) {
    if (!varDetailList.hasOwnProperty(varFullname)) {
      continue
    }

    var varDetail = varDetailList[varFullname]
    childrenMarkup = ''

    if (varDetail.IsCompositeType && varDetail.HasLoadedChildren) {
      childrenMarkup = listBasicVars(varDetail.Children)
    }

    markup += prepareVarMarkup(varFullname, varDetail, childrenMarkup)
  }

  markup = '<ul class="variable-list">' + markup + '</ul>'

  return markup
}

/**
 * Prepare the markup for a single variable.
 *
 * @param string varFullname
 * @param object varDetail
 * @return string
 */
function prepareVarMarkup (varFullname, varDetail, childrenMarkup) {
  var varType = varDetail.VarType

  // Example type description: string, string (private).
  if (varDetail.AccessModifier.length > 0) {
    varType += '(' + varDetail.AccessModifier + ')'
  }

  var markup = '<li class="variable" data-var-fullname="' + varFullname +
                 '" data-is-composite="' + varDetail.IsCompositeType + '">' +
                 '<span class="variable__display-name">' + varDetail.DisplayName + '</span>' +
                 '<span class="variable__type">' + varType + '</span>' +
                 '<span class="variable__value">' + varDetail.Value + '</span>' +
                 childrenMarkup +
               '</li>'

  return markup
}
