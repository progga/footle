/**
 * Package for talking to the DBGp engine.
 */

package core

import (
	"server/config"
	"log"
	"net"
	"strconv"
)

/**
 * Start listening on the standard DBGp port of 9000.
 */
func ListenForDBGpEngine(config config.Config) (sock net.Listener) {

	DBGpPort := config.GetDBGpPort()
	address := ":" + strconv.Itoa(DBGpPort)

	sock, err := net.Listen("tcp", address)
	if nil != err {
		log.Fatal(err)
	}

	return sock
}
