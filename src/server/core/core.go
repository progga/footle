
package core

import (
  "bytes"
  "fmt"
  "io"
  "log"
  "net"
  "../dbgp/message"
)

/**
 * Receive DBGP messages and pass it to all the user interfaces.
 */
func BroadcastMsgToUI(msg message.Message, toCmdLine, toHTTP chan<- message.Message) {

  toCmdLine <- msg
  toHTTP    <- msg
}

/**
 *
 *
 * Xdebug message format: Int Null XML-Snippet NULL
 */
func ReadMsgFromDBGpEngine(connection net.Conn) (msg string, err error) {

  var dbgpMsg bytes.Buffer
  var msgSize int

  count, err := fmt.Fscanf(connection, "%d\x00", &msgSize)
  if nil != err {
    log.Fatal(err)
  }
  if 0 == count {
    return msg, err
  }

  copy_count, err := io.CopyN(&dbgpMsg, connection, int64(msgSize))
  if 0 == copy_count || io.EOF == err {
    return msg, err
  } else if nil != err {
    log.Fatal(err)
  }

  count, err = fmt.Fscanf(connection, "\x00") // Read null byte.
  if nil != err {
    log.Fatal(err)
  }

  msg = dbgpMsg.String()
  return msg, err
}

/**
 *
 */
func ListenForDBGpEngine() (sock net.Listener) {

  sock, err := net.Listen("tcp", "127.0.0.1:9000");
  if  nil != err {
    log.Fatal(err)
  }

  return sock
}

/**
 *
 */
func StartTalkingToDBGpEngine(sock net.Listener) (connection net.Conn) {

  connection, err := sock.Accept();
  if nil != err {
    log.Fatal(err)
  }

  return connection
}

/**
 *
 */
func SendCmdsToDBGpEngine(conn *net.Conn, in <-chan string) {

  for DBGpCmd := range in {
    connection := *conn

    if isActiveConnection(connection) {
      writeCount, err := connection.Write([]byte(DBGpCmd))
      _ = writeCount

      if nil != err {
        log.Fatal(err)
      }
    } else {
      fmt.Println("Inactive connection.")
    }
  }
}

/**
 *
 */
func isActiveConnection(connection net.Conn) bool {

  ignore := []byte{}

  if nil == connection {
    return false
  }

  if readCount, err := connection.Write(ignore); nil != err {
    _ = readCount
    return false
  }

  return true
}
