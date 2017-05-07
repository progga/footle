
/**
 * Decode a DBGp message which is in XML format.
 *
 * Once decoded, we are able to access individual components of each message.
 */

package message

import (
  "encoding/xml"
  "golang.org/x/net/html/charset"
  "log"
  "strings"
)

/**
 * Decodes XML response from DBGp engine.
 */
func decodeResponse(xmlResponse string) (Response, error) {

  strReader := strings.NewReader(xmlResponse)

  decoder := xml.NewDecoder(strReader)
  // @see http://stackoverflow.com/questions/6002619/unmarshal-an-iso-8859-1-xml-input-in-go#answer-32224438
  // @see http://blog.tristanmedia.com/2014/10/using-go-to-parse-non-utf8-xml-feeds/
  decoder.CharsetReader = charset.NewReaderLabel

  var response Response
  err := decoder.Decode(&response)
  if (nil != err) {
    log.Print(err)
  }

  return response, err
}

/**
 * Decodes XML initialization message from DBGp engine.
 */
func decodeInit(xmlInit string) (Init, error) {

  strReader := strings.NewReader(xmlInit)

  decoder := xml.NewDecoder(strReader)
  decoder.CharsetReader = charset.NewReaderLabel

  var init Init
  err := decoder.Decode(&init)
  if (nil != err) {
    log.Print(err)
  }

  return init, err
}
