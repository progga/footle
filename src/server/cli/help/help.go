/**
 * @file
 * Helptext for the cli.
 */

package help

var cliCmdList []helptext = []helptext{
	helptext{[]string{"bye", "quit", "q"}, "Quits Footle."},
	helptext{[]string{"refresh"}, "Updates cli with any pending DBGp messages."},
	helptext{[]string{"verbose"}, "Dumps all traffic between Footle and the debugger engine."},
	helptext{[]string{"no-verbose"}, "Opposite of *verbose*."},
}

var footleCmdList []helptext = []helptext{
	helptext{[]string{"on"}, "Awake Footle.  It will then start responding to the debugger engine."},
	helptext{[]string{"off"}, "Put Footle to sleep.  It won't then respond to the debugger engine."},
	helptext{[]string{"continue"}, "End execution.  Ignore all breakpoints if needed."},
	helptext{[]string{"update_source"}, "Refresh source code of a displayed file.\nExample: update_source foo.php"},
}

var DBGpCmdList []helptext = []helptext{
	helptext{[]string{"breakpoint_set", "b"}, "Usage: breakpoint_set FILEPATH LINE-NUMBER"},
	helptext{[]string{"breakpoint_get", "bg"}, "Usage: breakpoint_get BREAKPOINT-ID"},
	helptext{[]string{"breakpoint_remove", "br"}, "Usage: breakpoint_remove BREAKPOINT-ID"},
	helptext{[]string{"breakpoint_list", "bl"}, "Fetches all breakpoints, including the pending ones."},
	helptext{[]string{"context_get", "vl"}, "Fetches all variables.\nUsage: context_get [local|global [stack-depth-number]]\nExample: context_get; context_get local; context_get global 3"},
	helptext{[]string{"dbgp"}, "Useful for executing raw DBGp commands.  Do *not* provide the transaction ID.\nUsage: dbgp DBGP-COMMAND [DBGP-COMMAND-ARGS]\nExample: dbgp breakpoint_list"},
	helptext{[]string{"eval", "ev"}, "Broken, don't use."},
	helptext{[]string{"property_get", "var"}, "Fetch the value of a variable.  Usage: property_get [local|global] VARIABLE-NAME\nExample: property_get $foo; property_get global $bar.  When neither *local* nor *global* context is mentioned, local is assumed."},
	helptext{[]string{"run", "r"}, "Carry on with execution."},
	helptext{[]string{"stk", "stack_get"}, "Fetch current stack trace."},
	helptext{[]string{"source", "sr", "src"}, "Fetch source code.\nUsage: source line-number line-count; source filepath.  The first format extracts from the current file under execution."},
	helptext{[]string{"status", "s"}, "Display debugger engine status."},
	helptext{[]string{"step_into", "si"}, "Move into a function or method."},
	helptext{[]string{"step_out", "so"}, "Move out of the current function or method."},
	helptext{[]string{"step_over", "sv", "sov"}, "Move to next line please."},
	helptext{[]string{"stop", "st"}, "End execution."},
}

var help Help

/**
 * Prepare the Help object.
 */
func init() {

	help = Help{}
	help.prepare(cliCmdList, footleCmdList, DBGpCmdList)
}

/**
 * Getter for the Help object maintained by this package.
 */
func Get() (h Help) {

	return help
}
