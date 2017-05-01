
package main

import (
  "net"
  "./cmdline"
  "./core"
  "./http"
  "./dbgp/message"
)

func main() {

  var activeDBGpConnection net.Conn

  MsgsForCmdLineUI  := make(chan message.Message)
  MsgsForHTTPUI     := make(chan message.Message)
  CmdsFromUI        := make(chan string)

  go cmdline.RunUI(CmdsFromUI)
  go cmdline.UpdateUIStatus(MsgsForCmdLineUI)

  go http.RunUI(CmdsFromUI)
  go http.UpdateUIStatus(MsgsForHTTPUI)

  sock := core.ListenForDBGpEngine()
  go core.SendCmdsToDBGpEngine(&activeDBGpConnection, CmdsFromUI)

  for {
    activeDBGpConnection = core.StartTalkingToDBGpEngine(sock)

    for {
      msg, err := core.ReadMsgFromDBGpEngine(activeDBGpConnection)
      if len(msg) == 0 || nil != err {
        break
      }

      if parsedMsg, err := message.Decode(msg); nil == err {
        core.BroadcastMsgToUI(parsedMsg, MsgsForCmdLineUI, MsgsForHTTPUI)
      }
    }

    activeDBGpConnection.Close()
  }
}
