
/**
 * Functions for receiving messages from DBGp engine.
 */

package core

import (
  "bytes"
  "fmt"
  "log"
  "io"
  "net"
  "../dbgp/message"
)

/**
 * Receive message from DBGp engine and pass it to user interfaces.
 */
func RecvMsgsFromDBGpEngine(sock net.Listener, activeDBGpConnection *net.Conn, MsgsForCmdLineUI, MsgsForHTTPUI chan<- message.Message) {

  for {
    *activeDBGpConnection = StartTalkingToDBGpEngine(sock)

    for {
      msg, err := ReadMsgFromDBGpEngine(*activeDBGpConnection)
      if len(msg) == 0 || nil != err {
        break
      }

      if parsedMsg, err := message.Decode(msg); nil == err {
        BroadcastMsgToUIs(parsedMsg, MsgsForCmdLineUI, MsgsForHTTPUI)
      }
    }

    (*activeDBGpConnection).Close()
  }
}

/**
 * Establish connection with a DBGp engine.
 */
func StartTalkingToDBGpEngine(sock net.Listener) (connection net.Conn) {

  connection, err := sock.Accept();
  if nil != err {
    log.Fatal(err)
  }

  return connection
}

/**
 * Read message from DBGp engine over a network connection.
 *
 * Xdebug message format: Int Null XML-Snippet NULL
 */
func ReadMsgFromDBGpEngine(connection net.Conn) (msg string, err error) {

  var DBGpMsg bytes.Buffer
  var msgSize int

  count, err := fmt.Fscanf(connection, "%d\x00", &msgSize)
  if nil != err {
    log.Fatal(err)
  }
  if 0 == count {
    return msg, err
  }

  copyCount, err := io.CopyN(&DBGpMsg, connection, int64(msgSize))
  if 0 == copyCount || io.EOF == err {
    return msg, err
  } else if nil != err {
    log.Fatal(err)
  }

  count, err = fmt.Fscanf(connection, "\x00") // Read null byte.
  if nil != err {
    log.Fatal(err)
  }

  msg = DBGpMsg.String()
  return msg, err
}

/**
 * Pass on a DBGP message to all the user interfaces.
 *
 * User interfaces include the command line interface and the HTTP interface.
 */
func BroadcastMsgToUIs(msg message.Message, toCmdLine, toHTTP chan<- message.Message) {

  toCmdLine <- msg
  toHTTP    <- msg
}
