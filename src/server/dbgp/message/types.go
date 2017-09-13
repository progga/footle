/**
 * Structure declarations that correspond to DBGp XML messages.
 *
 * These structures are used for decoding the XML messages into usable values.
 */

package message

import "encoding/xml"

type Message struct {
	MessageType string
	State       string
	Properties  Properties
	Context     Context
	Content     string
	Breakpoints map[int]Breakpoint
}

type Properties struct {
	Command      string
	ErrorCode    int
	ErrorMessage string
	Filename     string
	LineNumber   int
	BreakpointId int
	TxId         int
}

type Context struct {
	Local    map[string]Variable
	Global   map[string]Variable
	Constant map[string]string
}

type Variable struct {
	VarType string
	Literal string
	List    []Variable
}

type Init struct {
	XMLName   xml.Name `xml:"init"`
	FileURI   string   `xml:"fileuri,attr"`
	Language  string   `xml:"language,attr"`
	Protocol  string   `xml:"protocol_version,attr"`
	AppID     string   `xml:"appid,attr"`
	Engine    string   `xml:"engine"`
	Author    string   `xml:"author"`
	URL       string   `xml:"url"`
	Copyright string   `xml:"copyright"`
}

type Response struct {
	XMLName       xml.Name `xml:"response"`
	Command       string   `xml:"command,attr"`
	TransactionId int      `xml:"transaction_id,attr"`
	Status        string   `xml:"status,attr"`
	Reason        string   `xml:"reason,attr"`
	Id            int      `xml:"id,attr"`
	Message       ResponseMessage
	Breakpoints   []Breakpoint `xml:"breakpoint"`
	Error         Error        `xml:"error"`
	Content       string       `xml:",chardata"`
}

type ResponseMessage struct {
	XMLName  xml.Name `xml:"http://xdebug.org/dbgp/xdebug message"`
	Filename string   `xml:"filename,attr"`
	LineNo   int      `xml:"lineno,attr"`
}

type Breakpoint struct {
	Filename string `xml:"filename,attr"`
	LineNo   int    `xml:"lineno,attr"`
	Type     string `xml:"type,attr"`
	State    string `xml:"state,attr"`
	HitCount int    `xml:"hit_count,attr"`
	HitValue int    `xml:"hit_value,attr"`
	Id       int    `xml:"id,attr"`
}

type Error struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:"message"`
}
