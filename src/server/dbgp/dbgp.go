/**
 * Package for talking the DBGp protocol.
 */

package dbgp

import (
	"bytes"
	"fmt"
	"io"
	"log"
)

/**
 * Read message from DBGp engine over a network connection.
 *
 * DBGp message format: Int Null XML-Snippet NULL
 */
func Read(connection io.Reader) (msg string, err error) {

	var DBGpMsg bytes.Buffer
	var msgSize int

	count, err := fmt.Fscanf(connection, "%d\x00", &msgSize)
	if err != nil && err != io.EOF {
		log.Print(err)

		return msg, err
	} else if err == io.EOF {
		return msg, err
	} else if 0 == count {
		return msg, err
	}

	copyCount, err := io.CopyN(&DBGpMsg, connection, int64(msgSize))
	if 0 == copyCount || io.EOF == err {
		return msg, err
	} else if nil != err {
		log.Print(err)

		return msg, err
	}

	count, err = fmt.Fscanf(connection, "\x00") // Read null byte.
	if nil != err {
		log.Print(err)

		return msg, err
	}

	msg = DBGpMsg.String()
	return msg, err
}
