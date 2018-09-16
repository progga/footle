
/**
 * @file
 * Variable display.
 */

import sendCommand from './server-commands'
import {escapeSelector} from './common'

/**
 * Display variables.
 *
 * @param object varDetailList
 */
function updateVarsDisplay (variables) {
  var hasLocalVarDetails = Array.isArray(variables.Local) && (variables.Local.length > 0)
  var hasGlobalVarDetails = Array.isArray(variables.Global) && (variables.Global.length > 0)
  var varListMarkup

  if (hasLocalVarDetails) {
    varListMarkup = listBasicVars(variables.Local)

    jQuery('.variables').html(varListMarkup)
    jQuery('.variables').attr('data-var-context', 'local')
  } else if (hasGlobalVarDetails) {
    varListMarkup = listBasicVars(variables.Global)

    jQuery('.variables').html(varListMarkup)
    jQuery('.variables').attr('data-var-context', 'global')
  }
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
  var hasNoSubstance = !(Array.isArray(varDetailList) && (varDetailList.length > 0))
  if (hasNoSubstance) {
    return
  }

  var varDetail = varDetailList[0]

  if (!varDetail.hasOwnProperty('Fullname') || !varDetail.hasOwnProperty('Children')) {
    return
  }

  var childrenVars = varDetail.Children
  var varListMarkup = listBasicVars(childrenVars)

  var varFullname = varDetail.Fullname
  var varIdSelector = '#' + escapeSelector(escape(varFullname))
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
      var varName = unescape(jQuery(this).attr('data-var-fullname'))
      var varContext = jQuery(event.delegateTarget).attr('data-var-context')

      sendCommand('property_get', [varContext, varName])
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
  var hasNoSubstance = !(Array.isArray(varDetailList) && (varDetailList.length > 0))
  if (hasNoSubstance) {
    return
  }

  var markup = ''
  var childrenMarkup = ''

  for (let i = 0; i < varDetailList.length; i++) {
    var varDetail = varDetailList[i]
    childrenMarkup = ''

    if (varDetail.IsCompositeType && varDetail.HasLoadedChildren) {
      childrenMarkup = listBasicVars(varDetail.Children)
    }

    markup += prepareVarMarkup(varDetail, childrenMarkup)
  }

  markup = '<ul class="variable-list">' + markup + '</ul>'

  return markup
}

/**
 * Prepare the markup for a single variable.
 *
 * @param object varDetail
 * @return string
 */
function prepareVarMarkup (varDetail, childrenMarkup) {
  var varType = varDetail.VarType
  var varFullname = varDetail.Fullname

  // Example type description: string, string (private).
  if (varDetail.AccessModifier.length > 0) {
    varType += ' (' + varDetail.AccessModifier + ')'
  }

  if (varDetail.IsCompositeType && !varDetail.HasLoadedChildren) {
    childrenMarkup = '<ul class="variable-list wait--loading-children uk-icon-refresh uk-icon-spin"></ul>'
  }

  var markup = '<li id="' + escape(varFullname) + '"' +
                 ' class="variable" data-var-fullname="' + escape(varFullname) + '"' +
                 ' data-is-composite="' + varDetail.IsCompositeType + '"' +
                 ' data-has-loaded-children="' + varDetail.HasLoadedChildren + '">' +
                 '<span class="variable__display-name">' + varDetail.DisplayName + '</span>' +
                 '<span class="variable__type">' + varType + '</span>' +
                 '<span class="variable__value">' + varDetail.Value + '</span>' +
                 childrenMarkup +
               '</li>'

  return markup
}

export {displaySingleVar, setupVariableInteraction, updateVarsDisplay}
