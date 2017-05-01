
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
func Prepare(cmd string, args []string) (xdebug_cmd string, err error) {

  TxId := fetchNextTxId()

  switch strings.ToLower(cmd) {
    case "breakpoint", "b":
      xdebug_cmd, err = prepareBreakpointCmd(args, TxId)

    case "status", "s":
      xdebug_cmd, err = prepareCmdNoArgs("status", TxId)

    case "run", "r":
      xdebug_cmd, err = prepareCmdNoArgs("run", TxId)

    case "step_into", "si":
      xdebug_cmd, err = prepareCmdNoArgs("step_into", TxId)

    case "step_out", "so":
      xdebug_cmd, err = prepareCmdNoArgs("step_out", TxId)

    case "step_over", "sov", "sv":
      xdebug_cmd, err = prepareCmdNoArgs("step_over", TxId)

    case "eval", "ev":
      xdebug_cmd, err = prepareEvalCmd(args, TxId)

    default:
      xdebug_cmd, err = "", fmt.Errorf("Unknown command: %s", cmd)
  }

  return xdebug_cmd, err
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
 *
 * @param []string args
 * @param int TxId
 * @return string
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
 * Xdebug status command.
 *
 * @param int TxId
 * @return string xdebug_cmd
 * @return error err
 */
func prepareStatusCmd(TxId int) (xdebug_cmd string, err error) {

  xdebug_cmd = fmt.Sprintf("status -i %d\x00", TxId)

  return xdebug_cmd, err
}

/**
 * Xdebug run command.
 *
 * @param int TxId
 * @return string xdebug_cmd
 * @return error err
 */
func prepareRunCmd(TxId int) (xdebug_cmd string, err error) {

  xdebug_cmd = fmt.Sprintf("run -i %d\x00", TxId)

  return xdebug_cmd, err
}

/**
 * Xdebug step_into command.
 *
 * @param int TxId
 * @return string xdebug_cmd
 * @return error err
 */
func prepareStepIntoCmd(TxId int) (xdebug_cmd string, err error) {

  xdebug_cmd = fmt.Sprintf("step_into -i %d\x00", TxId)

  return xdebug_cmd, err
}

/**
 * Xdebug step_out command.
 *
 * @param int TxId
 * @return string xdebug_cmd
 * @return error err
 */
func prepareStepOutCmd(TxId int) (xdebug_cmd string, err error) {

  xdebug_cmd = fmt.Sprintf("step_out -i %d\x00", TxId)

  return xdebug_cmd, err
}

/**
 * Xdebug step_over command.
 *
 * @param int TxId
 * @return string xdebug_cmd
 * @return error err
 */
func prepareStepOverCmd(TxId int) (xdebug_cmd string, err error) {

  xdebug_cmd = fmt.Sprintf("step_over -i %d\x00", TxId)

  return xdebug_cmd, err
}

/**
 * Xdebug eval command.
 *
 * @param []string args
 * @param int TxId
 * @return string xdebug_cmd
 * @return error err
 */
func prepareEvalCmd(args []string, TxId int) (xdebug_cmd string, err error) {

  if 0 == len(args) {
    return xdebug_cmd, fmt.Errorf("Unsufficient number of args for eval.")
  }

  xdebug_cmd = fmt.Sprintf("eval -i %d -- %s\x00", TxId, args[0])

  return xdebug_cmd, err
}

/**
 *
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
