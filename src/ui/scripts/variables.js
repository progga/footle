
/**
 * @file
 * Variable display.
 */

'use strict'

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
 * Include a variable's children to the display.
 *
 * Example: {foo: ['bar', 'buz']}. Assuming "foo" is already part of the
 * display, this function will add ['bar', 'buz'] to the variable tree.
 *
 * @param object varDetailList
 *   Record of a single variable and its children upto a certain depth.
 *   This may be an object, but it is expected to have only one property.
 */
function displaySingleVar (varDetailList) {
  var varFullname = Object.keys(varDetailList)[0]

  if (varFullname === undefined) {
    return
  }

  if (!varDetailList[varFullname].hasOwnProperty('Children')) {
    return
  }

  var childrenVars = varDetailList[varFullname].Children
  var varListMarkup = listBasicVars(childrenVars)

  var varIdSelector = '#' + escapeSelector(varFullname)
  var varChildrenSelector = varIdSelector + ' > .variable-list'

  jQuery(varChildrenSelector).replaceWith(varListMarkup)
  jQuery(varIdSelector).attr('data-has-loaded-children', 'true')
}

/**
 * Setup variable interaction events.
 *
 * Current events:
 * - Click handler for collapsing/uncollapsing the variable tree.
 */
function setupVariableInteraction () {
  // When a variable with children is clicked, collapse it.
  jQuery('.variables').on('click', '.variable[data-is-composite="true"]', function (event) {
    // Has the click been on a variable with children?  Only act on clicks
    // that are on the list item surrounding a variable name or the variable
    // name itself.
    if (!jQuery(event.target).is('.variable[data-is-composite="true"], .variable[data-is-composite="true"] > .variable__display-name')) {
      return false
    }

    jQuery(this).toggleClass('expanded')

    var hasNotYetLoadedChildren = jQuery(this).is('.expanded[data-has-loaded-children="false"]')
    if (hasNotYetLoadedChildren) {
      var varName = jQuery(this).attr('data-var-fullname')

      sendCommand('property_get', [varName])
    }

    // We do *not* want to expand/collapse the parent variables of the clicked
    // variable.
    return false
  })
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
    varType += ' (' + varDetail.AccessModifier + ')'
  }

  if (varDetail.IsCompositeType && !varDetail.HasLoadedChildren) {
    childrenMarkup = '<ul class="variable-list wait--loading-children uk-icon-refresh uk-icon-spin"></ul>'
  }

  var markup = '<li id="' + varFullname + '"' +
                 ' class="variable" data-var-fullname="' + varFullname + '"' +
                 ' data-is-composite="' + varDetail.IsCompositeType + '"' +
                 ' data-has-loaded-children="' + varDetail.HasLoadedChildren + '">' +
                 '<span class="variable__display-name">' + varDetail.DisplayName + '</span>' +
                 '<span class="variable__type">' + varType + '</span>' +
                 '<span class="variable__value">' + varDetail.Value + '</span>' +
                 childrenMarkup +
               '</li>'

  return markup
}
