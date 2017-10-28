/**
 * Tests for command validation.
 *
 * These commands come from the command line user interface.
 */

package command

import "testing"

/**
 * Tests for Validate().
 */
func TestValidate(t *testing.T) {

  // The "dbgp" command.
  err := Validate("dbgp", []string{"foo", "bar"})

  if err != nil {
    t.Error(err)
  }

  // The "run" command.
  err = Validate("run", []string{})

  if err != nil {
    t.Error(err)
  }
}

/**
 * Tests for validateBreakpointArgs()
 */
func TestValidateBreakpointArgs(t *testing.T) {

	// Pass case.
	err := validateBreakpointArgs([]string{
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
 * Tests for validateBreakpointGetArgs()
 */
func TestValidateBreakpointGetArgs(t *testing.T) {

	// Pass case.
	err := validateBreakpointGetArgs([]string{"28"})
	if nil != err {
		t.Error(err)
	}

	// Fail case.
	err = validateBreakpointGetArgs([]string{})
	if nil == err {
		t.Error("Failed to spot missing argument for the breakpoint_get command.")
	}

	// Another fail case.  Invalid breakpoint ID.
	err = validateBreakpointGetArgs([]string{"0"})
	if nil == err {
		t.Error("Failed to spot invalid breakpoint ID for the breakpoint_get command.")
	}
}

/**
 * Tests for validateCmdWithNoArg().
 */
func TestValidateCmdWithNoArg(t *testing.T) {

	err := validateCmdWithNoArg("run", []string{})

	if nil != err {
		t.Error(err)
	}

	err = validateCmdWithNoArg("run", []string{
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
	err := validateSourceArgs([]string{"10", "5"})

	if nil != err {
		t.Error(err)
	}

	// Fail case.
	err = validateSourceArgs([]string{})

	if nil == err {
		t.Error("Failed to spot lack of arguments.")
	}

	// Fail case.
	err = validateSourceArgs([]string{"1"})

	if nil == err {
		t.Error("Failed to spot insufficient number of arguments.")
	}
}

/**
 * Tests for validateRawDBGpArgs().
 */
func TestValidateRawDBGpArgs(t *testing.T) {

  // Pass case.
  err := validateRawDBGpArgs([]string{"foo"})

  if nil != err {
    t.Error(err)
  }

  // Fail case.  The "DBGp" command expects at least one argument.
  err = validateRawDBGpArgs([]string{})

  if nil == err {
    t.Error(err)
  }
}
