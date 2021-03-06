
/**
 * @file
 * Make the raw file list more usable.
 *
 * Prepend breadcrumb and link to parent directory to file list.  Also apply
 * styling to improve the appearance of the file list.
 */

import * as breakpoint from './breakpoints.js'
import * as tab from './tabs.js'
import RecentFiles from './recent-files.js'

/**
 * Setup file browser.
 *
 * @param object ignoredEvent
 *    Optional, it is okay to call a Javascript function without its arguments.
 */
function setup (ignoredEvent) {
  setupFileLinks()
  improveFileListUX()

  const recentFiles = new RecentFiles(localStorage).get()
  displayRecentFiles(recentFiles)
}

/**
 * Attach click handlers to filenames listed as recently used.
 */
function setupRecent () {
  jQuery('.file-list--recent').off('click', '.file--recent__link').on('click', '.file--recent__link', function (event) {
    // Even though we don't use absolute filenames in recent file links, that's
    // what we get here.  But we don't want a leading slash in the filepath for
    // display purposes.
    const relativeFilepath = this.pathname.replace(/^\//, '')
    tab.add(relativeFilepath, postFileOpenTasks)

    return false
  })
}

/**
 * Setup click handler on all *file* links.
 *
 * These file links are found as part of the file browser and recent file list.
 *
 * When a directory name has been clicked, load its file list.  When a filename
 * has been clicked, open that file in a new tab.
 *
 * Notice that directory names end in a slash, filenames do not.  Only exception
 * is the parent directory link.
 */
function setupFileLinks () {
  jQuery('pre', window.file_browser.document).off('click', 'a:not([href$="/"])').on('click', 'a:not([href$="/"])', function (event) {
    const relativeFilepath = this.pathname.replace('/files/', '')
    tab.add(relativeFilepath, postFileOpenTasks)

    return false
  })
}

/**
 * Things we want to do once a file is opened in a tab.
 *
 * At the moment, we are:
 * - Updating the recently used file list.
 * - Redrawing any existing breakpoints.
 */
function postFileOpenTasks (filename, filepath) {
  updateRecent(filepath)
  breakpoint.highlightFile(filepath)
}

/**
 * Add current dir name and parent dir link for better UX.
 */
function improveFileListUX () {
  const hasFinishedSetup = window.file_browser.document.querySelector('body > ul.breadcrumb--file-path:first-child')
  if (hasFinishedSetup) {
    return
  }
  // Hide file list until stylesheet is ready and all modification is complete.
  jQuery('body', window.file_browser.document).hide()

  // Load stylesheet for file links which are inside the file browser iframe.
  // Show the file list once the stylesheet has been loaded.
  jQuery('head', window.file_browser.document).append(jQuery('<link rel="stylesheet" href="/style/css/ui.css" />').on('load', function () {
    jQuery('body', window.file_browser.document).addClass('file-list')
    jQuery('body', window.file_browser.document).show()
  }))

  // Add breadcrumb...
  var crumbs = prepareCrumbs(document.URL, window.file_browser.document.URL)
  if (crumbs.length <= 1) {
    // No breadcrumb when we are listing files at the docroot.  Note that
    // the docroot is always the first crumb.
    return
  }

  var breadcrumb = prepareBreadcrumbMarkup(crumbs)
  jQuery('pre', window.file_browser.document).before(breadcrumb)

  // ... and then a file link pointing to the parent directory.
  var parentDirURL = crumbs[crumbs.length - 2].url
  var parentDirLink = '<a class="link--parent-dir" href="' + parentDirURL + '/">..</a>&#10;'
  jQuery('pre', window.file_browser.document).prepend(parentDirLink)
}

/**
 * Prepare items for breadcrumb for the current path.
 *
 * @param string siteURL
 *    Example: https://example.net/
 * @param string currentDirURL
 *    Example: https://example.net/foo/bar/
 *
 * @return Array
 *    Array of objects.  Each object has two keys: dir, url.
 */
function prepareCrumbs (siteURL, currentDirURL) {
  var relativePathParts = currentDirURL.replace(siteURL, '').split('/').filter(dir => dir !== '')
  var crumbs = []

  var pathPartsCumulative = []
  for (const pathPart of relativePathParts) {
    pathPartsCumulative.push(pathPart)
    crumbs.push({ dir: pathPart, url: siteURL + pathPartsCumulative.join('/') })
  }

  return crumbs
}

/**
 * Prepare HTML markup for breadcrumb from given crumbs.
 *
 * @param Array crumbList
 *    Array of objects.  Each object has two keys: dir, url that represents a
 *    directory path.
 *
 * @return string
 *    HTML markup for breadcrumb.
 */
function prepareBreadcrumbMarkup (crumbs) {
  if (crumbs.length <= 0) {
    return ''
  }

  // Clone crumbs as we do not want side effects.
  var crumbList = JSON.parse(JSON.stringify(crumbs))

  // Rename the first crumb from "files" to a more meaningful "Docroot".
  crumbList[0].dir = 'Docroot'

  var breadcrumb = '<ul class="breadcrumb--file-path uk-breadcrumb">'
  var lastCrumb = crumbList.pop()

  for (const crumb of crumbList) {
    breadcrumb += '<li class="link--dir"><a href="' + crumb.url + '">' + crumb.dir + '</a></li>'
  }

  breadcrumb += '<li class="label--dir uk-active"><span>' + lastCrumb.dir + '</span></li></ul>'

  return breadcrumb
}

/**
 * We want to list the last five files opened.
 *
 * These are the files that have been explicitely opened by clicking their
 * names.  Files opened *automatically* as part of the step through execution
 * process are excluded from this list.
 */
function displayRecentFiles (filelist) {
  const listMarkup = filelist.map(filename => `<li class="file--recent"><a href="${filename}" class="file--recent__link">${filename}</a></li>`).join('\n')

  jQuery('.file-list--recent').html(listMarkup)

  // Only display the "Recent files" header when we have got some such files.
  if (filelist.length) {
    jQuery('.recent-files__header').removeClass('uk-invisible')
  }
}

/**
 * The list of recent files is constantly updated.
 *
 * Everytime a filename is *clicked*, it moves to the top of the list.
 */
function updateRecent (filename) {
  const recentFiles = new RecentFiles(localStorage)
  recentFiles.add(filename)

  const filelist = recentFiles.get()
  displayRecentFiles(filelist)
}

export { setup, setupRecent }
