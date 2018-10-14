
/**
 * @file
 * Update source code for a displayed file.
 */

import * as tab from './tabs.js'
import * as util from './common.js'

/**
 * Update the displayed source code of the given file.
 *
 * @param string filepath
 *  Relative filepath.
 */
function update (filepath) {
  var isUnknownFile = !tab.hasFileMapping(filepath)
  if (isUnknownFile) {
    return
  }

  var formattedFilepath = '/formatted-file/' + filepath
  var fileTabElement = tab.getContentElementForFile(filepath)
  var fileContentElement = jQuery('.file-content', fileTabElement)

  jQuery(fileContentElement).load(formattedFilepath, function () {
    var fileTabLinkElementSelector = '#' + util.escapeSelector(filepath)
    jQuery(fileTabLinkElementSelector).addClass('uk-animation-slide-top')

    // Remove the animation so that it can be applied again.
    window.setTimeout(function () {
      jQuery(fileTabLinkElementSelector).removeClass('uk-animation-slide-top')
    }, 1000) // Because the uk-shake animation lasts for 500ms.
  })
}

export {update}
