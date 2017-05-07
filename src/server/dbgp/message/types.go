
/**
 * Structure declarations that correspond to DBGp XML messages.
 *
 * These structures are used for decoding the XML messages into usable values.
 */

package message

import "encoding/xml"

type Message struct {
  Message_type string
  State string
  Properties Properties
  Context Context
}

type Properties struct {
  Command string
  Error_code int
  Error_message string
  Filename string
  Line_number int
  TxId int
}

type Context struct {
  Local map[string]Variable
  Global map[string]Variable
  Constant map[string]string
}

type Variable struct {
  Var_type string
  Literal string
  List []Variable
}

type Init struct {
  XMLName xml.Name `xml:"init"`
  FileURI string `xml:"fileuri,attr"`
  Language string `xml:"language,attr"`
  Protocol string `xml:"protocol_version,attr"`
  AppID string `xml:"appid,attr"`
  Engine string `xml:"engine"`
  Author string `xml:"author"`
  URL string `xml:"url"`
  Copyright string `xml:"copyright"`
}

type Response struct {
  XMLName xml.Name `xml:"response"`
  Command string `xml:"command,attr"`
  Transaction_id int `xml:"transaction_id,attr"`
  Status string `xml:"status,attr"`
  Reason string `xml:"reason,attr"`
  Id int `xml:"id"`
  Message ResponseMessage
  Error Error `xml:"error"`
}

type ResponseMessage struct {
  XMLName xml.Name `xml:"http://xdebug.org/dbgp/xdebug message"`
  Filename string `xml:"filename,attr"`
  LineNo int `xml:"lineno,attr"`
}

type Error struct {
  Code int `xml:"code,attr"`
  Message string `xml:"message"`
}
