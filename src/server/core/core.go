/**
 * Package for talking to the DBGp engine.
 */

package core

import (
	conn "server/core/connection"
)

/**
 * Process commands coming from UIs.
 *
 * Some commands (e.g. run) are meant for the DBGp engine.  These are forwarded
 * to the appropriate channel.  Other commands (e.g. on) are meant to control
 * Footle's behavior.  These are acted up on.
 */
func ProcessUICmds(CmdsFromUI, DBGpCmds chan string, DBGpConnection *conn.Connection) {

	for cmd := range CmdsFromUI {
		if cmd == "on" {
			DBGpConnection.Activate()
		} else if cmd == "off" {
			DBGpConnection.Deactivate()
		} else if cmd == "continue" {
			DBGpConnection.Disconnect()
		} else {
			DBGpCmds <- cmd
		}
	}
}
