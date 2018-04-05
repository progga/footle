/**
 * @file
 * Functions for dealing with command aliases.
 *
 * Example aliases: b for breakpoint_set, r for run, etc.
 */

package command

import "strings"

/**
 * Mapping between DBGp commands and their aliases.
 */
var shortCmdFullCmdMap map[string]string = map[string]string{
	"b":   "breakpoint_set",
	"bg":  "breakpoint_get",
	"br":  "breakpoint_remove",
	"bl":  "breakpoint_list",
	"vl":  "context_get",
	"ev":  "eval",
	"var": "property_get",
	"r":   "run",
	"stk": "stack_get",
	"sr":  "source",
	"src": "source",
	"s":   "status",
	"si":  "step_into",
	"so":  "step_out",
	"sv":  "step_over",
	"sov": "step_over",
	"st":  "stop",
}

/**
 * Determine the real name for an alias.
 *
 * Non-aliases are assumed to be DBGp command names.
 */
func resolveAlias(potentialAlias string) (DBGpCmd string) {

	DBGpCmd, ok := shortCmdFullCmdMap[potentialAlias]
	if !ok {
		DBGpCmd = strings.ToLower(potentialAlias)
	}

	return DBGpCmd
}
