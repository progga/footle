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
func (c Config) GetDocroot() string {

	return c.GetArg("docroot")
}

/**
 * Getter for network port where Footle listens for HTTP requests.
 */
func (c Config) GetHTTPPort() int {

	port := c.GetArg("port")
	portNumber, err := strconv.Atoi(port)

	if err != nil {
		return -1
	}

	return portNumber
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
