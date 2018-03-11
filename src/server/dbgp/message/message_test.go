/**
 * Tests for DBGp message decoding.
 */

package message

import "testing"

/**
 * Tests for Decode()
 */
func TestDecode(t *testing.T) {

	xml :=
		`<?xml version="1.0" encoding="iso-8859-1"?>
<init xmlns="urn:debugger_protocol_v1" xmlns:xdebug="http://xdebug.org/dbgp/xdebug"
  fileuri="file:///srv/www/drupal/drupal8/index.php"
  language="PHP"
  protocol_version="1.0"
  appid="27891">
  <engine version="2.2.5"><![CDATA[Xdebug]]></engine>
  <author><![CDATA[Derick Rethans]]></author>
  <url><![CDATA[http://xdebug.org]]></url>
  <copyright><![CDATA[Copyright (c) 2002-2014 by Derick Rethans]]></copyright>
</init>`

	message, err := Decode(xml)
	if nil != err {
		t.Error(err)
	}

	if nil != err && "init" != message.MessageType {
		t.Error("Missed Init message.")
	}

	xml =
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

	message, err = Decode(xml)
	if nil != err {
		t.Error(err)
	}

	if nil != err && "response" != message.MessageType {
		t.Error("Missed Response message.")
	}

	if nil != err && 14 != message.Properties.LineNumber {
		t.Error("Missed line number.")
	}

	xml =
		`<foo
  command="run"
  transaction_id="3"
  status="break"
  reason="ok">
</foo>`

	message, err = Decode(xml)
	if nil == err {
		t.Error("Missed unknown message.")
	} else {
		t.Log(err)
	}
}

/**
 * Tests for decodeInit()
 */
func TestDecodeInit(t *testing.T) {

	xml :=
		`<?xml version="1.0" encoding="iso-8859-1"?>
<init xmlns="urn:debugger_protocol_v1" xmlns:xdebug="http://xdebug.org/dbgp/xdebug"
  fileuri="file:///srv/www/drupal/drupal8/index.php"
  language="PHP"
  protocol_version="1.0"
  appid="27891">
  <engine version="2.2.5"><![CDATA[Xdebug]]></engine>
  <author><![CDATA[Derick Rethans]]></author>
  <url><![CDATA[http://xdebug.org]]></url>
  <copyright><![CDATA[Copyright (c) 2002-2014 by Derick Rethans]]></copyright>
</init>`

	init, err := decodeInit(xml)
	if nil != err {
		t.Error(err)
	}

	if "file:///srv/www/drupal/drupal8/index.php" != init.FileURI {
		t.Error(`decodeInit(<init ... fileuri="file:///srv/www/drupal/drupal8/index.php""...>...</init>) cannot find file URI.`)
	}
}

/**
 * Tests for decodeResponse()
 */
func TestDecodeResponse(t *testing.T) {

	// Fail case.
	xml :=
		`<?xml version="1.0" encoding="iso-8859-1"?>
<init xmlns="urn:debugger_protocol_v1" xmlns:xdebug="http://xdebug.org/dbgp/xdebug"
  fileuri="file:///srv/www/drupal/drupal8/index.php"
  language="PHP"
  protocol_version="1.0"
  appid="27891">
  <engine version="2.2.5"><![CDATA[Xdebug]]></engine>
  <author><![CDATA[Derick Rethans]]></author>
  <url><![CDATA[http://xdebug.org]]></url>
  <copyright><![CDATA[Copyright (c) 2002-2014 by Derick Rethans]]></copyright>
</init>`

	if _, err := decodeResponse(xml); nil == err {
		t.Error(err)
	}

	// DBGP "status" command.
	xml =
		`<?xml version="1.0" encoding="iso-8859-1"?>
<response xmlns="urn:debugger_protocol_v1" xmlns:xdebug="http://xdebug.org/dbgp/xdebug"
  command="status"
  transaction_id="0"
  status="starting"
  reason="ok">
</response>`

	response, err := decodeResponse(xml)
	if nil != err {
		t.Error(err)
	}

	if "status" != response.Command {
		t.Error(`decodeResponse(<response ... command="status"...></response>): Command is not "status"`)
	}

	// DBGP "breakpoint_get" command.
	xml =
		`<?xml version="1.0" encoding="iso-8859-1"?>
<response xmlns="urn:debugger_protocol_v1" xmlns:xdebug="http://xdebug.org/dbgp/xdebug" command="breakpoint_get" transaction_id="2">
  <breakpoint
    type="line"
    filename="file:///srv/www/drupal/drupal8/index.php"
    lineno="14"
    state="enabled"
    hit_count="0"
    hit_value="0"
    id="68310001">
  </breakpoint>
</response>`

	response, err = decodeResponse(xml)
	if nil != err {
		t.Error(err)
	}

	if "breakpoint_get" != response.Command {
		t.Error(`decodeResponse(<response ... command="breakpoint_get"...></response>): Command is not "breakpoint_get"`)
	}

	if 14 != response.Breakpoints[0].LineNo {
		t.Errorf("Failed to spot Line number. %d given.", response.Breakpoints[0].LineNo)
	}

	if 68310001 != response.Breakpoints[0].Id {
		t.Errorf("Failed to spot Breakpoint ID. %d given.", response.Breakpoints[0].Id)
	}
}

