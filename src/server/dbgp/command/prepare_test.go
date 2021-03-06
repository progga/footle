/**
 * Tests for DBGp command preparation.
 */

package command

import "testing"

/**
 * Tests for fetchNextTxId().
 */
func TestFetchNextTxId(t *testing.T) {

	TxId := fetchNextTxId()
	if 1 != TxId {
		t.Errorf("Expected TxId value of 1, got %q", TxId)
	}

	TxId = fetchNextTxId()
	if 2 != TxId {
		t.Errorf("Expected TxId value of 2, got %q", TxId)
	}

	TxId = fetchNextTxId()
	if 3 != TxId {
		t.Errorf("Expected TxId value of 3, got %q", TxId)
	}
}

/**
 * Tests for prepareBreakpointCmd().
 */
func TestPrepareBreakpointCmd(t *testing.T) {

	// Pass case.
	cmd, _ := prepareBreakpointCmd([]string{
		"/home/foo/code/php/bar.php",
		"9",
	}, 5)

	expected_cmd := "breakpoint_set -i 5 -t line -f /home/foo/code/php/bar.php -n 9\x00"
	if cmd != expected_cmd {
		t.Errorf("Incorrect breakpoint command. Expected %q, got %q.", expected_cmd, cmd)
	}

	// Fail case.
	cmd, err := prepareBreakpointCmd([]string{"foo"}, 3)
	if nil == err {
		t.Error("Missed insufficient number of args.")
	}
}

/**
 * Tests for prepareBreakpointGetCmd().
 */
func TestPrepareBreakpointGetCmd(t *testing.T) {

	// Pass case.
	cmd, _ := prepareBreakpointGetCmd([]string{"9"}, 5)

	expected_cmd := "breakpoint_get -i 5 -d 9\x00"
	if cmd != expected_cmd {
		t.Errorf("Incorrect breakpoint command. Expected %q, got %q.", expected_cmd, cmd)
	}

	// Fail case.
	cmd, err := prepareBreakpointGetCmd([]string{}, 3)
	if nil == err {
		t.Error("Missed insufficient number of args.")
	}
}

/**
 * Tests for prepareBreakpointRemoveCmd().
 */
func TestPrepareBreakpointRemoveCmd(t *testing.T) {

	// Pass case.
	cmd, _ := prepareBreakpointRemoveCmd([]string{"9"}, 5)

	expected_cmd := "breakpoint_remove -i 5 -d 9\x00"
	if cmd != expected_cmd {
		t.Errorf("Incorrect breakpoint remove command. Expected %q, got %q.", expected_cmd, cmd)
	}

	// Fail case.
	cmd, err := prepareBreakpointRemoveCmd([]string{}, 3)
	if nil == err {
		t.Error("Missed insufficient number of args.")
	}
}

/**
 * Tests for prepareEvalCmd().
 *
 * prepareEvalCmd() itself is incomplete at the moment.  So this test evolve.
 */
func TestPrepareEvalCmd(t *testing.T) {

	// Pass case.
	args := []string{"$a = 2 + 2"}
	TxId := 4
	cmd, err := prepareEvalCmd(args, TxId)

	expected := "eval -i 4 -- $a = 2 + 2\x00"
	if expected != cmd {
		t.Error("Eval command preparation failed.")
	}

	// Fail case.
	args = []string{}
	TxId = 4
	_, err = prepareEvalCmd(args, TxId)

	if nil == err {
		t.Error("Missed insufficient number of args.")
	}
}

/**
 * Tests for prepareSourceCmd().
 */
func TestPrepareSourceCmd(t *testing.T) {

	TxId := 4
	cmd, _ := prepareSourceCmd([]string{"10", "2"}, TxId)

	expected := "source -i 4 -b 10 -e 12\x00"
	if expected != cmd {
		t.Errorf("Source command preparation failed. Expected: %s, got: %s", expected, cmd)
	}

	cmd, _ = prepareSourceCmd([]string{"foo/bar/baz.php"}, TxId)

	expected = "source -i 4 -f foo/bar/baz.php\x00"
	if expected != cmd {
		t.Errorf("Source command preparation failed. Expected: %s, got: %s", expected, cmd)
	}
}

/**
 * Tests for prepareCmdNoArgs().
 */
func TestPrepareCmdNoArgs(t *testing.T) {

	// Pass case.
	cmd := "foo"
	TxId := 33
	DBGpCmd, _ := prepareCmdNoArgs(cmd, TxId)

	expected := "foo -i 33\x00"
	if expected != DBGpCmd {
		t.Error("prepareCmdNoArgs(foo, 33) != foo -i 33")
	}

	// Fail case.
	cmd = " "
	_, err := prepareCmdNoArgs(cmd, TxId)

	if nil == err {
		t.Error("Failed to spot empty command string.")
	}
}

/**
 * Tests for prepareRawDBGpCmd().
 */
func TestPrepareRawDBGpCmd(t *testing.T) {

	// Pass case.
	args := []string{"breakpoint_list"}
	TxId := 4
	cmd, _ := prepareRawDBGpCmd(args, TxId)

	expected := "breakpoint_list -i 4\x00"
	if expected != cmd {
		t.Errorf("Raw DBGp command preparation failed. Expected: %s, got: %s", expected, cmd)
	}

	// Pass case.
	args = []string{"property_get", "-n", "foo"}
	TxId = 5
	cmd, _ = prepareRawDBGpCmd(args, TxId)

	expected = "property_get -n foo -i 5\x00"
	if expected != cmd {
		t.Errorf("Raw DBGp command preparation failed. Expected: %s, got: %s", expected, cmd)
	}

	// Fail case.
	args = []string{}
	TxId = 6
	_, err := prepareRawDBGpCmd(args, TxId)

	if err == nil {
		t.Errorf("Raw DBGp command preparation failed to spot empty argument.")
	}
}

