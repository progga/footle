/**
 * Tests for DBGp message parsing.
 */

package dbgp

import (
	"fmt"
	"strings"
	"testing"
)

/**
 * Tests for Read().
 */
func TestRead(t *testing.T) {

	// Pass case.
	xml :=
		`<?xml version="1.0" encoding="iso-8859-1"?>
<response xmlns="urn:debugger_protocol_v1" xmlns:xdebug="http://xdebug.org/dbgp/xdebug"
  command="run"
  transaction_id="3"
  status="break"
  reason="ok">
  <xdebug:message
    filename="file:///srv/www/drupal/drupal8/index.php"
    lineno="14">
  </xdebug:message>
</response>`

	DBGpMsg := fmt.Sprintf("%d\x00%s\x00", len(xml), xml)
	stringReader := strings.NewReader(DBGpMsg)
	msg, err := Read(stringReader)

	if xml != msg {
		t.Error("dbgp.read() failed to parse response type message.")
	}

	if nil != err {
		t.Error(err)
	}

	// Fail case.
	DBGpMsg = fmt.Sprintf("%d\x00%s\x00", len(xml)-10, xml)
	stringReader = strings.NewReader(DBGpMsg)
	msg, err = Read(stringReader)

	if nil == err {
		t.Error("Failed to spot broken DBGp message.")
	}
}
