/**
 * Footle the DBGp debugger.
 *
 * Here we launch go routines for command line UI, HTTP UI, receiving
 * messages from DBGp engine and sending DBGp commands to DBGp engine.
 */

package main

import (
	"server/cli"
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
	DBGpMessages := make(chan message.Message)
	bye := make(chan struct{})

	launchUIs(config, &MsgsForCmdLineUI, &MsgsForHTTPUI, CmdsFromUI, bye)

	// Talk to DBGp engine.
	DBGpConnection := conn.GetConnection()
	DBGpConnection.Activate()

	go core.RecvMsgsFromDBGpEngine(DBGpConnection, DBGpMessages)
	go core.SendCmdsToDBGpEngine(DBGpConnection, DBGpCmds)

	// Let Footle deal with all commands from UIs first.  Some commands will then
	// head for the DBGp engine while some will change Footle's internal state.
	go core.ProcessUICmds(CmdsFromUI, DBGpCmds, DBGpMessages, DBGpConnection)

	// Process incoming DBGP messages before selectively passing them to the UIs.
	go core.ProcessDBGpMessages(DBGpCmds, DBGpMessages, MsgsForCmdLineUI, MsgsForHTTPUI)

	<-bye
}

/**
 * Launch all user interfaces.
 *
 * Start the HTTP and/or the Cli interfaces depending on user preferences.
 */
func launchUIs(config config.Config, MsgsForCmdLineUI, MsgsForHTTPUI *chan message.Message, CmdsFromUI chan string, bye chan struct{}) {

	if config.HasCmdLine() {
		*MsgsForCmdLineUI = make(chan message.Message)

		go cli.RunUI(CmdsFromUI, bye)
		go cli.UpdateUIStatus(*MsgsForCmdLineUI)
	}

	if config.HasHTTP() {
		*MsgsForHTTPUI = make(chan message.Message)

		go http.Listen(CmdsFromUI, config)
		go http.TellBrowsers(*MsgsForHTTPUI, config)
	}
}
