/**
 * @file
 * Tests for resolving DBGp command aliases.
 */

package command

import "testing"

/**
 * Tests for Extract().
 */
func TestExtract(t *testing.T) {

	passCases := []struct {
		CmdFromUI string
		CmdName   string
	}{
		// {Input, expected}
		{"b foo.php 20", "breakpoint_set"},
		{"r", "run"},
		{"breakpoint_set qux.py 99", "breakpoint_set"},
	}

	for _, test := range passCases {
		if DBGpCmdName, err := Extract(test.CmdFromUI); err != nil {
			t.Errorf("Extract(%s) = %s", test.CmdFromUI, DBGpCmdName)
		}
	}

	failCase := "foo bar baz"
	if DBGpCmdName, err := Extract(failCase); err == nil {
		t.Errorf("Extract(%s) = %s", failCase, DBGpCmdName)
	}
}
