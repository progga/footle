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

	return c.getArg("docroot")
}

/**
 * Getter for network port where Footle listens for HTTP requests.
 */
func (c Config) GetHTTPPort() int {

	port := c.getArg("port")
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

	return c.getFlag("has-http")
}

/**
 * Getter for the presence of the HTTP interface.
 */
func (c Config) HasHTTP() bool {

	return c.getFlag("has-http")
}

/**
 * Get any command line argument.
 */
func (c Config) getArg(item string) string {

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
func (c Config) getFlag(item string) bool {

	value, itemPresent := c.flags[item]

	if itemPresent {
		return value
	} else {
		return false
	}
}
