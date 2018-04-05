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
 * Context ID number for fetching local or global variables.
 *
 * The context Ids for the local and global contexts are used in the context_get
 * command.  These are used to fetch local and global variables.
 */
const localContextId = 0
const globalContextId = 1

/**
 * Context labels available to UIs.
 */
const localContextLabel = "local"
const globalContextLabel = "global"

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

	DBGpCmd = resolveAlias(cmd)

	switch DBGpCmd {
	case "breakpoint_set":
		DBGpCmd, err = prepareBreakpointCmd(args, TxId)

	case "breakpoint_list":
		DBGpCmd, err = prepareCmdNoArgs("breakpoint_list", TxId)

	case "breakpoint_get":
		DBGpCmd, err = prepareBreakpointGetCmd(args, TxId)

	case "breakpoint_remove":
		DBGpCmd, err = prepareBreakpointRemoveCmd(args, TxId)

	case "context_get":
		DBGpCmd, err = prepareContextGetCmd(args, TxId)

	case "eval":
		DBGpCmd, err = prepareEvalCmd(args, TxId)

	case "run":
		DBGpCmd, err = prepareCmdNoArgs("run", TxId)

	case "dbgp":
		DBGpCmd, err = prepareRawDBGpCmd(args, TxId)

	case "property_get":
		DBGpCmd, err = preparePropertyGetCmd(args, TxId)

	case "stack_get":
		DBGpCmd, err = prepareCmdNoArgs("stack_get", TxId)

	case "source":
		DBGpCmd, err = prepareSourceCmd(args, TxId)

	case "status":
		DBGpCmd, err = prepareCmdNoArgs("status", TxId)

	case "stop":
		DBGpCmd, err = prepareCmdNoArgs("stop", TxId)

	case "step_into":
		DBGpCmd, err = prepareCmdNoArgs("step_into", TxId)

	case "step_out":
		DBGpCmd, err = prepareCmdNoArgs("step_out", TxId)

	case "step_over":
		DBGpCmd, err = prepareCmdNoArgs("step_over", TxId)

	default:
		DBGpCmd, err = "", fmt.Errorf("Unknown command: %s", cmd)
	}

	return DBGpCmd, err
}

/**
 * Determine the transaction ID for the next DBGp command.
 *
 * Uses global variable "lastTxId"
 *
 * @todo Make it goroutine safe by wrapping it in a lock.
 */
func fetchNextTxId() (nextTxId int) {

	lastTxId++
	nextTxId = lastTxId % math.MaxInt32

	lastTxId = nextTxId

	return nextTxId
}

/**
 * The DBGp Breakpoint set command.
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
 * The DBGp Breakpoint get command.
 */
func prepareBreakpointGetCmd(args []string, TxId int) (DBGpCmd string, err error) {

	if 1 > len(args) {
		return DBGpCmd, fmt.Errorf("Need at least one argument to prepare the breakpoint_get cmd.")
	}

	breakpointId := args[0]

	DBGpCmd = fmt.Sprintf("breakpoint_get -i %d -d %s\x00", TxId, breakpointId)

	return DBGpCmd, err
}

/**
 * The DBGp breakpoint_remove command.
 */
func prepareBreakpointRemoveCmd(args []string, TxId int) (DBGpCmd string, err error) {

	if 1 > len(args) {
		return DBGpCmd, fmt.Errorf("Need at least one argument to prepare the breakpoint_remove cmd.")
	}

	breakpointId := args[0]

	DBGpCmd = fmt.Sprintf("breakpoint_remove -i %d -d %s\x00", TxId, breakpointId)

	return DBGpCmd, err
}

/**
 * DBGp Eval command.
 */
func prepareEvalCmd(args []string, TxId int) (DBGpCmd string, err error) {

	if 0 == len(args) {
		return DBGpCmd, fmt.Errorf("Insufficient number of args for eval.")
	}

	DBGpCmd = fmt.Sprintf("eval -i %d -- %s\x00", TxId, args[0])

	return DBGpCmd, err
}

/**
 * Any DBGp command.
 *
 * Raw DBGp commands are passed to the DBGp engine without further validation.
 * This is useful for troubleshooting Footle.
 *
 * There is no need to provide the transaction ID as part of the raw command as
 * we add it here.
 */
func prepareRawDBGpCmd(args []string, TxId int) (DBGpCmdWTxId string, err error) {

	if len(args) < 1 {
		return DBGpCmdWTxId, fmt.Errorf("No raw DBGp command given.")
	}

	rawDBGpCmd := strings.Join(args, space)
	DBGpCmdWTxId = fmt.Sprintf("%s -i %d\x00", rawDBGpCmd, TxId)

	return DBGpCmdWTxId, err
}

/**
 * DBGp Source command.
 */
func prepareSourceCmd(args []string, TxId int) (DBGpCmd string, err error) {

	if 2 != len(args) {
		err = fmt.Errorf("Insufficient number of args for source.")
		return DBGpCmd, err
	}

	beginLine, err := strconv.ParseInt(args[0], 10, 64)
	lineCount, err := strconv.ParseInt(args[1], 10, 64)
	endLine := beginLine + lineCount

	DBGpCmd = fmt.Sprintf("source -i %d -b %d -e %d\x00", TxId, beginLine, endLine)

	return DBGpCmd, err
}

/**
 * DBGp property_get command.
 *
 * It fetches the value of a single variable.
 *
 * Example: property_get -i 9 -n "foo"
 */
func preparePropertyGetCmd(args []string, TxId int) (DBGpCmd string, err error) {

	if len(args) < 1 {
		err = fmt.Errorf("Insufficient number of args for property_get.")
		return DBGpCmd, err
	}

	// Some variable names may contain a space character (e.g. foo["bar buz"]).
	// Such names will appear as separate argument items.  We reconstruct the
	// original variable name by joining the items.
	variableName := strings.Join(args, space)

	// Escapse following chars with backslash: single quote, double quote, null,
	// and backslash as per the DBGp protocol.
	escapseRule := strings.NewReplacer(`'`, `\'`, `"`, `\"`, "\x00", "\\\x00", `\`, `\\`)
	variableName = escapseRule.Replace(variableName)

	DBGpCmd = fmt.Sprintf("property_get -i %d -n \"%s\"\x00", TxId, variableName)

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

/**
 * DBGp context_get command.
 *
 * It fetches all local or global variables.
 *
 * First argument, when present, can be "global" or anything else to denote a
 * "local" context.  Second argument, when present, must be a number.
 *
 * Example: context_get -i 9 -c 1, context_get -i 9 -c 0 -d 0
 */
func prepareContextGetCmd(args []string, TxId int) (DBGpCmd string, err error) {

	if err = validateContextGetArgs(args); err != nil {
		return DBGpCmd, err
	}

	contextId := localContextId // Default is local context.
	argCount := len(args)
	stackDepth := 0

	if argCount == 0 {
		DBGpCmd = fmt.Sprintf("context_get -i %d -c %d\x00", TxId, localContextId)

		return DBGpCmd, err
	}

	if argCount >= 1 && args[0] == globalContextLabel {
		contextId = globalContextId
	}

	if argCount == 1 {
		DBGpCmd = fmt.Sprintf("context_get -i %d -c %d\x00", TxId, contextId)
	} else if argCount >= 2 {
		stackDepth, err = strconv.Atoi(args[1])

		DBGpCmd = fmt.Sprintf("context_get -i %d -c %d -d %d\x00", TxId, contextId, stackDepth)
	}

	return DBGpCmd, err
}
