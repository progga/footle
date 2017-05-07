
/**
 * Footle the DBGp debugger.
 *
 * Here we launch go routines for command line UI, HTTP UI, receiving
 * messages from DBGp engine and sending DBGp command to DBGp engine.
 */

package main

import (
  "net"
  "./cmdline"
  "./core"
  "./http"
  "./dbgp/message"
)

/**
 * Launch go routines.
 *
 * Launch the debugging process and its user interfaces.
 *
 * End execution when the "bye" channel is closed.
 */
func main() {

  var activeDBGpConnection net.Conn

  MsgsForCmdLineUI  := make(chan message.Message)
  MsgsForHTTPUI     := make(chan message.Message)
  CmdsFromUI        := make(chan string)
  bye               := make(chan struct{})

  go cmdline.RunUI(CmdsFromUI, bye)
  go cmdline.UpdateUIStatus(MsgsForCmdLineUI)

  go http.RunUI(CmdsFromUI)
  go http.UpdateUIStatus(MsgsForHTTPUI)

  sock := core.ListenForDBGpEngine()
  go core.RecvMsgsFromDBGpEngine(sock, &activeDBGpConnection, MsgsForCmdLineUI, MsgsForHTTPUI)
  go core.SendCmdsToDBGpEngine(&activeDBGpConnection, CmdsFromUI)

  <- bye
}
