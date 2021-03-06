
/**
 * @file
 * Tab management.
 *
 * Each source file is displayed in a tab.  To manage this, we keep a mapping
 * of filenames and their corresponding tabs.
 */

import * as server from './server-commands.js'
import * as feedback from './feedback.js'

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
function add (filepath, postTabOpenAction) {
  const filename = filepath.split(/[\\/]/).pop()
  const formattedFilepath = '/formatted-file/' + filepath

  if (hasFileMapping(filepath)) {
    if (postTabOpenAction) {
      postTabOpenAction(filename, filepath)
    }

    return
  }

  jQuery.get(formattedFilepath, function (data) {
    // The mapping may have been updated in the meantime.  Hence this recheck.
    if (hasFileMapping(filepath)) {
      if (postTabOpenAction) {
        postTabOpenAction(filename, filepath, data)
      }

      return
    }

    /* Tab link. */
    const tabLink = jQuery(`<li id="${filepath}" class="tab-selector" data-filepath="${filepath}" title="${filepath}"><a href="#"><span class="tab-refresh" title="Reload"></span>${filename}<span class="tab-closer" title="Close"></span></a></li>`)
    // Make sure to add a jQuery object instead of plain markup so that we can
    // use the object later.
    jQuery('#tab-selector-wrapper').append(tabLink)

    /* Tab content. */
    const tabContent = `<li id="body-of-${filepath}" class="tab-content" data-filepath="${filepath}"><div class="file-content">${data}</div></li>`
    jQuery('#tab-content-wrapper').append(tabContent)

    /* Record the presence of a tab for this file. */
    addFileMapping(filepath, tabLink)

    /* Activate tab. */
    open(filepath)

    // The tab content area should start after the fix positioned tab
    // selector area.
    recordTabHeight()

    if (postTabOpenAction) {
      postTabOpenAction(filename, filepath, data)
    }
  }).fail(function (jqXHR, textStatus, errorThrown) {
    feedback.show(`Failed to fetch file ${filename}.  More in console log.`)
    console.log(jqXHR)
  })
}

/**
 * Reload a tab's content when its refresh link is clicked.
 */
function setupRefresher () {
  jQuery('#tab-selector-wrapper').on('click', '.tab-refresh', function (event) {
    event.preventDefault()
    event.stopPropagation()

    const filepath = this.offsetParent.id
    server.sendCommand('update_source', [filepath])
  })
}

/**
 * Close a tab when its close link is clicked.
 */
function setupCloser () {
  jQuery('#tab-selector-wrapper').on('click', '.tab-closer', function (event) {
    event.preventDefault()
    event.stopPropagation()

    remove(this.offsetParent.id)
  })
}

/**
 * Arrange to restore scroll position when we return to a tab.
 *
 * UIKit does not keep track of the last scroll position of each tab.  To
 * understand, take this sequence of steps:
 * - Open tab A.
 * - Scroll a few lines.  Note the position of the scrollbar.
 * - Open tab B but do *not* scroll.
 * - Return to tab A.
 * The scrollbar will *not* be where it was when we last left tab A.  Here
 * we attempt to fill this gap for the *vertical* scrollbar.  We are ignoring
 * the horizontal one for now as that is less of a concern.
 *
 * Note:
 * - The change.uk.tab UIKit event is not working in Footle although it
 *   works elsewhere.  We are using the show.uk.switcher event instead.
 * - UIKit 3 has a better solution as it offers more tab switch related events.
 *   The "beforeshow", "show", and "beforehide" events (in that order) are
 *   relevant.
 */
function setupScrollRestoration () {
  let previousTab, lastScrolltop

  // Before we switch tab, keep track of the scrolling position of the last tab.
  jQuery('#tab-selector-wrapper').on('click', '.tab-selector', function (event) {
    lastScrolltop = document.getElementsByTagName('body')[0].scrollTop
  })

  // We have just switched tab, so restore the scrolling position of this tab.
  jQuery('#tab-selector-wrapper').on('show.uk.switcher', function (event, activeTab) {
    let scrolltop = activeTab[0].getAttribute('data-last-scrolltop')
    if (!scrolltop) {
      scrolltop = 0
    }
    document.getElementsByTagName('body')[0].scrollTo(0, scrolltop)

    // Now that we know the last scroll position of the *previous* tab, save it
    // in that tab for later use.
    if (previousTab) {
      previousTab.setAttribute('data-last-scrolltop', lastScrolltop)
    }

    // Save this tab, because we will add its last scroll position to it *after*
    // we switch to another tab.  See the conditional block above.
    previousTab = activeTab[0]
  })
}

/**
 * Note down a file and its associated tab element.
 *
 * @param string filepath
 * @param object tabElement
 *    jQuery object for a tab.
 */
function addFileMapping (filepath, tabElement) {
  fileTabMapping[filepath] = tabElement
}

/**
 * Remove association record between a file and its tab element.
 *
 * @param string filepath
 */
function removeFileMapping (filepath) {
  if (!hasFileMapping(filepath)) {
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
function hasFileMapping (filepath) {
  if (Object.prototype.hasOwnProperty.call(fileTabMapping, filepath)) {
    return fileTabMapping[filepath]
  }

  return false
}

/**
 * Open the tab for the given file.
 *
 * @param string filepath
 */
function open (filepath) {
  if (!hasFileMapping(filepath)) {
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
function remove (filepath) {
  let tabElement

  if (!(tabElement = hasFileMapping(filepath))) {
    return
  }

  /* Delete tab and its content */
  const tabContentElement = getContentElement(tabElement)
  tabElement.remove()
  tabContentElement.remove()

  removeFileMapping(filepath)

  recordTabHeight()

  /* When we are closing the active tab, return to file browser in first tab. */
  const isActiveTab = tabElement.hasClass('uk-active')
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
function getContentElement (tabElement) {
  const tabIndex = tabElement.index()
  const tabContentElement = jQuery('.tab-content').get(tabIndex)

  return tabContentElement
}

/**
 * Given a filename, find the tab for that file.
 *
 * @param string filename
 * @return object/null
 */
function getContentElementForFile (filename) {
  const tabContentId = 'body-of-' + filename

  const tabContent = document.getElementById(tabContentId)

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
 * @see --tab-selector-height CSS variable.
 * @todo Adjust top padding on window resize.
 */
function recordTabHeight () {
  const height = document.getElementById('tab-selector-wrapper').offsetHeight
  document.querySelector(':root').style.setProperty('--tab-selector-height', `${height}px`)
}

export { add, getContentElement, getContentElementForFile, hasFileMapping, setupRefresher, setupCloser, setupScrollRestoration }
