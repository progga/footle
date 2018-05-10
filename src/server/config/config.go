/**
 * @file
 * Definition for the config class.
 */

package config

import "strconv"

/**
 * The Config class stores command line arguments and flags.
 */
type Config struct {
	args  map[string]string
	flags map[string]bool
}

/**
 * Getter for path of source code to debug.
 */
func (c Config) GetCodebase() string {

	return c.GetArg("codebase")
}

/**
 * Getter for remote path of source code to debug.
 *
 * Only relevant when Footle and the DBGp engine are running in different
 * machines.  In this type of setups, the file paths returned by the DBGp engine
 * will start with the remote path.
 *
 * Example:
 *   Codebase in Footle's machine: /home/foo/bar/
 *   Codebase in DBGp engine's machine: /var/www/html/
 *   In this case /var/www/html/ is the remote codebase.
 */
func (c Config) GetRemoteCodebase() string {

	return c.GetArg("remote-codebase")
}

/**
 * Determine the source code path returned by the DBGp engine.
 *
 * When the DBGp engine and Footle are in different machines, source code paths
 * returned by the DBGp engine will start with a path from that machine.  This
 * path is likely to be different from local paths seen by Footle.
 */
func (c Config) DetermineCodeDir() (codeDir string) {

	codeDir = c.GetRemoteCodebase()

	if codeDir == "" {
		codeDir = c.GetCodebase()
	}

	return codeDir
}

/**
 * Getter for network port where Footle listens for HTTP requests.
 */
func (c Config) GetHTTPPort() int {

	return c.getInt("http-port")
}

/**
 * Getter for network port to listen for DBGp server.
 */
func (c Config) GetDBGpPort() int {

	return c.getInt("dbgp-port")
}

/**
 * Getter for alternate HTTP UI path.
 */
func (c Config) GetUIPath() string {

	return c.GetArg("ui-path")
}

/**
 * Return value of configuration item as an integer.
 *
 * Useful for fetching numeric configurations such as port number.
 */
func (c Config) getInt(item string) int {

	value := c.GetArg(item)
	number, err := strconv.Atoi(value)

	if err != nil {
		return -1
	}

	return number
}

/**
 * Getter for the presence of the command line interface.
 */
func (c Config) HasCmdLine() bool {

	return c.GetFlag("has-cmdline")
}

/**
 * Getter for the presence of the HTTP interface.
 */
func (c Config) HasHTTP() bool {

	return c.GetFlag("has-http")
}

/**
 * Predicate for determine verbosity.
 *
 * The "high" verbose mode displays all incoming and outgoing communication with
 * DBGp server.  The "low" and "medium" verbose modes are unused at the moment.
 *
 * By default, there is no verbosity.
 */
func (c Config) IsVerbose() bool {

	verbosityLevel := c.GetArg("verbosity")

	return (verbosityLevel == "high")
}

/**
 * Turn on verbose mode.
 */
func (c Config) GoVerbose() {

	c.SetArg("verbosity", "high")
}

/**
 * Turn off verbose mode.
 */
func (c Config) GoSilent() {

	c.SetArg("verbosity", "")
}

/**
 * Get any command line argument.
 */
func (c Config) GetArg(item string) string {

	value, itemPresent := c.args[item]

	if itemPresent {
		return value
	} else {
		return ""
	}
}

/**
 * Get any command line flag.
 */
func (c Config) GetFlag(item string) bool {

	value, itemPresent := c.flags[item]

	if itemPresent {
		return value
	} else {
		return false
	}
}

/**
 * Set flag to true.
 */
func (c Config) SetFlag(item string) {

	c.flags[item] = true
}

/**
 * Set flag to false.
 */
func (c Config) UnsetFlag(item string) {

	c.flags[item] = false
}

/**
 * Setter for argument.
 */
func (c Config) SetArg(item string, value string) {

	c.args[item] = value
}
