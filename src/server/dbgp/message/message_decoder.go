
package message

import (
  "encoding/xml"
  "golang.org/x/net/html/charset"
  "log"
  "strings"
)

/**
 * Parses XDebug's XML responses.
 *
 * @param string xmlResponse
 * @return Response
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
 * Parses XDebug's XML initialization message.
 *
 * @param string xmlResponse
 * @return Response
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
