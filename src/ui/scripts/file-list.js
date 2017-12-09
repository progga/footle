
/**
 * @file
 * Make the raw file list more usable.
 *
 * Prepend breadcrumb and link to parent directory to file list.  Also apply
 * styling to improve the appearance of the file list.
 */

'use strict'

/**
 * Setup click handler on all *file* links in the file browser.
 *
 * @param object ignoredEvent
 *    Optional, it is okay to call a Javascript function without its arguments.
 */
function setupFileList (ignoredEvent) {
  setupFileLinks()
  improveFileListUX()
}

/**
 * Load file list or file.
 *
 * When a directory name has been clicked, load its file list.  When a filename
 * has been clicked, open that file in a new tab.
 *
 * Notice that directory names end in a slash, filenames do not.  Only exception
 * is the parent directory link.
 */
function setupFileLinks () {
  jQuery('pre', window.file_browser.document).on('click', 'a:not([href$="/"])', function (event) {
    var relativeFilepath = this.pathname.replace('/files/', '')
    addTab(relativeFilepath)

    return false
  })
}

/**
 * Add current dir name and parent dir link for better UX.
 */
function improveFileListUX () {
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
  for (let pathPart of relativePathParts) {
    pathPartsCumulative.push(pathPart)
    crumbs.push({dir: pathPart, url: siteURL + pathPartsCumulative.join('/')})
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

  for (let crumb of crumbList) {
    breadcrumb += '<li class="link--dir"><a href="' + crumb.url + '">' + crumb.dir + '</a></li>'
  }

  breadcrumb += '<li class="label--dir uk-active"><span>' + lastCrumb.dir + '</span></li></ul>'

  return breadcrumb
}
