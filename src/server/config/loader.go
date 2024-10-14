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
	codebase, remoteCodebase, uiPath, verbosity, httpPort, DBGpPort, hasCmdLine, hasHTTP := getFlagsAndArgs()

	config.SetArg("codebase", codebase)
	config.SetArg("remote-codebase", remoteCodebase)
	config.SetArg("http-port", strconv.Itoa(httpPort))
	config.SetArg("dbgp-port", strconv.Itoa(DBGpPort))
	config.SetArg("ui-path", uiPath)
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
 *  - codebase: Parent directory of code that will be debugged.
 *  - Remote codebase: This is what the debugger engine sees when the engine is
 *    in a remote machine.
 *  - HTTP port: Network port of the HTTP interface.
 *  - DBGp port: Network port to listen for the DBGp server.
 *  - UI path: Location of the HTTP UI.
 *
 * Flag:
 *  - cli: We want the command line.
 *  - nohttp : No HTTP.
 *  - v, vv, vvv: Verbosity level.
 */
func getFlagsAndArgs() (codebase, remoteCodebase, uiPath, verbosity string, httpPort, DBGpPort int, hasCmdLine, hasHTTP bool) {

	codebaseArg := flag.String("codebase", "", "[Optional] Path of directory whose code you want to debug; e.g. /var/www/html/ (default is current dir)")
	remoteCodebaseArg := flag.String("codebase-remote", "", "[Optional] When Footle and the DBGp server (e.g. xdebug) are in different machines, this is the path of the source code directory in the remote machine.  This scenario is *not* recommended.  Try as a last resort.  Footle assumes that a copy of the source code is present in the local machine.  To tell Footle where this local copy is, either run footle from inside that copy or use the -codebase option.")
	DBGpPortArg := flag.Int("port-dbgp", 9003, "[Optional] Network port to listen for the DBGp server.")
	httpPortArg := flag.Int("port-http", 1234, "[Optional] Network port for Footle's Web interface.")
	uiPathArg := flag.String("ui-path", "", "[Optional] Location of an alternate HTTP UI.  Only relevant during UI development.")

	hasCmdLineFlag := flag.Bool("cli", false, "[Optional] Launch command line debugger.")
	noHTTPFlag := flag.Bool("nohttp", false, "[Optional] Do *not* launch HTTP interface of the debugger.")

	LowVerbosityFlag := flag.Bool("v", false, "[Optional] Low verbosity.  Unused.")
	MediumVerbosityFlag := flag.Bool("vv", false, "[Optional] Medium verbosity.  Unused.")
	HighVerbosityFlag := flag.Bool("vvv", false, "[Optional] High verbosity.  Include communication with DBGp server.")

	flag.Parse()

	codebase = *codebaseArg
	remoteCodebase = *remoteCodebaseArg
	httpPort = *httpPortArg
	DBGpPort = *DBGpPortArg
	uiPath = *uiPathArg
	hasCmdLine = *hasCmdLineFlag
	hasHTTP = !*noHTTPFlag

	if codebase == "" {
		currentDir, err := os.Getwd()

		if err == nil {
			codebase = currentDir
		}
	} else if _, err := os.Stat(codebase); os.IsNotExist(err) {
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
