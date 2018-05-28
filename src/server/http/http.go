/**
 * Provides the HTTP interface for the debugger.
 */

package http

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/config"
	"server/core/breakpoint"
	footlecmd "server/core/cmd"
	"server/dbgp/command"
	"server/dbgp/message"
	"server/http/file"
)

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
func Listen(out chan string, config config.Config) {

	codeDir := config.GetCodebase()
	port := config.GetHTTPPort()

	uiPath, err := determineUIPath(config)
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
	http.HandleFunc("/formatted-file/", makeFormattedFileHandler(http.Dir(codeDir)))

	http.HandleFunc("/steering-wheel", makeReceiveHandler(out))
	http.HandleFunc("/message-stream", makeTransmitHandler(arrival, departure))
	http.HandleFunc("/breakpoints", makeBreakpointListingHandler(codeDir))

	address := fmt.Sprintf(":%d", port)
	http.ListenAndServe(address, nil)
}

/**
 * Pass on DBGp command output to HTTP clients.
 *
 * Uses global variable "clientList."
 */
func TellBrowsers(in <-chan message.Message, config config.Config) {

	codeDir := config.DetermineCodeDir()

	for msg := range in {
		adjustedMsg := adjustFilepath(msg, codeDir)
		adjustedMsg.Context.Local = escapeVarValue(msg.Context.Local)

		jsonMsg, err := json.Marshal(adjustedMsg)

		if nil == err {
			broadcast(string(jsonMsg), clientList)
		}
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
func receive(writeStream http.ResponseWriter, request *http.Request, debugger chan string) {

	cmd := request.FormValue("cmd")

	cmdAlias, cmdArgs, err := command.Break(cmd)
	if nil != err {
		fmt.Fprintf(writeStream, "%s", err)

		return
	}

	isFootleCmd := footlecmd.Is(cmdAlias)

	err = command.Validate(cmdAlias, cmdArgs)
	isInvalidDBGpCmd := !isFootleCmd && err != nil

	if isInvalidDBGpCmd {
		fmt.Fprintf(writeStream, "%s", err)

		return
	}

	fmt.Fprintf(writeStream, "Got it.")

	debugger <- cmd
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
func makeFormattedFileHandler(codebase http.Dir) http.HandlerFunc {

	return func(writeStream http.ResponseWriter, request *http.Request) {

		filePath := request.URL.Path[len("/formatted-file/"):]

		output, err := file.GrabIt(codebase, filePath)

		if nil != err {
			http.Error(writeStream, err.Error(), http.StatusInternalServerError)
			return
		}

		writeStream.Header().Set("Content-Type", "text/html")
		io.WriteString(writeStream, output)
	}
}

/**
 * Prepare handler for listing breakpoints.
 *
 * This is useful during UI initialization where only the new UI client needs to
 * process the breakpoint list.  This saves existing UI clients from the burden
 * of reprocessing the breakpoints every time a new HTTP UI client is added.
 */
func makeBreakpointListingHandler(codeDir string) http.HandlerFunc {

	return func(writeStream http.ResponseWriter, request *http.Request) {

		msgWBreakpoint := breakpoint.PrepareFakeMsg()
		adjustedMsg := adjustFilepath(msgWBreakpoint, codeDir)

		if jsonMsg, err := json.Marshal(adjustedMsg); err == nil {
			writeStream.Header().Set("Content-Type", "application/json")
			writeStream.Header().Set("Cache-control", "no-cache")

			fmt.Fprintf(writeStream, "%s", jsonMsg)
		} else {
			fmt.Fprintf(writeStream, "%s", err)
		}
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
func determineUIPath(config config.Config) (uiPath string, err error) {

	uiDocrootPath := config.GetUIPath()
	if filepath.IsAbs(uiDocrootPath) {
		if isDir(uiDocrootPath) {
			return uiDocrootPath, err
		} else {
			err = fmt.Errorf("The %s directory does not exist.", uiDocrootPath)
			return uiDocrootPath, err
		}
	}

	binPath, err := os.Executable()
	if nil != err {
		return uiPath, err
	}

	realBinPath, err := filepath.EvalSymlinks(binPath)
	if nil != err {
		return uiPath, err
	}

	// /bar/baz/footle/bin/footle + .. = /bar/baz/footle/bin
	// So /bar/baz/footle/bin/footle + .. + "../ui" = /bar/baz/footle/ui
	uiPath = filepath.Join(realBinPath, "..", uiDocrootPath)
	if isDir(uiPath) {
		return uiPath, err
	} else {
		err = fmt.Errorf("The %s directory does not exist.", uiPath)
		return uiPath, err
	}
}

/**
 * Does the given directory exist?
 */
func isDir(path string) bool {

	fileinfo, err := os.Stat(path)

	if os.IsNotExist(err) || !fileinfo.IsDir() {
		return false
	}

	return true
}

/**
 * Set filepath relative to codebase.
 *
 * HTTP clients are always given relative filepaths whereas the DBGp engine
 * deals with absolute filepaths.  Here we convert absolute filepaths to
 * relative.
 *
 * Filepaths are present in:
 *  - response.Properties.Filename
 *  - response.Breakpoints
 *  - response.Stacktrace
 */
func adjustFilepath(response message.Message, codeDir string) message.Message {

	codeDirUri := "file://" + codeDir

	// @todo filepath.HasPrefix() is deprecated.  Replace when a suitable
	// replacement is found.
	hasFilename := filepath.HasPrefix(response.Properties.Filename, codeDirUri)
	hasBreakpoints := len(response.Breakpoints) > 0
	hasStacktrace := len(response.Stacktrace) > 0

	// Adjust response.Properties.Filename
	if hasFilename {
		relativePath, err := filepath.Rel(codeDirUri, response.Properties.Filename)

		if nil == err {
			response.Properties.Filename = relativePath
		}
	}

	// Now adjust response.Breakpoints
	//
	// Modify a *copy* of the breakpoint list.  Otherwise it will modify the
	// original message too.  This is because the breakpoint list is a map
	// *reference* and not a copy.
	var adjustedBreakpoints map[int]message.Breakpoint
	if hasBreakpoints {
		adjustedBreakpoints = make(map[int]message.Breakpoint)

		for breakpointId, breakpoint := range response.Breakpoints {
			relativePath, err := filepath.Rel(codeDirUri, breakpoint.Filename)

			if nil == err {
				breakpoint.Filename = relativePath
				adjustedBreakpoints[breakpointId] = breakpoint
			}
		}
		response.Breakpoints = adjustedBreakpoints
	}

	// Lastly, adjust response.Stacktrace
	var adjustedStacktrace []message.StackLevel
	if hasStacktrace {
		adjustedStacktrace = []message.StackLevel{}

		for _, StackLevel := range response.Stacktrace {
			relativePath, err := filepath.Rel(codeDirUri, StackLevel.Filename)

			if nil == err {
				StackLevel.Filename = relativePath
				adjustedStacktrace = append(adjustedStacktrace, StackLevel)
			}
		}
		response.Stacktrace = adjustedStacktrace
	}

	return response
}

/**
 * HTML escapse variable values.
 *
 * Let's not burden HTTP clients with HTML escaping.
 */
func escapeVarValue(vars []message.Variable) (variables []message.Variable) {

	if len(vars) == 0 {
		return
	}

	var escapedVar message.Variable

	for _, varDetails := range vars {
		escapedVar = varDetails
		escapedVar.Value = html.EscapeString(varDetails.Value)

		if len(varDetails.Children) > 0 {
			escapedVar.Children = escapeVarValue(varDetails.Children)
		}

		variables = append(variables, escapedVar)
	}

	return
}
