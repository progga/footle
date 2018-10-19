/**
 * Tests for the HTTP interface.
 *
 * We *must* run "go generate http.go" before trying the tests in this file.
 */

package http

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

/**
 * Mock ResponseRecorder.
 *
 * httptest.ResponseRecorder does not implement http.CloseNotifier.  So we use
 * a wrapper over httptest.ResponseRecorder instead.  This wrapper satisfies
 * http.CloseNotifier.  This response recorder is useful where the HTTP handler
 * needs to know as soon as the HTTP connection closes.  This is particularly
 * useful for long running HTTP connections as is served by transmit().
 */
type mockResponseRecorder struct {
	*httptest.ResponseRecorder
	closeNotify chan bool
}

func (writer mockResponseRecorder) CloseNotify() <-chan bool { return writer.closeNotify }
func (writer mockResponseRecorder) Close()                   { writer.closeNotify <- true }
func (writer mockResponseRecorder) Flush()                   {}

/**
 * Tests for receive().
 *
 * Tests both the output text and the extracted command name coming out of
 * receive().
 */
func TestReceive(t *testing.T) {

	// Fail case. "foo" is an invalid command.
	formValues := url.Values{"cmd": {"foo"}}
	formReader := strings.NewReader(formValues.Encode())
	request := httptest.NewRequest("POST", "/steering-wheel", formReader)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	writer := httptest.NewRecorder()
	commands := make(chan string)

	receive(writer, request, commands)

	expectedResponse := "Unknown command."
	response := writer.Body.String()

	if expectedResponse != response {
		t.Errorf("receive(foo) says: %s", response)
	}

	// Pass case that receives the "status" command.
	formValues = url.Values{"cmd": {"status"}}
	formReader = strings.NewReader(formValues.Encode())
	request = httptest.NewRequest("POST", "/steering-wheel", formReader)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; param=value")
	writer = httptest.NewRecorder()
	commands = make(chan string)

	go receive(writer, request, commands)
	DBGpCmd := <-commands

	expectedResponse = "Got it."
	response = writer.Body.String()

	if expectedResponse != response {
		t.Errorf("receive(status) says: %s", response)
	}

	expectedCmd := "status"
	if expectedCmd != DBGpCmd {
		t.Errorf("receive(status) commanded: %s", DBGpCmd)
	}
}

/**
 * Tests for transmit().
 *
 * Tests the whole sequence of events for our Server-sent-event server.
 */
func TestTransmit(t *testing.T) {

	request := httptest.NewRequest("GET", "/message-stream", nil)
	writer := mockResponseRecorder{
		ResponseRecorder: httptest.NewRecorder(),
		closeNotify:      make(chan bool),
	}

	// HTTP client has arrived.
	arrival := make(chan client)
	departure := make(chan client)
	go transmit(writer, request, arrival, departure)

	ear0 := <-arrival
	ear0 <- "Foo bar."

	time.Sleep(time.Millisecond)
	streamedOutput := writer.Body.String()

	expected := "data: Foo bar.\n\n"
	if expected != streamedOutput {
		t.Errorf("transmit() said: %s", streamedOutput)
	}

	// HTTP client has departed.  Its channel should be sent for removal from
	// client list...
	writer.Close()
	time.Sleep(time.Millisecond)

	earChannel := <-departure
	if ear0 != earChannel {
		t.Error("transmit() did not deal with channel at departure.")
	}

	// ...and a closing message has been sent.
	close(ear0)
	time.Sleep(time.Millisecond)
	streamedOutput = writer.Body.String()

	expected = "data: Foo bar.\n\nevent: close\ndata: The end\n\n"
	if expected != streamedOutput {
		t.Errorf("At the end, transmit() said: %s", streamedOutput)
	}
}

/**
 * Tests for broadcast().
 */
func TestBroadcast(t *testing.T) {

	httpClientList := make(map[client]bool)

	ear0 := make(chan string)
	ear1 := make(chan string)
	ear2 := make(chan string)

	httpClientList[ear0] = true
	httpClientList[ear1] = true
	httpClientList[ear2] = true

	go broadcast("Foo", httpClientList)

	var msg0, msg1, msg2 string

	// Record what we have just heard.
	for i := 0; i < 3; i++ {
		select {
		case msg0 = <-ear0:

		case msg1 = <-ear1:

		case msg2 = <-ear2:
		}
	}

	if "Foo" != msg0 || "Foo" != msg1 || "Foo" != msg2 {
		t.Errorf("Wrong broadcast: %s, %s, %s", msg0, msg1, msg2)
	}
}

/**
 * Tests for manageClients().
 *
 * Tests arrival and departure recording of HTTP clients.
 */
func TestManageClients(t *testing.T) {

	httpClientList := make(map[client]bool)
	arrival := make(chan client)
	departure := make(chan client)

	go manageClients(httpClientList, arrival, departure)

	ear0 := make(chan string)
	ear1 := make(chan string)
	ear2 := make(chan string)

	arrival <- ear0
	// Sleep() is needed to give the manageClients() goroutine a chance update
	// the client list.
	time.Sleep(time.Millisecond)
	if 1 != len(httpClientList) {
		t.Error("manageClients() failed to record arrival.")
	}

	departure <- ear0
	time.Sleep(time.Millisecond)
	if 0 != len(httpClientList) {
		t.Error("manageClients() failed to record departure.")
	}

	arrival <- ear1
	arrival <- ear2
	time.Sleep(time.Millisecond)
	if 2 != len(httpClientList) {
		t.Error("manageClients() failed to record two arrivals.")
	}

	departure <- ear1
	departure <- ear2
	time.Sleep(time.Millisecond)
	if 0 != len(httpClientList) {
		t.Error("manageClients() failed to record two departures.")
	}
}
