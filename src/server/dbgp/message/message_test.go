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

	if nil != err && "init" != message.Message_type {
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

	if nil != err && "response" != message.Message_type {
		t.Error("Missed Response message.")
	}

	if nil != err && 14 != message.Properties.Line_number {
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
		t.Error(`parseInit(<init ... fileuri="file:///srv/www/drupal/drupal8/index.php""...>...</init>) cannot find file URI.`)
	}
}

/**
 * Tests for decodeResponse()
 */
func TestDecodeResponse(t *testing.T) {

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
		t.Error(`parseResponse(<response ... command="status"...></response>): Command is not "status"`)
	}
}
