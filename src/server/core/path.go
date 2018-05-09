/**
 * @file
 * Filepath and URI related functions.
 */

package core

import (
	"path/filepath"
	"server/config"
	"strings"
)

/**
 * Turn a relative filepath into an absolute path.
 *
 *   - foo/bar.txt -> /codebase/foo/bar.txt
 *   - /foo/bar.txt -> /foo/bar.txt
 *   - file:///codebase/foo/bar.txt -> /codebase/foo/bar.txt
 */
func toAbsolutePath(relativePath string, config config.Config) (absolutePath string) {

	fileUri := toAbsoluteUri(relativePath, config)
	absoluteFilename := strings.TrimPrefix(fileUri, "file://")

	return absoluteFilename
}

/**
 * Turn a relative filepath into an absolute URI.
 *
 * Examples:
 *   - foo/bar.txt -> file://codebase/foo/bar.txt
 *   - /foo/bar.txt -> file:///foo/bar.txt
 *   - file://codebase/foo/bar.txt -> file://codebase/foo/bar.txt
 *
 * @todo Add Unit tests.
 */
func toAbsoluteUri(relativePath string, config config.Config) (absoluteUri string) {

	isAbsoluteUri := strings.HasPrefix(relativePath, "file://")
	if isAbsoluteUri {
		absoluteUri = relativePath
		return absoluteUri
	}

	isAbsolutePath := filepath.IsAbs(relativePath)
	if isAbsolutePath {
		absoluteUri = "file://" + relativePath

		return absoluteUri
	}

	codebase := config.DetermineCodeDir()
	absoluteUri = "file://" + filepath.Join(codebase, relativePath)

	return absoluteUri
}