/**
 * Tests for decodeResponse() against stack_get.
 *
 * Verify the decoded response for the DBGp stack_get command.
 */
func TestDecodeResponseForStackGet(t *testing.T) {

	xml :=
		`<response xmlns="urn:debugger_protocol_v1" xmlns:xdebug="http://xdebug.org/dbgp/xdebug" command="stack_get" transaction_id="10">
  <stack
    where="Drupal\Core\DrupalKernel::bootEnvironment"
    level="0"
    type="file"
    filename="file:///srv/www/drupal/drupal8/core/lib/Drupal/Core/DrupalKernel.php"
    lineno="882">
  </stack>
  <stack
    where="Drupal\Core\DrupalKernel-&gt;handle"
    level="1"
    type="file"
    filename="file:///srv/www/drupal/drupal8/core/lib/Drupal/Core/DrupalKernel.php"
    lineno="615">
  </stack>
  <stack
    where="{main}"
    level="2"
    type="file"
    filename="file:///srv/www/drupal/drupal8/index.php"
    lineno="19">
  </stack>
</response>`

	response, err := decodeResponse(xml)
	if nil != err {
		t.Error(err)
	}

	if "stack_get" != response.Command {
		t.Error(`decodeResponse(<response ... command="stack_get"...></response>): Command is not "stack_get"`)
	}

	if response.StackDetail[0].LineNo != 882 {
		t.Errorf("Failed to spot line number in call stack record. %d given.", response.StackDetail[0].LineNo)
	}

	if response.StackDetail[0].Level != 0 {
		t.Errorf("Failed to spot stack level number. %d given.", response.StackDetail[0].Level)
	}

	if response.StackDetail[2].Filename != "file:///srv/www/drupal/drupal8/index.php" {
		t.Errorf("Failed to spot filename in call stack record. %s given.", response.StackDetail[2].Filename)
	}
}

/**
 * Tests for extractVariableValue().
 */
func TestExtractVariableValue(t *testing.T) {

	// Text value.  Should be decoded.
	varDetails := VariableDetails{
		Encoding: "base64",
		Value:    "Zm9v", // "foo"
	}

	value, isBase64 := extractVariableValue(varDetails)

	if value != "foo" {
		t.Errorf("Failed to extract value of the variable.  Expecting \"foo\", got \"%s\"", value)
	}

	if isBase64 {
		t.Error("Should have decoded Base64 content.")
	}

	// Binary value.  Should not be decoded.
	varDetails = VariableDetails{
		Encoding: "base64",
		Value:    "Zm9vAGJhcg==", // "foo\x00bar"
	}

	value, isBase64 = extractVariableValue(varDetails)

	if !isBase64 {
		t.Error("Failed to spot binary content.")
	}

	// Base64 encoding is not always used for plain content.
	varDetails = VariableDetails{
		Encoding: "none",
		Value:    "99.99",
	}

	value, isBase64 = extractVariableValue(varDetails)

	if value != "99.99" {
		t.Errorf("Variable value was not encoded.  Expected \"99.99\", got \"%s\"", value)
	}

	if isBase64 {
		t.Error("Failed to spot plain encoding.")
	}
}
