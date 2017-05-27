/**
 * Footle the DBGp debugger.
 *
 * Here we launch go routines for command line UI, HTTP UI, receiving
 * messages from DBGp engine and sending DBGp commands to DBGp engine.
 */

package main

import (
	"./cmdline"
	"./core"
	"./dbgp/message"
	"./http"
	"flag"
	"net"
)

/**
 * Launch go routines.
 *
 * Launch the debugger and its user interfaces.
 *
 * End execution when the "bye" channel is closed.
 */
func main() {

	// Setup command line flags and arguments.
	docroot, hasCmdLine, hasHTTP := getFlagsAndArgs()

	// Initializations.
	var activeDBGpConnection net.Conn

	var MsgsForCmdLineUI, MsgsForHTTPUI chan message.Message

	CmdsFromUI := make(chan string)
	bye := make(chan struct{})

	// Launch all interfaces.
	if hasCmdLine {
		MsgsForCmdLineUI = make(chan message.Message)

		go cmdline.RunUI(CmdsFromUI, bye)
		go cmdline.UpdateUIStatus(MsgsForCmdLineUI)
	}

	if hasHTTP {
		MsgsForHTTPUI = make(chan message.Message)

		go http.Listen(docroot, CmdsFromUI)
		go http.Tell(MsgsForHTTPUI)
	}

	// Talk to DBGp engine.
	sock := core.ListenForDBGpEngine()
	go core.RecvMsgsFromDBGpEngine(sock, &activeDBGpConnection, MsgsForCmdLineUI, MsgsForHTTPUI)
	go core.SendCmdsToDBGpEngine(&activeDBGpConnection, CmdsFromUI)

	<-bye
}

/**
 * Setup command line flags and arguments.
 *
 * Return the values of these flags and arguments.
 *
 * Arg:
 *  - docroot : Docroot of code that will be debugged.
 *
 * Flag:
 *  - cmdline : We want the command line.
 *  - nohttp  : No HTTP.
 */
func getFlagsAndArgs() (docroot string, hasCmdLine, hasHTTP bool) {

	docrootArg := flag.String("docroot", "", "Path of directory whose code you want to debug; e.g. /var/www/html/")
	hasCmdLineFlag := flag.Bool("cmdline", false, "Launch command line debugger.")
	noHTTPFlag := flag.Bool("nohttp", false, "Do *not* launch HTTP interface of the debugger.")

	flag.Parse()

	docroot = *docrootArg
	hasCmdLine = *hasCmdLineFlag
	hasHTTP = !*noHTTPFlag

	return docroot, hasCmdLine, hasHTTP
}
