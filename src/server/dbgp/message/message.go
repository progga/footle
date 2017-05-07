
/**
 * Turn DBGp messages into usable data structures.
 */

package message

import (
  "fmt"
  "strings"
)

/**
 * Parse dbgp XML message.
 */
func Decode(xmlContent string) (message Message, err error) {

  has_response := (-1 != strings.LastIndex(xmlContent, "</response>"))
  if has_response {
    response, err := decodeResponse(xmlContent)

    if nil == err {
      message = prepareResponseMessage(response)
    }
  }

  has_init := (-1 != strings.LastIndex(xmlContent, "</init>"))
  if has_init {
    init, err := decodeInit(xmlContent)

    if nil == err {
      message = prepareInitMessage(init)
    }
  }

  if ! has_response && ! has_init {
    err = fmt.Errorf("Unknown message: %s", xmlContent)
  }

  return message, err
}

/**
 * Prepare a message structure based on DBGp engine's initialization attempt.
 */
func prepareInitMessage(init Init) (message Message) {

  message.Message_type = "init"
  message.State = "starting"
  message.Properties.Filename = init.FileURI

  return message
}

/**
 * Prepare a message structure based on DBGp engine's response.
 */
func prepareResponseMessage(response Response) (message Message) {

  message.Message_type             = "response"
  message.State                    = response.Status
  message.Properties.Filename      = response.Message.Filename
  message.Properties.Line_number   = response.Message.LineNo
  message.Properties.Error_message = response.Error.Message
  message.Properties.Error_code    = response.Error.Code
  message.Properties.TxId          = response.Transaction_id
  message.Properties.Command       = response.Command

  return message
}
