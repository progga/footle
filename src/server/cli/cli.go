/**
 * Command line user interface for a DBGp debugger.
 */

package cli

import (
	"encoding/base64"
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"log"
	"server/config"
	"server/dbgp/command"
	"server/dbgp/message"
)

const READLINE_PROMPT = "> "

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

	rl, err := readline.New(READLINE_PROMPT)
	if err != nil {
		log.Fatal(err)
	}

	config := config.Get()

	for {
		cmd, err := rl.Readline()
		if err != nil && err != io.EOF {
			log.Println(err)
			continue
		} else if err == io.EOF {
			close(bye)
			return
		}

		cmdAlias, cmdArgs, err := command.Break(cmd)
		if nil != err {
			log.Println(err)
			continue
		}

		if "bye" == cmd || "quit" == cmd || "q" == cmd {
			close(bye)
			return
		} else if "refresh" == cmd || "" == cmd {
			continue
		} else if cmd == "verbose" {
			config.GoVerbose()
			continue
		} else if cmd == "no-verbose" {
			config.GoSilent()
			continue
		} else if cmd == "on" || cmd == "off" || cmd == "continue" || cmdAlias == "update_source" {
			// Commands for controlling Footle.
			out <- cmd
			continue
		}

		err = command.Validate(cmdAlias, cmdArgs)
		if nil != err {
			log.Println(err)
			continue
		}

		out <- cmd
	}
}

/**
 * Display incoming DBGP messages.
 */
func UpdateUIStatus(in <-chan message.Message) {

	for msg := range in {
		fmt.Printf("%s\n\r%s", msg, READLINE_PROMPT)

		// Some commands such as "source" send XML character data
		// as inner XML content.
		decoded, err := base64.StdEncoding.DecodeString(msg.Content)
		if nil == err && 0 < len(decoded) {
			fmt.Printf("%s\n\r%s", string(decoded), READLINE_PROMPT)
		}
	}
}
