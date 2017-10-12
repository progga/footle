/**
 * @file
 * Configuration management.
 */

package config

import (
	"flag"
	"strconv"
)

/**
 * The configuration object.
 */
var config Config

/**
 * Prepares and returns the configuration object.
 *
 * Attempts to implement the Singleton pattern i.e. it always returns a single
 * instance of the configuration object.
 */
func Get() Config {

	if config.args != nil {
		return config
	}

	// Initialize configuration.
	config.args = make(map[string]string)
	config.flags = make(map[string]bool)

	// Now load the values.
	docroot, port, hasCmdLine, hasHTTP := getFlagsAndArgs()

	config.args["docroot"] = docroot
	config.args["port"] = strconv.Itoa(port)
	config.flags["has-cmdline"] = hasCmdLine
	config.flags["has-http"] = hasHTTP

	return config
}

/**
 * Setup command line flags and arguments.
 *
 * Return the values of these flags and arguments.
 *
 * Arg:
 *  - docroot : Docroot of code that will be debugged.
 *  - port: Network port of the HTTP interface.
 *
 * Flag:
 *  - cmdline: We want the command line.
 *  - nohttp : No HTTP.
 */
func getFlagsAndArgs() (docroot string, port int, hasCmdLine, hasHTTP bool) {

	docrootArg := flag.String("docroot", "", "Path of directory whose code you want to debug; e.g. /var/www/html/")
	portArg := flag.Int("port", 9090, "Network port for Footle's Web interface.")
	hasCmdLineFlag := flag.Bool("cmdline", false, "Launch command line debugger.")
	noHTTPFlag := flag.Bool("nohttp", false, "Do *not* launch HTTP interface of the debugger.")

	flag.Parse()

	docroot = *docrootArg
	port = *portArg
	hasCmdLine = *hasCmdLineFlag
	hasHTTP = !*noHTTPFlag

	return docroot, port, hasCmdLine, hasHTTP
}
