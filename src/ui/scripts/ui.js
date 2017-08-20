
/**
 * @file
 * Behaviour for tabbed file browser and viewer.
 *
 * File browser is displayed in an iframe.  This iframe sits inside the first
 * tab.  When a filename is clicked in the file browser, it should be displayed
 * in another tab.  This tab will contain the HTML formatted file content.
 */

"use strict";

/**
 * List of files and their corresponding tab element.
 */
var fileTabMapping = {}

jQuery(function() {

  /* We have missed the very first "load" event for the iframe.  So we
     need to explicitely call clickSetter() for setting up the click
     handlers on the file links inside the file browser. */
  clickSetter();

  /* iframe is loaded again when a directory is opened. */
  jQuery("iframe").on("load", clickSetter);

  /* Close a tab when its close link is clicked. */
  jQuery(".tab-nav").on("click", ".tab-closer", function(event) {
    event.preventDefault();
    event.stopPropagation();

    removeTabForFile(this.offsetParent.id)
  });
})

/**
 * Setup click handler on all *file* links in the file browser.
 *
 * @param object ignoredEvent
 *    Optional, it is okay to call a Javascript function without its arguments.
 */
function clickSetter(ignoredEvent) {

  /* Directory names end in a slash, filenames do not. */
  jQuery("pre :not(a[href$=\"/\"])", window.file_browser.document).on("click", function(event) {
    addTab(this.pathname);

    event.preventDefault();
  });
}

/**
 * Add file content inside a tab.
 *
 * @param string filepath
 */
function addTab(filepath) {

  var filename = filepath.split(/[\\/]/).pop();
  var formattedFilepath = filepath.replace("/files/", "/formatted-file/");
  var filepathRelativeToDocroot = filepath.replace("/files/", "");

  if (hasFileTabMapping(filepathRelativeToDocroot)) {
    return;
  }

  jQuery.get(formattedFilepath, function(data) {

    /* Tab link */
    var tabLink = jQuery("<li " + "id=\"" + filepathRelativeToDocroot + "\" class=\"tab-selector\"><a href=\"#\">" + filename + "<span class=\"tab-closer\">X</span></a></li>");
    jQuery(".tab-nav").append(tabLink);

    /* Tab content. */
    jQuery("#tab-content-wrapper").append("<li class=\"tab-content\"><div class=\"file-content\">" + data + "</div></li>");

    /* Record the presence of a tab for this file. */
    addFileTabMapping(filepathRelativeToDocroot, tabLink);

    /* Activate tab. */
    activateTabForFile(filepathRelativeToDocroot);
  });
}

/**
 * Note down a file and its associated tab element.
 *
 * @param string filepath
 * @param object tabElement
 *    jQuery object for a tab.
 */
function addFileTabMapping(filepath, tabElement) {

  fileTabMapping[filepath] = tabElement;
}

/**
 * Remove association record between a file and its tab element.
 *
 * @param string filepath
 */
function removeFileTabMapping(filepath) {

  if (!hasFileTabMapping(filepath)) {
    return;
  }

  delete fileTabMapping[filepath];
}

/**
 * Is there a tab for the given file?
 *
 * Returns the tab element when a tab is present.
 */
function hasFileTabMapping(filepath) {

  if (fileTabMapping.hasOwnProperty(filepath)) {
    return fileTabMapping[filepath];
  }

  return false;
}

/**
 * Open the tab for the given file.
 */
function activateTabForFile(filepath) {

  if (!hasFileTabMapping(filepath)) {
    return;
  }

  fileTabMapping[filepath].click();
}

/**
 * Manage tab closing.
 *
 * - Close the tab.
 * - Remove association between a file and its tab element.
 * - When the current tab is closed, go back to the file browser.
 */
function removeTabForFile(filepath) {

  var tabElement;

  if (!(tabElement = hasFileTabMapping(filepath))) {
    return;
  }

  /* Delete tab and its content */
  var tabIndex = tabElement.index();
  jQuery(".tab-content").get(tabIndex).remove();
  tabElement.remove();

  removeFileTabMapping(filepath);

  /* When we are closing the active tab, return to file browser in first tab. */
  var isActiveTab = tabElement.hasClass("uk-active");
  if (isActiveTab) {
    jQuery(".tab-selector").get(0).click();
  }
}
