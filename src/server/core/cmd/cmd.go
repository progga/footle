/**
 * @file
 * Simplifies dealing with internal Footle commands.
 *
 * These commands are specific to Footle.  They are unrelated to DBGp commands.
 * They drive Footle's internal state.  Examples include telling Footle to
 * disengage from the debugger engine (off), telling all the UIs to update a
 * certain file (update_source), etc.
 */

package cmd

import (
	"fmt"
	"strings"
)

/**
 * Is this an internal Footle command?
 */
func Is(cmdName string) (result bool) {

	result = cmdName == "on" ||
		cmdName == "off" ||
		cmdName == "continue" ||
		cmdName == "update_source"

	return result
}

/**
 * Validate the given command and its argument.
 */
func Validate(cmdName string, args []string) (valid bool, err error) {

	argCount := len(args)

	if (cmdName == "on" || cmdName == "off" || cmdName == "continue") && argCount == 0 {
		valid = true
	} else if cmdName == "update_source" && argCount == 1 {
		valid = true
	}

	if valid {
		return valid, err
	}

	cmd := cmdName + strings.Join(args, " ")

	if cmdName == "on" || cmdName == "off" || cmdName == "continue" {
		err = fmt.Errorf("Invalid command: %s; The right format is: %s", cmd, cmdName)
	} else if cmdName == "update_source" {
		err = fmt.Errorf("Invalid command: %s; The right format is: update_source FILENAME", cmd)
	} else {
		err = fmt.Errorf("Invalid command: %s", cmd)
	}

	return valid, err
}
