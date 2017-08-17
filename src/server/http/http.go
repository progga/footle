/**
 * Provides the HTTP interface for the debugger.
 */

package http

import (
	"../dbgp/command"
	"../dbgp/message"
	"./file"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const DOCROOT_PATH = "../../ui/"

type client chan<- string

/**
 * List of HTTP clients that are currently listening for Server sent events.
 */
var clientList map[client]bool

/**
 * Initializes current HTTP client list.
 */
func init() {

	clientList = make(map[client]bool)
}

/**
 * Setup HTTP server.
 *
 * Five handlers are used:
 *   - HTTP interface for Footle.
 *   - A file browser for selecting files that will be debugged.
 *   - File content rendered as HTML.
 *   - Debugging command receiver.  This is supposed to be called over Ajax.
 *   - Debugging output sender.  This is supposed to be consumed using
 *     Server sent events.
 *
 * Uses global variable "clientList."
 */
func Listen(codeDir string, port int, out chan string) {

	uiPath, err := determineUIPath()
	if nil != err {
		log.Fatal(err)
	}

	arrival := make(chan client)
	departure := make(chan client)
	go manageClients(clientList, arrival, departure)

	http.Handle("/", http.FileServer(http.Dir(uiPath)))
	// Serve the files that will be debugged.
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(codeDir))))
	// HTML markup for the same files.
	http.HandleFunc("/formatted-file/", makeFormattedFileHandler(port))

	http.HandleFunc("/steering-wheel", makeReceiveHandler(out))
	http.HandleFunc("/message-stream", makeTransmitHandler(arrival, departure))

	address := fmt.Sprintf(":%d", port)
	http.ListenAndServe(address, nil)
}

/**
 * Pass on DBGp command output to HTTP clients.
 *
 * Uses global variable "clientList."
 */
func Tell(in <-chan message.Message) {

	for msg := range in {
		broadcast(fmt.Sprint(msg), clientList)
	}
}

/**
 * Wrapper for the "receive" handler.
 *
 * Apart from the usual arguments for an HTTP handler, it passes a channel
 * to receiver().  This channel can be used to write whatever is received
 * by receive().
 */
func makeReceiveHandler(out chan string) http.HandlerFunc {

	return func(writeStream http.ResponseWriter, request *http.Request) {

		receive(writeStream, request, out)
	}
}

/**
 * Processes HTTP POST calls to the "/steering-wheel" path.
 *
 * Extracts whatever is sent by HTTP clients and tries to prepare a DBGp
 * command out of it.  This command is then written to the output channel so
 * that it can be sent to the DBGp engine.
 */
func receive(writeStream http.ResponseWriter, request *http.Request, out chan string) {

	msg := request.FormValue("msg")

	shortCmd, cmdArgs, err := command.Break(msg)
	if nil != err {
		fmt.Fprintf(writeStream, "%s", err)

		return
	}

	DBGpCmd, err := command.Prepare(shortCmd, cmdArgs)
	if nil != err {
		fmt.Fprintf(writeStream, "%s", err)

		return
	}

	fmt.Fprintf(writeStream, "Got it.")

	out <- DBGpCmd
}

/**
 * Wrapper over transmit().
 *
 * In addition to the usual arguments for an HTTP handler, it passes two
 * channels to transmit().
 */
func makeTransmitHandler(arrival, departure chan client) http.HandlerFunc {

	return func(writeStream http.ResponseWriter, request *http.Request) {

		transmit(writeStream, request, arrival, departure)
	}
}

/**
 * Serves the "/transmit" path.
 *
 * Passes the output of DBGp commands to HTTP clients as Server sent events.
 * Also announces the arrival and departure of each HTTP client.
 *
 * For each client, a new channel is created.  This channel is then passed to
 * the other parts of Footle that writes the output of DBGp commands to this
 * channel.
 */
func transmit(writeStream http.ResponseWriter, request *http.Request, arrival, departure chan client) {

	myEar := make(chan string)
	arrival <- myEar

	flusher, ok := writeStream.(http.Flusher)

	if ok {
		writeStream.Header().Set("Content-Type", "text/event-stream")
		writeStream.Header().Set("Cache-control", "no-cache")
		writeStream.Header().Set("Connection", "keep-alive")

		flusher.Flush()
	} else {
		fmt.Fprintf(writeStream, "Unable to Flush.")

		departure <- myEar
		return
	}

	closedConnectionNotification := writeStream.(http.CloseNotifier).CloseNotify()

	go func() {
		<-closedConnectionNotification
		departure <- myEar
	}()

	for msg := range myEar {
		fmt.Fprintf(writeStream, "data: %s\n\n", msg)
		flusher.Flush()
	}

	// Only relevant when myEar has closed before writeStream.
	fmt.Fprintf(writeStream, "event: close\ndata: The end\n\n")
}

/**
 * Prepare handler for displaying a file as HTML.
 *
 * Split a file into its lines and display them as individual HTML element.
 * Example:
 * <div class="lines">
 *   <pre class="line line__0">&lt?php</pre>
 *   <pre class="line line__1">use Drupal\Core\DrupalKernel;</pre>
 *   <pre class="line line__1">$autoloader = require_once &#39;autoload.php&#39;;</pre>
 *   ...
 * </div>
 */
func makeFormattedFileHandler(port int) http.HandlerFunc {

	return func(writeStream http.ResponseWriter, request *http.Request) {

		filePath := request.URL.Path[len("/formatted-file/"):]

		output, err := file.GrabIt(filePath, port)

		if nil != err {
			http.Error(writeStream, err.Error(), http.StatusInternalServerError)
			return
		}

		writeStream.Header().Set("Content-Type", "text/html")
		io.WriteString(writeStream, output)
	}
}

/**
 * Writes a string message to all current client channels.
 */
func broadcast(msg string, httpClientList map[client]bool) {

	for clientChannel, _ := range httpClientList {
		clientChannel <- msg
	}
}

/**
 * Update the list of HTTP clients currently listening.
 *
 * When an HTTP client first starts listening for Server sent events, we
 * add it as a new client and vice-versa.
 */
func manageClients(httpClientList map[client]bool, arrival, departure <-chan client) {

	for {
		select {
		case clientChannel := <-arrival:
			httpClientList[clientChannel] = true

		case clientChannel := <-departure:
			delete(httpClientList, clientChannel)
			close(clientChannel)
		}
	}
}

/**
 * Find document root of the HTML UI.
 */
func determineUIPath() (uiPath string, err error) {

	binPath, err := os.Executable()
	if nil != err {
		return uiPath, err
	}

	realBinPath, err := filepath.EvalSymlinks(binPath)
	if nil != err {
		return uiPath, err
	}

	uiPath = realBinPath + "/" + DOCROOT_PATH

	return uiPath, err
}
