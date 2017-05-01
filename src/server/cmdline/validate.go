
package cmdline

import (
  "fmt"
  "strconv"
)

/**
 *
 */
func Validate(cmd string, args []string) (err error) {

  switch cmd {
    default:
      err = fmt.Errorf("Unknown command.")

    case "breakpoint_set", "b":
      err = validateBreakpointArgs(args)

    case "run", "r":
      err = validateCmdWithNoArg("run", args)

    case "status", "s":
      err = validateCmdWithNoArg("status", args)

    case "step_into", "si":
      err = validateCmdWithNoArg("step_into", args)

    case "step_out", "so":
      err = validateCmdWithNoArg("step_out", args)

    case "step_over", "sov", "sv":
      err = validateCmdWithNoArg("step_over", args)

    case "eval", "ev":
      err = validateCmdWithNoArg("eval", args)
  }

  return err
}

/**
 *
 */
func validateBreakpointArgs(args []string) (err error) {

  if (len(args) != 2) {
    err = fmt.Errorf("Usage: breakpoint_set filepath line-number")
    return err
  }

  line_number, err := strconv.ParseInt(args[1], 10, 64);
  if nil != err || ! (0 < line_number) {
    err = fmt.Errorf("Expecting line number as the second argument. %s given.", args[1])
  }

  return err
}

/**
 *
 */
func validateCmdWithNoArg(cmd string, args []string) (err error) {

  if 0 != len(args) {
    err = fmt.Errorf("The \"%s\" command does not take any argument.", cmd)
  }

  return err
}
