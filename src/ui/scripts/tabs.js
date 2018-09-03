
/**
 * @file
 *  Tab management.
 *
 * Each source file is displayed in a tab.  To manage this, we keep a mapping
 * of filenames and their corresponding tabs.
 */

'use strict'

/**
 * List of files and their corresponding tab elements.
 */
var fileTabMapping = {}

/**
 * Add file content inside a tab.
 *
 * @param string filepath
 *    Relative filepath.
 * @param callback postTabOpenAction
 *    Call this function once the tab is fully prepared.
 */
function addTab (filepath, postTabOpenAction) {
  var filename = filepath.split(/[\\/]/).pop()
  var formattedFilepath = '/formatted-file/' + filepath

  if (hasFileTabMapping(filepath)) {
    return
  }

  jQuery.get(formattedFilepath, function (data) {
    if (hasFileTabMapping(filepath)) {
      // The mapping may have been updated in the meantime.  Hence this recheck.
      return
    }

    /* Tab link. */
    var tabLink = jQuery(`<li id="${filepath}" class="tab-selector" data-filepath="${filepath}" title="${filepath}"><a href="#"><span class="tab-refresh" title="Reload"></span>${filename}<span class="tab-closer" title="Close"></span></a></li>`)
    // Make sure to add a jQuery object instead of plain markup so that we can
    // use the object later.
    jQuery('#tab-selector-wrapper').append(tabLink)

    /* Tab content. */
    var tabContent = `<li id="body-of-${filepath}" class="tab-content" data-filepath="${filepath}"><div class="file-content">${data}</div></li>`
    jQuery('#tab-content-wrapper').append(tabContent)

    /* Record the presence of a tab for this file. */
    addFileTabMapping(filepath, tabLink)

    /* Activate tab. */
    openTabForFile(filepath)

    /* Update recent file list. */
    updateRecentFiles(filepath)

    // The tab content area should start after the fix positioned tab
    // selector area.
    adjustTabContentPosition()

    if (postTabOpenAction) {
      postTabOpenAction()
    }
  })
}

/**
 * Reload a tab's content when its refresh link is clicked.
 */
function setupTabRefresher () {
  jQuery('#tab-selector-wrapper').on('click', '.tab-refresh', function (event) {
    event.preventDefault()
    event.stopPropagation()

    var filepath = this.offsetParent.id
    sendCommand('update_source', [filepath])
  })
}

/**
 * Close a tab when its close link is clicked.
 */
function setupTabCloser () {
  jQuery('#tab-selector-wrapper').on('click', '.tab-closer', function (event) {
    event.preventDefault()
    event.stopPropagation()

    removeTabForFile(this.offsetParent.id)
  })
}

/**
 * Note down a file and its associated tab element.
 *
 * @param string filepath
 * @param object tabElement
 *    jQuery object for a tab.
 */
function addFileTabMapping (filepath, tabElement) {
  fileTabMapping[filepath] = tabElement
}

/**
 * Remove association record between a file and its tab element.
 *
 * @param string filepath
 */
function removeFileTabMapping (filepath) {
  if (!hasFileTabMapping(filepath)) {
    return
  }

  delete fileTabMapping[filepath]
}

/**
 * Is there a tab for the given file?
 *
 * Returns the tab element when a tab is present.
 *
 * @param string filepath
 * @return bool
 */
function hasFileTabMapping (filepath) {
  if (fileTabMapping.hasOwnProperty(filepath)) {
    return fileTabMapping[filepath]
  }

  return false
}

/**
 * Open the tab for the given file.
 *
 * @param string filepath
 */
function openTabForFile (filepath) {
  if (!hasFileTabMapping(filepath)) {
    return
  }

  fileTabMapping[filepath].click()
}

/**
 * Manage tab closing.
 *
 * - Close the tab.
 * - Remove association between a file and its tab element.
 * - When the current tab is closed, go back to the file browser.
 *
 * @param string filepath
 */
function removeTabForFile (filepath) {
  var tabElement

  if (!(tabElement = hasFileTabMapping(filepath))) {
    return
  }

  /* Delete tab and its content */
  var tabContentElement = getTabContentElement(tabElement)
  tabElement.remove()
  tabContentElement.remove()

  removeFileTabMapping(filepath)

  adjustTabContentPosition()

  /* When we are closing the active tab, return to file browser in first tab. */
  var isActiveTab = tabElement.hasClass('uk-active')
  if (isActiveTab) {
    jQuery('.tab-selector').get(0).click()
  }
}

/**
 * Find the content element for the given tab.
 *
 * Each tab element has a corresponding content element where the tab's content
 * is displayed.  Here we find that content element.
 *
 * The find the content element, we look at the index of the tab element.  So
 * when the tab is the fifth child of its parent, the tab content element is
 * also the fifth child.  That is how we identify the content element.
 *
 * @param object tabElement
 *    jQuery object for a tab header.
 * @return object
 *    jQuery object for the given tab's body.
 */
function getTabContentElement (tabElement) {
  var tabIndex = tabElement.index()
  var tabContentElement = jQuery('.tab-content').get(tabIndex)

  return tabContentElement
}

/**
 * Given a filename, find the tab for that file.
 *
 * @param string filename
 * @return object/null
 */
function getTabContentElementForFile (filename) {
  var tabContentId = 'body-of-' + filename

  var tabContent = document.getElementById(tabContentId)

  if (tabContent) {
    return jQuery(tabContent)
  } else {
    return tabContent
  }
}

/**
 * Position tab content after tab selector.
 *
 * The tab selector is sticky at the top thanks to its *fixed* positioning.
 * Unless we do something, the tab content area will flow under it.  We
 * avoid it by pushing down the tab content area as much as the height of
 * the tab selector.  To push it down, we adjust the top padding of the
 * *whole tab* area everytime the height of the tab selector area *may* change.
 *
 * @todo Adjust top padding on window resize.
 */
function adjustTabContentPosition () {
  var tabSelectorHeight = jQuery('#tab-selector-wrapper').outerHeight(true)
  jQuery('#tab').css('padding-top', tabSelectorHeight)
}
