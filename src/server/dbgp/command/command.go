/**
 * Command pkg.  Turns short commands into full blown DBGp commands.
 */

package command

import (
	"fmt"
	"strings"
)

const space = " "

/**
 * Breaks a given string into a command and its arguments.
 *
 * Example:
 *    Input: foo bar qux
 *
 *    Output:
 *      Command: foo
 *      args:
 *        - bar
 *        - qux
 */
func Break(cmd string) (shortCmd string, cmdArgs []string, err error) {

	trimmedCmd := strings.TrimSpace(cmd)

	cmdParts := strings.Split(trimmedCmd, space)

	if len(cmdParts) < 1 {
		err = fmt.Errorf("Cannot break short command %s", cmd)

		return shortCmd, cmdArgs, err
	}

	shortCmd = cmdParts[0]
	cmdArgs = cmdParts[1:]

	return shortCmd, cmdArgs, err
}

/**
 * Turns short commands into full blown DBGp commands.
 *
 * Example of short command: r
 * Example of corresponding DBGp command: run -i 59
 */
func Prepare(shortCmd string, cmdArgs []string) (DBGpCmd string, err error) {

	if err = Validate(shortCmd, cmdArgs); nil != err {
		return DBGpCmd, err
	}

	DBGpCmd, err = PrepareDBGpCmd(shortCmd, cmdArgs)

	return DBGpCmd, err
}

/**
 * Given a DBGp command, extract the command name.
 *
 * For example, if the UI has received "b foo.php 2", which is a shorthand for
 * "breakpoint_set foo.php 2", return "breakpoint_set".
 */
func Extract(cmdFromUI string) (DBGpCmdName string, err error) {

	cmdName, cmdArgs, err := Break(cmdFromUI)
	if err != nil {
		return DBGpCmdName, err
	}

	err = Validate(cmdName, cmdArgs)
	if err != nil {
		return DBGpCmdName, err
	}

	DBGpCmdName = resolveAlias(cmdName)

	return DBGpCmdName, err
}
