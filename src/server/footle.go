/**
 * Footle the DBGp debugger.
 *
 * Here we launch go routines for command line UI, HTTP UI, receiving
 * messages from DBGp engine and sending DBGp commands to DBGp engine.
 */

package main

import (
	"server/cmdline"
	"server/config"
	"server/core"
	conn "server/core/connection"
	"server/dbgp/message"
	"server/http"
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
	config := config.Get()

	// Initializations.
	var MsgsForCmdLineUI, MsgsForHTTPUI chan message.Message

	CmdsFromUI := make(chan string)
	DBGpCmds := make(chan string)
	bye := make(chan struct{})

	// Launch all interfaces.
	if config.HasCmdLine() {
		MsgsForCmdLineUI = make(chan message.Message)

		go cmdline.RunUI(CmdsFromUI, bye)
		go cmdline.UpdateUIStatus(MsgsForCmdLineUI)
	}

	if config.HasHTTP() {
		MsgsForHTTPUI = make(chan message.Message)

		go http.Listen(CmdsFromUI, config)
		go http.TellBrowsers(MsgsForHTTPUI, config)
	}

	// Talk to DBGp engine.
	DBGpConnection := conn.GetConnection()
	DBGpConnection.Activate()

	go core.RecvMsgsFromDBGpEngine(DBGpConnection, MsgsForCmdLineUI, MsgsForHTTPUI)
	go core.SendCmdsToDBGpEngine(DBGpConnection, DBGpCmds)

	// Let Footle deal with all commands from UIs first.  Some commands will then
	// head for the DBGp engine while some will change Footle's internal state.
	go core.ProcessUICmds(CmdsFromUI, DBGpCmds, DBGpConnection)

	<-bye
}
