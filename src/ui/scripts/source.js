
/**
 * @file
 * Update source code for a displayed file.
 */

'use strict'

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
    jQuery(fileTabLinkElementSelector).addClass('uk-animation-shake')

    // Remove the animation so that it can be applied again.
    window.setTimeout(function () {
      jQuery(fileTabLinkElementSelector).removeClass('uk-animation-shake')
    }, 1000) // Because the uk-shake animation lasts for 500ms.
  })
}
