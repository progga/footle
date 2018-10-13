
/**
 * @file
 * Update source code for a displayed file.
 */

import {hasFileTabMapping, getTabContentElementForFile} from './tabs.js'
import {escapeSelector} from './common.js'

/**
 * Update the displayed source code of the given file.
 *
 * @param string filepath
 *  Relative filepath.
 */
function updateSourceFile (filepath) {
  var isUnknownFile = !hasFileTabMapping(filepath)
  if (isUnknownFile) {
    return
  }

  var formattedFilepath = '/formatted-file/' + filepath
  var fileTabElement = getTabContentElementForFile(filepath)
  var fileContentElement = jQuery('.file-content', fileTabElement)

  jQuery(fileContentElement).load(formattedFilepath, function () {
    var fileTabLinkElementSelector = '#' + escapeSelector(filepath)
    jQuery(fileTabLinkElementSelector).addClass('uk-animation-slide-top')

    // Remove the animation so that it can be applied again.
    window.setTimeout(function () {
      jQuery(fileTabLinkElementSelector).removeClass('uk-animation-slide-top')
    }, 1000) // Because the uk-shake animation lasts for 500ms.
  })
}

export default updateSourceFile
