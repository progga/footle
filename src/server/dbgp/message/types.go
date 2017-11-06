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
	DisplayName       string // Short name is useful for display purposes.
	VarType           string
	Value             string // Only for basic types such as int, float, string, etc.
	AccessModifier    string // private, protected, public, etc.
	IsCompositeType   bool   // Is it an array, object, structure, etc.?
	Children          map[string]Variable
	ChildCount        int
	HasLoadedChildren bool // DBGp servers return children upto a certain depth.
	IsBase64          bool
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
	Breakpoints   []Breakpoint      `xml:"breakpoint"`
	Error         Error             `xml:"error"`
	Variables     []VariableDetails `xml:"property"`
	Content       string            `xml:",chardata"`
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

type VariableDetails struct {
	Name        string            `xml:"name,attr"`
	Fullname    string            `xml:"fullname,attr"`
	VarType     string            `xml:"type,attr"`
	Facet       string            `xml:"facet,attr"` // public, private, etc.
	Classname   string            `xml:"classname,attr"`
	Constant    int               `xml:"constant,attr"`
	HasChildren bool              `xml:"children,attr"`
	Size        int               `xml:"size,attr"`
	Page        int               `xml:"page,attr"`
	Pagesize    int               `xml:"pagesize,attr"`
	Address     int               `xml:"address,attr"`
	Key         string            `xml:"key,attr"`
	Encoding    string            `xml:"encoding,attr"`
	NumChildren int               `xml:"numchildren,attr"`
	Value       string            `xml:",chardata"`
	Variables   []VariableDetails `xml:"property"`
}
