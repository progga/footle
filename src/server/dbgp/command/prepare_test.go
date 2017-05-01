
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
  cmd, _ := prepareBreakpointCmd([]string {
    "/home/foo/code/php/bar.php",
    "9",
  }, 5)

  expected_cmd := "breakpoint_set -i 5 -t line -f /home/foo/code/php/bar.php -n 9\x00"
  if cmd != expected_cmd {
    t.Errorf("Incorrect breakpoint command. Expected %q, got %q.", expected_cmd, cmd)
  }

  // Fail case.
  cmd, err := prepareBreakpointCmd([]string {"foo"}, 3)
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
  args := []string {"$a = 2 + 2"}
  TxId := 4
  cmd, err := prepareEvalCmd(args, TxId)

  expected := "eval -i 4 -- $a = 2 + 2\x00"
  if expected != cmd {
    t.Errorf("Eval command preparation failed.")
  }

  // Fail case.
  args = []string {}
  TxId = 4
  _, err = prepareEvalCmd(args, TxId)

  if nil == err {
    t.Error("Missed insufficient number of args.")
  }
}

/**
 * Tests for prepareCmdNoArgs().
 */
func TestPrepareCmdNoArgs(t *testing.T) {

  // Pass case.
  cmd  := "foo"
  TxId := 33
  xdebugCmd, _ := prepareCmdNoArgs(cmd, TxId)

  expected := "foo -i 33\x00"
  if expected != xdebugCmd {
    t.Error("prepareCmdNoArgs(foo, 33) != foo -i 33")
  }

  // Fail case.
  cmd = " "
  _, err := prepareCmdNoArgs(cmd, TxId)

  if nil == err {
    t.Error("Failed to spot empty command string.")
  }
}
