/**
 * Wrapper around network socket.
 *
 * Encapsulates network socket and its associated network connection.  This
 * makes it easy to turn the socket on and off whenever needed.  Any established
 * network connection can also be easily dropped.
 *
 * We can also block execution when the socket has been turned off.
 *
 * The network connection data structure is maintained as a Singleton.
 */

package connection

import "log"
import "net"
import "server/config"
import "strconv"

type Connection struct {
	sock       net.Listener
	connection net.Conn

	isActive bool
	wait     chan bool

	config config.Config

	initialized bool
}

/**
 * This is used to maintain a single instance of the network connection
 * data structure.
 */
var connection Connection

/**
 * Initialize and/or return single instance of the connection structure.
 */
func GetConnection() *Connection {

	if connection.initialized {
		return &connection
	}

	config := config.Get()
	connection = Connection{
		sock:        nil,
		connection:  nil,
		isActive:    false,
		wait:        nil,
		config:      config,
		initialized: true,
	}

	return &connection
}

/**
 * Block execution when the network socket is not listening.
 */
func (c *Connection) WaitUntilActive() {

	if c.isActive {
		return
	}

	<-c.wait
}

/**
 * Activate the network socket.
 *
 * Close the Connection.wait channel to indicate that the socket has been
 * initialized.
 */
func (c *Connection) Activate() {

	if c.isActive {
		return
	}

	c.signalActivation()
	c.startListeningForDBGpEngine()
}

/**
 * Deactivate the socket and any established connection.
 */
func (c *Connection) Deactivate() {

	if !c.isActive {
		return
	}

	c.isActive = false
	c.wait = make(chan bool)

	c.Disconnect()
	c.stopListening()
}

/**
 * Establish connection with a DBGp engine.
 */
func (c *Connection) Connect() *net.Conn {

	conn, err := c.sock.Accept()

	if err != nil {
		log.Println(err)
	}

	c.connection = conn
	return &c.connection
}

/**
 * Drop any established network connection.
 */
func (c *Connection) Disconnect() error {

	var err error

	if c.connection == nil {
		return err
	}

	err = c.connection.Close()

	return err
}

/**
 * Are we talking to a DBGp engine?
 *
 * We go on air once a DBGp engine connects to Footle.
 *
 * Write an empty byte array to test if a connection has been established.
 */
func (c *Connection) IsOnAir() bool {

	ignore := []byte{}

	if nil == c.connection {
		return false
	}

	if readCount, err := c.connection.Write(ignore); nil != err {
		_ = readCount
		return false
	}

	return true
}

/**
 * Return an instance of the network connection, active or not.
 */
func (c *Connection) Get() *net.Conn {

	return &c.connection
}

/**
 * End wait by WaitUntilActive().
 */
func (c *Connection) signalActivation() {

	if c.isActive {
		return
	}

	c.isActive = true

	if c.wait != nil {
		close(c.wait)
	}
}

/**
 * Start listening for the DBGp engine.
 *
 * Listen on a port (default 9003) where the DBGp engine is expected to knock.
 */
func (c *Connection) startListeningForDBGpEngine() {

	DBGpPort := c.config.GetDBGpPort()
	address := ":" + strconv.Itoa(DBGpPort)

	sock, err := net.Listen("tcp", address)
	if nil != err {
		log.Fatal(err)
	}

	c.sock = sock
}

/**
 * Stop listening for the DBGp engine.
 */
func (c *Connection) stopListening() {

	if c.sock == nil {
		return
	}

	c.sock.Close()
}
