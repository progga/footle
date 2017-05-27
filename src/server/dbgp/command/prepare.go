/**
 * Functions for DBGp command preparation.
 */

package command

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

/**
 * DBGp transaction ID.  Its value was last used for the DBGp command's
 * transaction ID.
 *
 * @see fetchNextTxId()
 */
var lastTxId int

/**
 * Prepare DBGp command from the given values.
 *
 * The arguments for this function can be considered a short form of the DBGp
 * commands.  Here the short forms are expanded to their full forms.
 *
 * Example: "b /home/foo/php/bar.php 9" becomes
 * "breakpoint_set -i 1 -t line -f /home/foo/php/bar.php -n 9"
 */
func PrepareDBGpCmd(cmd string, args []string) (DBGpCmd string, err error) {

	TxId := fetchNextTxId()

	switch strings.ToLower(cmd) {
	case "breakpoint", "b":
		DBGpCmd, err = prepareBreakpointCmd(args, TxId)

	case "eval", "ev":
		DBGpCmd, err = prepareEvalCmd(args, TxId)

	case "run", "r":
		DBGpCmd, err = prepareCmdNoArgs("run", TxId)

	case "source", "src", "sr":
		DBGpCmd, err = prepareSourceCmd(args, TxId)

	case "status", "s":
		DBGpCmd, err = prepareCmdNoArgs("status", TxId)

	case "stop", "st":
		DBGpCmd, err = prepareCmdNoArgs("stop", TxId)

	case "step_into", "si":
		DBGpCmd, err = prepareCmdNoArgs("step_into", TxId)

	case "step_out", "so":
		DBGpCmd, err = prepareCmdNoArgs("step_out", TxId)

	case "step_over", "sov", "sv":
		DBGpCmd, err = prepareCmdNoArgs("step_over", TxId)

	default:
		DBGpCmd, err = "", fmt.Errorf("Unknown command: %s", cmd)
	}

	return DBGpCmd, err
}

/**
 * Determine the transaction ID for the next DBGp command.
 *
 * @return int
 */
func fetchNextTxId() (nextTxId int) {

	lastTxId++
	nextTxId = lastTxId % math.MaxInt32

	lastTxId = nextTxId

	return nextTxId
}

/**
 * The DBGp Breakpoint command.
 */
func prepareBreakpointCmd(args []string, TxId int) (DBGpCmd string, err error) {

	if 2 > len(args) {
		return DBGpCmd, fmt.Errorf("Need at least two args for preparing breakpoint cmd.")
	}

	filepath := args[0]
	lineNumber := args[1]

	DBGpCmd = fmt.Sprintf("breakpoint_set -i %d -t line -f %s -n %s\x00", TxId, filepath, lineNumber)

	return DBGpCmd, err
}

/**
 * DBGp Eval command.
 */
func prepareEvalCmd(args []string, TxId int) (DBGpCmd string, err error) {

	if 0 == len(args) {
		return DBGpCmd, fmt.Errorf("Unsufficient number of args for eval.")
	}

	DBGpCmd = fmt.Sprintf("eval -i %d -- %s\x00", TxId, args[0])

	return DBGpCmd, err
}

/**
 * DBGp Source command.
 */
func prepareSourceCmd(args []string, TxId int) (DBGpCmd string, err error) {

	if 2 != len(args) {
		err = fmt.Errorf("Unsufficient number of args for source.")
		return DBGpCmd, err
	}

	beginLine, err := strconv.ParseInt(args[0], 10, 64)
	lineCount, err := strconv.ParseInt(args[1], 10, 64)
	endLine := beginLine + lineCount

	DBGpCmd = fmt.Sprintf("source -i %d -b %d -e %d\x00", TxId, beginLine, endLine)

	return DBGpCmd, err
}

/**
 * Any DBGp command that does not take any argument other than the TX ID.
 *
 * Example: run, stop, etc.
 */
func prepareCmdNoArgs(cmd string, TxId int) (DBGpCmd string, err error) {

	cmd = strings.TrimSpace(cmd)

	if "" == cmd {
		err = fmt.Errorf("Command cannot be empty.")

		return DBGpCmd, err
	}

	DBGpCmd = fmt.Sprintf("%s -i %d\x00", cmd, TxId)

	return DBGpCmd, err
}
