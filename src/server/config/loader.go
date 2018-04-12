/**
 * @file
 * Configuration management.
 */

package config

import (
	"flag"
	"log"
	"os"
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

	// Now load the configuration passed from the command line.
	docroot, remoteDocroot, verbosity, httpPort, DBGpPort, hasCmdLine, hasHTTP := getFlagsAndArgs()

	config.SetArg("docroot", docroot)
	config.SetArg("remote-docroot", remoteDocroot)
	config.SetArg("http-port", strconv.Itoa(httpPort))
	config.SetArg("dbgp-port", strconv.Itoa(DBGpPort))
	config.SetArg("verbosity", verbosity)

	if hasCmdLine {
		config.SetFlag("has-cmdline")
	} else {
		config.UnsetFlag("has-cmdline")
	}

	if hasHTTP {
		config.SetFlag("has-http")
	} else {
		config.UnsetFlag("has-http")
	}

	return config
}

/**
 * Setup command line flags and arguments.
 *
 * Return the values of these flags and arguments.
 *
 * Arg:
 *  - docroot : Docroot of code that will be debugged.
 *  - HTTP port: Network port of the HTTP interface.
 *  - DBGp port: Network port to listen for the DBGp server.
 *
 * Flag:
 *  - cmdline: We want the command line.
 *  - nohttp : No HTTP.
 *  - v, vv, vvv: Verbosity level.
 */
func getFlagsAndArgs() (docroot, remoteDocroot, verbosity string, httpPort, DBGpPort int, hasCmdLine, hasHTTP bool) {

	docrootArg := flag.String("docroot", "", "[Optional] Path of directory whose code you want to debug; e.g. /var/www/html/ (default is current dir)")
	remoteDocrootArg := flag.String("docroot-remote", "", "[Optional] When Footle and the DBGp server (e.g. xdebug) are in different machines, this is the path of the source code directory in the remote machine.  This scenario is *not* recommended.  Try as a last resort.")
	DBGpPortArg := flag.Int("port-dbgp", 9000, "[Optional] Network port to listen for the DBGp server.")
	httpPortArg := flag.Int("port-http", 9090, "[Optional] Network port for Footle's Web interface.")
	hasCmdLineFlag := flag.Bool("cmdline", false, "[Optional] Launch command line debugger.")
	noHTTPFlag := flag.Bool("nohttp", false, "[Optional] Do *not* launch HTTP interface of the debugger.")

	LowVerbosityFlag := flag.Bool("v", false, "[Optional] Low verbosity.  Unused.")
	MediumVerbosityFlag := flag.Bool("vv", false, "[Optional] Medium verbosity.  Unused.")
	HighVerbosityFlag := flag.Bool("vvv", false, "[Optional] High verbosity.  Include communication with DBGp server.")

	flag.Parse()

	docroot = *docrootArg
	remoteDocroot = *remoteDocrootArg
	httpPort = *httpPortArg
	DBGpPort = *DBGpPortArg
	hasCmdLine = *hasCmdLineFlag
	hasHTTP = !*noHTTPFlag

	if docroot == "" {
		currentDir, err := os.Getwd()

		if err == nil {
			docroot = currentDir
		}
	} else if _, err := os.Stat(docroot); os.IsNotExist(err) {
		log.Fatal(err)
	}

	if *HighVerbosityFlag {
		verbosity = "high"
	} else if *MediumVerbosityFlag {
		verbosity = "medium"
	} else if *LowVerbosityFlag {
		verbosity = "low"
	}

	return
}
