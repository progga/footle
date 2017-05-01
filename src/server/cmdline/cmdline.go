
package cmdline

import (
  "fmt"
  "log"
  "strings"
  "github.com/chzyer/readline"
  "../dbgp/command"
  "../dbgp/message"
)

const space = " "

var incomingMsgQueue []message.Message

/**
 *
 */
func RunUI(out chan<- string) {

  rl, err := readline.New("> ")
  if nil != err {
    log.Fatal(err)
  }

  for {
    if len(incomingMsgQueue) > 0 {
      // First, display all messages stored in the queue.
      for _, msg := range incomingMsgQueue {
        fmt.Println(msg)
      }
      incomingMsgQueue = incomingMsgQueue[:0]
    }

    cmd, err := rl.Readline()
    if nil != err {
      fmt.Println(err)
      continue
    }

    cmd = strings.TrimSpace(cmd)

    if "bye" == cmd {
      return
    } else if "refresh" == cmd || "" == cmd {
      continue
    }

    cmdParts := strings.Split(cmd, space)
    if len(cmdParts) < 1 {
      continue
    }

    footleCmd := cmdParts[0]
    cmdArgs   := cmdParts[1:]
    if err = Validate(footleCmd, cmdArgs); nil != err {
      fmt.Println(err)
      continue
    }

    DBGpCmd, err := command.Prepare(footleCmd, cmdArgs)
    if nil != err {
      fmt.Println(err)
      continue
    }

    out <- DBGpCmd
  }
}

/**
 * Second go routine for storing incoming DBGP messages.
 */
func UpdateUIStatus(in <-chan message.Message) {

  for msg := range in {
    // @todo Lock the queue before updating.
    incomingMsgQueue = append(incomingMsgQueue, msg)
  }
}
