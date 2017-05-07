
/**
 * Functions for DBGp command preparation.
 */

package command

import (
  "fmt"
  "math"
  "strings"
)

/**
 * Xdebug transaction ID.  Its value was last used for the xdebug command's
 * transaction ID.
 *
 * @see fetchNextTxId()
 */
var lastTxId int

/**
 * Prepare xdebug command from the given values.
 *
 * The arguments for this function can be considered a short form of the xdebug
 * commands.  Here the short forms are expanded to their full forms.
 *
 * Example: "b /home/foo/php/bar.php 9" becomes
 * "breakpoint_set -i 1 -t line -f /home/foo/php/bar.php -n 9"
 *
 * @param string cmd
 * @param []string args
 * @return string xdebug_cmd
 * @return error err
 */
func Prepare(cmd string, args []string) (xdebugCmd string, err error) {

  TxId := fetchNextTxId()

  switch strings.ToLower(cmd) {
    case "breakpoint", "b":
      xdebugCmd, err = prepareBreakpointCmd(args, TxId)

    case "status", "s":
      xdebugCmd, err = prepareCmdNoArgs("status", TxId)

    case "run", "r":
      xdebugCmd, err = prepareCmdNoArgs("run", TxId)

    case "stop", "st":
      xdebugCmd, err = prepareCmdNoArgs("stop", TxId)

    case "step_into", "si":
      xdebugCmd, err = prepareCmdNoArgs("step_into", TxId)

    case "step_out", "so":
      xdebugCmd, err = prepareCmdNoArgs("step_out", TxId)

    case "step_over", "sov", "sv":
      xdebugCmd, err = prepareCmdNoArgs("step_over", TxId)

    case "eval", "ev":
      xdebugCmd, err = prepareEvalCmd(args, TxId)

    default:
      xdebugCmd, err = "", fmt.Errorf("Unknown command: %s", cmd)
  }

  return xdebugCmd, err
}

/**
 * Determine the transaction ID for the next xdebug command.
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
 * The full xdebug breakpoint command.
 */
func prepareBreakpointCmd(args []string, TxId int) (xdebug_cmd string, err error) {

  if 2 > len(args) {
    return xdebug_cmd, fmt.Errorf("Need at least two args for preparing breakpoint cmd.")
  }

  filepath    := args[0]
  line_number := args[1]

  xdebug_cmd = fmt.Sprintf("breakpoint_set -i %d -t line -f %s -n %s\x00", TxId, filepath, line_number)

  return xdebug_cmd, err
}

/**
 * Xdebug eval command.
 */
func prepareEvalCmd(args []string, TxId int) (xdebug_cmd string, err error) {

  if 0 == len(args) {
    return xdebug_cmd, fmt.Errorf("Unsufficient number of args for eval.")
  }

  xdebug_cmd = fmt.Sprintf("eval -i %d -- %s\x00", TxId, args[0])

  return xdebug_cmd, err
}

/**
 * Any Xdebug command that does not take any argument other than the TX ID.
 *
 * Example: run, stop, etc.
 */
func prepareCmdNoArgs(cmd string, TxId int) (xdebugCmd string, err error) {

  cmd = strings.TrimSpace(cmd)

  if "" == cmd {
    err = fmt.Errorf("Command cannot be empty.")

    return xdebugCmd, err
  }

  xdebugCmd = fmt.Sprintf("%s -i %d\x00", cmd, TxId)

  return xdebugCmd, err
}
