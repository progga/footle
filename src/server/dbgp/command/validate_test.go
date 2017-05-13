
/**
 * Tests for command validation.
 *
 * These commands come from the command line user interface.
 */

package command

import "testing"

/**
 * Tests for validateBreakpointArgs()
 */
func TestValidateBreakpointArgs(t *testing.T) {

  // Pass case.
  err := validateBreakpointArgs([]string {
    "/home/foo/bar.php",
    "28",
  })

  if nil != err {
    t.Error(err)
  }

  // Fail case.
  err = validateBreakpointArgs([]string{
    "/home/foo/bar.php",
  })

  if nil == err {
    t.Error("Failed to spot missing argument for the breakpoint_set command.")
  }
}

/**
 * Tests for validateCmdWithNoArg().
 */
func TestValidateCmdWithNoArg(t *testing.T) {

  err := validateCmdWithNoArg("run", []string {
  })

  if nil != err {
    t.Error(err)
  }

  err = validateCmdWithNoArg("run", []string {
    "foo",
    "bar",
  })

  if nil == err {
    t.Error("Failed to spot non-zero arguments for the \"run\" command.")
  }
}

/**
 * Tests for validateSourceArgs().
 */
func TestValidateSourceArgs(t *testing.T) {

  // Pass case.
  err := validateSourceArgs([]string { "10", "5" })

  if nil != err {
    t.Error(err)
  }

  // Fail case.
  err = validateSourceArgs([]string {})

  if nil == err {
    t.Error("Failed to spot lack of arguments.")
  }

  // Fail case.
  err = validateSourceArgs([]string {"1"})

  if nil == err {
    t.Error("Failed to spot insufficient number of arguments.")
  }
}
