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

	DBGpCmd := resolveAlias(cmd)

	switch DBGpCmd {
	default:
		err = fmt.Errorf("Unknown command.")

	case "breakpoint_set":
		err = validateBreakpointArgs(args)

	case "breakpoint_get":
		err = validateBreakpointGetArgs(args)

	case "breakpoint_remove":
		err = validateBreakpointRemoveArgs(args)

	case "breakpoint_list":
		err = validateCmdWithNoArg("breakpoint_list", args)

	case "context_get":
		err = validateContextGetArgs(args)

	case "eval":
		err = validateCmdWithNoArg("eval", args)

	case "run":
		err = validateCmdWithNoArg("run", args)

	case "dbgp":
		err = validateRawDBGpArgs(args)

	case "property_get":
		err = validatePropertyGetArgs(args)

	case "feature_set":
		err = validateFeatureSetArgs(args)

	case "stack_get":
		err = validateCmdWithNoArg("stack_get", args)

	case "source":
		err = validateSourceArgs(args)

	case "stop":
		err = validateCmdWithNoArg("stop", args)

	case "status":
		err = validateCmdWithNoArg("status", args)

	case "step_into":
		err = validateCmdWithNoArg("step_into", args)

	case "step_out":
		err = validateCmdWithNoArg("step_out", args)

	case "step_over":
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
	_ = breakpointId
	if nil != err {
		err = fmt.Errorf("Expecting breakpoint ID as the first argument. %s given.", args[0])
	}

	return err
}

/**
 * Validate the Breakpoint remove command.
 */
func validateBreakpointRemoveArgs(args []string) (err error) {

	if len(args) != 1 {
		err = fmt.Errorf("Usage: breakpoint_remove breakpoint-id")
		return err
	}

	breakpointId, err := strconv.ParseInt(args[0], 10, 64)
	_ = breakpointId
	if nil != err {
		err = fmt.Errorf("Expecting a breakpoint ID as the only argument. %s given.", args[0])
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
 * Valid formats:
 *   - source line-number line-count
 *   - source filepath
 * Examples:
 *   - source 14 5: This should return 5 lines *starting* at line number 14.
 *   - source foo/bar/baz.php: This should return the complete source code of
 *     baz.php.
 */
func validateSourceArgs(args []string) (err error) {

	argCount := len(args)
	if argCount != 1 && argCount != 2 {
		err = fmt.Errorf("The \"source\" command takes a filepath OR two numbers as argument.")
		return err
	}

	if argCount == 2 {
		lineNumber, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		lineCount, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}

		if lineNumber < 1 {
			err = fmt.Errorf("Invalid line number.")
		}

		if lineCount < 1 {
			err = fmt.Errorf("Invalid line count.")
		}
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

	if len(args) < 1 {
		err = fmt.Errorf("The \"property_get\" command takes a variable name as an argument.")
		return err
	}

	return err
}

/**
 * Validate the arguments for the context_get command.
 *
 * Acceptable command formats: context_get, context_get local/global,
 * context local/global N
 */
func validateContextGetArgs(args []string) (err error) {

	argCount := len(args)

	if argCount == 0 {
		return err
	}

	if argCount == 1 && (args[0] != localContextLabel && args[0] != globalContextLabel) {
		err = fmt.Errorf("Invalid context.  Acceptable values: %s, %s.  %s given.", localContextLabel, globalContextLabel, args[0])

		return err
	}

	if argCount == 2 {
		_, err = strconv.Atoi(args[1])
	} else if argCount > 2 {
		err = fmt.Errorf("Too many arguments.")
	}

	return err
}

/**
 * feature_set command.
 *
 * Acceptable command format: feature_set feature-name feature-value
 */
func validateFeatureSetArgs(args []string) (err error) {

	argCount := len(args)
	if argCount != 2 {
		err = fmt.Errorf("The feature_set command takes two arguments.")

		return err
	}

	return err
}
