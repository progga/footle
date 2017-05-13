
/**
 * Command line user interface for a DBGp debugger.
 */

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

/**
 * Stored DBGp messages.
 *
 * Storage for any DBGp message that has arrived after the previous command has
 * been issued from the command line.  Stored messages are usually the result of
 * the previous command.
 *
 * This is needed because we do not *block* the command line once a command has
 * been issued.  The response for the previous command is displayed when the
 * next command is issued.
 *
 * @see UpdateUIStatus()
 */
var incomingMsgQueue []message.Message

/**
 * Execute the command line interface.
 *
 * Read commands, validate them, turn them into DBGp commands, and pass them
 * to the other part of the debugger that deals with talking to the DBGp
 * engine.
 *
 * Special commands:
 *   - The "bye" command exits the debugger.
 *   - The empty string or "refresh" command lists any DBGp message that has
 *     arrived after the previous command has been issued.
 *
 * @param chan<- string out
 *   DBGp commands are written to this channel.
 * @param chan struct{} bye
 *   Event channel.  It is closed to broadcast the global exit event.
 */
func RunUI(out chan<- string, bye chan struct{}) {

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
      close(bye)

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
    if err = command.Validate(footleCmd, cmdArgs); nil != err {
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
 * Store incoming DBGP messages.
 *
 * The stored messages can be viewed in the command line interface.
 */
func UpdateUIStatus(in <-chan message.Message) {

  for msg := range in {
    // @todo Lock the queue before updating.
    incomingMsgQueue = append(incomingMsgQueue, msg)
  }
}
