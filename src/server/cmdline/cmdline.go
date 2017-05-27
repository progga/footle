/**
 * Command line user interface for a DBGp debugger.
 */

package cmdline

import (
	"../dbgp/command"
	"../dbgp/message"
	"encoding/base64"
	"fmt"
	"github.com/chzyer/readline"
	"log"
)

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

				// Some commands such as "source" send XML character data
				// as inner XML content.
				decoded, err := base64.StdEncoding.DecodeString(msg.Content)
				if nil == err && 0 < len(decoded) {
					fmt.Printf("%s", string(decoded))
				}
			}
			incomingMsgQueue = incomingMsgQueue[:0]
		}

		cmd, err := rl.Readline()
		if nil != err {
			fmt.Println(err)
			continue
		}

		shortCmd, cmdArgs, err := command.Break(cmd)
		if nil != err {
			fmt.Println(err)
			continue
		}

		if "bye" == shortCmd || "quit" == shortCmd || "q" == shortCmd {
			close(bye)

			return
		} else if "refresh" == cmd || "" == cmd {
			continue
		}

		DBGpCmd, err := command.Prepare(shortCmd, cmdArgs)
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
