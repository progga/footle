/**
 * Validate commands entered from the user interface.
 *
 * Commands coming from the user interface are shorter versions of the actual
 * DBGp commands.  This affords ease of use.  Here we verify these shorter
 * commands so that we can later prepare the full DBGp commands.
 */

package command

import (
	"fmt"
	"strconv"
)

/**
 * Wrapper for validating any command coming from the interface.
 */
func Validate(cmd string, args []string) (err error) {

	switch cmd {
	default:
		err = fmt.Errorf("Unknown command.")

	case "breakpoint_set", "b":
		err = validateBreakpointArgs(args)

	case "breakpoint_get", "bg":
		err = validateBreakpointGetArgs(args)

	case "breakpoint_list", "bl":
		err = validateCmdWithNoArg("breakpoint_list", args)

	case "context_get", "vl":
		err = validateCmdWithNoArg("context_get", args)

	case "eval", "ev":
		err = validateCmdWithNoArg("eval", args)

	case "run", "r":
		err = validateCmdWithNoArg("run", args)

	case "dbgp":
		err = validateRawDBGpArgs(args)

	case "property_get", "var":
		err = validatePropertyGetArgs(args)

	case "source", "src", "sr":
		err = validateSourceArgs(args)

	case "stop", "st":
		err = validateCmdWithNoArg("stop", args)

	case "status", "s":
		err = validateCmdWithNoArg("status", args)

	case "step_into", "si":
		err = validateCmdWithNoArg("step_into", args)

	case "step_out", "so":
		err = validateCmdWithNoArg("step_out", args)

	case "step_over", "sov", "sv":
		err = validateCmdWithNoArg("step_over", args)
	}

	return err
}

/**
 * Validate the Breakpoint command.
 */
func validateBreakpointArgs(args []string) (err error) {

	if len(args) != 2 {
		err = fmt.Errorf("Usage: breakpoint_set filepath line-number")
		return err
	}

	line_number, err := strconv.ParseInt(args[1], 10, 64)
	if nil != err || line_number < 1 {
		err = fmt.Errorf("Expecting line number as the second argument. %s given.", args[1])
	}

	return err
}

/**
 * Validate the Breakpoint get command.
 */
func validateBreakpointGetArgs(args []string) (err error) {

	if len(args) != 1 {
		err = fmt.Errorf("Usage: breakpoint_get breakpoint-id")
		return err
	}

	breakpointId, err := strconv.ParseInt(args[0], 10, 64)
	if nil != err || breakpointId < 1 {
		err = fmt.Errorf("Expecting breakpoint ID as the first argument. %s given.", args[0])
	}

	return err
}

/**
 * Validate any command that does not take any argument except TX ID.
 *
 * Example: run, stop, etc.
 */
func validateCmdWithNoArg(cmd string, args []string) (err error) {

	if 0 != len(args) {
		err = fmt.Errorf("The \"%s\" command does not take any argument.", cmd)
	}

	return err
}

/**
 * Validate the Source command.
 *
 * Valid format: source line-number line-count
 * Example: source 14 5
 *   This should return 5 lines starting at line number 14.
 */
func validateSourceArgs(args []string) (err error) {

	if 2 != len(args) {
		err = fmt.Errorf("The \"source\" command takes two numbers as argument.")
		return err
	}

	lineNumber, err := strconv.ParseInt(args[0], 10, 64)
	lineCount, err := strconv.ParseInt(args[1], 10, 64)

	if lineNumber < 1 {
		err = fmt.Errorf("Invalid line number.")
	}

	if lineCount < 1 {
		err = fmt.Errorf("Invalid line count.")
	}

	return err
}

/**
 * Validate the Raw DBGp command.
 *
 * Valid format: dbgp dbgp-command [dbgp-command-args]
 * Example: dbgp breakpoint_list
 *   This should execute the breakpoint_list command without validation.
 */
func validateRawDBGpArgs(args []string) (err error) {

	if len(args) < 1 {
		err = fmt.Errorf("The \"dbgp\" command expects at least one argument.")
		return err
	}

	return err
}

/**
 * Validate the property_get command.
 */
func validatePropertyGetArgs(args []string) (err error) {

	if len(args) != 1 {
		err = fmt.Errorf("The \"property_get\" command takes a variable name as an argument.")
		return err
	}

	return err
}
