
/**
 * @file
 * Manage the recently opened file list.
 */

/**
 * Class for managing the recently opened file list.
 */
class RecentFiles {
  /**
   * As it says on the tin.
   */
  constructor (storage, recentFileCount = 5) {
    this.storage = storage
    this.fileCount = recentFileCount
  }

  /**
   * Fetch the existing list.
   */
  get () {
    const word = this.storage.getItem('recent-files')

    const recentFiles = word ? word.split(',') : []

    return recentFiles
  }

  /**
   * Add the given filename to the saved list.
   */
  add (filename) {
    const recentFiles = this.get()

    const moreRecentFiles = this.update(recentFiles, filename)

    const hasChanged = !((recentFiles.length === moreRecentFiles) && moreRecentFiles.every((item, index) => item === recentFiles[index]))
    if (hasChanged) {
      this.storage.setItem('recent-files', moreRecentFiles)
    }
  }

  /**
   * Add the given filename to a unique list of maximum five files.
   */
  update (recentFiles, filename) {
    const updated = [...new Set([filename, ...recentFiles])].slice(0, this.fileCount)

    return updated
  }
}

export default RecentFiles