/**
 * Tests for preparePropertyGetCmd().
 *
 * The preparePropertyGetCmd() expects a variable name as well as the
 * transaction Id.  It may optionally be given a global or local context name.
 *
 * It is safe to always wrap the variable name in double quotes in case the
 * name contains a space or null (unlikely case) character.  We escape the
 * following characters with backslash in the variable name: single quote,
 * double quote, null, backslash.
 *
 * @see https://xdebug.org/docs-dbgp.php#escaping-rules
 */
func TestPreparePropertyGetCmd(t *testing.T) {

	// Pass case.
	args := []string{"foo"}
	TxId := 4
	cmd, _ := preparePropertyGetCmd(args, TxId)

	expected := "property_get -i 4 -c 0 -n \"foo\"\x00"
	if cmd != expected {
		t.Errorf("property_get command preparation failed.  Expected: %s, got: %s", expected, cmd)
	}

	// Pass case where variable name (e.g. foo['bar baz'\x00\]) has single quote,
	// space character, null, and backslash.  This tests variable name escaping.
	args = []string{"foo['bar", "baz'\x00\\]"}
	TxId = 44
	cmd, _ = preparePropertyGetCmd(args, TxId)

	expected = "property_get -i 44 -c 0 -n \"foo[\\'bar baz\\'\\\x00\\\\]\"\x00"
	if cmd != expected {
		t.Errorf("property_get command preparation failed.  Expected: %s, got: %s", expected, cmd)
	}

	// Pass case where global context name is used to fetch a global variable.
	args = []string{globalContextLabel, "foo"}
	TxId = 4
	cmd, _ = preparePropertyGetCmd(args, TxId)

	expected = "property_get -i 4 -c 1 -n \"foo\"\x00"
	if cmd != expected {
		t.Errorf("property_get command preparation failed.  Expected: %s, got: %s", expected, cmd)
	}

	// Pass case where local context name is used to fetch a local variable.
	args = []string{localContextLabel, "qux"}
	TxId = 4
	cmd, _ = preparePropertyGetCmd(args, TxId)

	expected = "property_get -i 4 -c 0 -n \"qux\"\x00"
	if cmd != expected {
		t.Errorf("property_get command preparation failed.  Expected: %s, got: %s", expected, cmd)
	}

	// Pass case where we attempt to fetch a local variable named "global".
	args = []string{localContextLabel, "global"}
	TxId = 4
	cmd, _ = preparePropertyGetCmd(args, TxId)

	expected = "property_get -i 4 -c 0 -n \"global\"\x00"
	if cmd != expected {
		t.Errorf("property_get command preparation failed.  Expected: %s, got: %s", expected, cmd)
	}

	// Fail case.
	args = []string{}
	TxId = 5
	_, err := preparePropertyGetCmd(args, TxId)

	if err == nil {
		t.Error("Failed to spot lack of a variable name.")
	}

	// Fail case where we attempt to fetch a local variable named "global" without
	// mentioning the context name.
	args = []string{"global"}
	TxId = 4
	cmd, _ = preparePropertyGetCmd(args, TxId)

	if err == nil {
		t.Error("Failed to spot lack of a context name.")
	}
}

/**
 * Tests for prepareContextGetCmd().
 *
 * Example DBGp commands: context_get -i N -c N, context_get -i N -c N -d N
 */
func TestPrepareContextGetCmd(t *testing.T) {

	// Pass case.
	args := []string{}
	TxId := 5
	cmd, _ := prepareContextGetCmd(args, TxId)

	expected := "context_get -i 5 -c 0\x00"
	if cmd != expected {
		t.Errorf("context_get command preparation failed.  Expected: %s, got: %s", expected, cmd)
	}

	args = []string{"global", "0"}
	TxId = 5
	cmd, _ = prepareContextGetCmd(args, TxId)

	expected = "context_get -i 5 -c 1 -d 0\x00"
	if cmd != expected {
		t.Errorf("context_get command preparation failed.  Expected: %s, got: %s", expected, cmd)
	}

	// Fail case.
	args = []string{"global", "Foo"}
	TxId = 5
	_, err := prepareContextGetCmd(args, TxId)

	if err == nil {
		t.Error("Failed to spot invalid stack depth.")
	}
}

/**
 * Tests for prepareFeatureSetCmd().
 *
 * Format: feature_set -i 9 -n FOO -v bar
 */
func TestPrepareFeatureSetCmd(t *testing.T) {

	// Pass case.
	args := []string{"foo", "bar"}
	TxId := 5
	cmd, err := prepareFeatureSetCmd(args, TxId)

	expected := "feature_set -i 5 -n foo -v bar\x00"
	if cmd != expected {
		t.Errorf("feature_set command preparation failed.  Expected: %s, got: %s", expected, cmd)
	}

	// Fail case.
	args = []string{}
	_, err = prepareFeatureSetCmd(args, TxId)
	if err == nil {
		t.Error("Failed to spot missing arguments.")
	}
}
