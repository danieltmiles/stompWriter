package stompWriter

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/franela/goblin"
	"github.com/gmallard/stompngo"
	"github.com/monsooncommerce/mockstomp"
	. "github.com/onsi/gomega"
)

func testServer(serverGivesThis string) (*httptest.Server, *int) {
	var server *httptest.Server

	serverGivesBadStatus := false
	serverHitCount := 0
	serverGivesTheseBytes := []byte(serverGivesThis)

	server = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if serverGivesBadStatus {
					w.WriteHeader(http.StatusBadRequest)
				}
				serverHitCount++
				w.Write(serverGivesTheseBytes)
			}))
	return server, &serverHitCount
}

func makeMockStomp() *mockstomp.MockStompConnection {
	mockStomp := mockstomp.MockStompConnection{}
	return &mockStomp
}

func TestStompWriter(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	var server *httptest.Server
	var mockStomp *mockstomp.MockStompConnection

	g.Describe("Initialization", func() {
		g.BeforeEach(func() {
			server, _ = testServer("200 OK")
			mockStomp = makeMockStomp()
			getStompConnection = func(netCon net.Conn, connectionHeaders stompngo.Headers) (stompConnectioner, error) {
				return mockStomp, nil
			}
		})
		g.It("should be initialized correctly", func() {
			hostname, port, _ := net.SplitHostPort(server.URL[7:])
			username := "username"
			password := "password"
			queueName := "queueName"

			stompWriter, _ := New(hostname, port, username, password, queueName)
			Expect(stompWriter.hostname).To(Equal(hostname))
			Expect(stompWriter.port).To(Equal(port))
			Expect(stompWriter.password).To(Equal(password))
			Expect(stompWriter.queueName).To(Equal(queueName))

			expectedHeaders := stompngo.Headers{
				"accept-version", "1.1",
				"login", username,
				"passcode", password,
				"host", hostname,
				"heart-beat", "5000,5000",
			}
			Expect(stompWriter.connectionHeaders).To(Equal(expectedHeaders))

			Expect(stompWriter.Connection).To(Equal(mockStomp))
		})
		g.It("should fail when configured improperly", func() {
			hostname := ""
			port := "999"
			username := "username"
			password := "password"
			queueName := "queueName"

			stompWriter, err := Configure(hostname, port, username, password, queueName, "myAppName")
			Expect(stompWriter).To(Equal((*StompWriter)(nil)))
			Expect(err.Error()).To(Equal("Logger configuration not properly set"))
		})
		g.It("should send request properly", func() {
			hostname, port, _ := net.SplitHostPort(server.URL[7:])
			username := "username"
			password := "password"
			queueName := "queueName"

			stompWriter, _ := New(hostname, port, username, password, queueName)

			logLine := "logLinelogLinelogLine"
			stompWriter.Write([]byte(logLine))

			Expect(len(mockStomp.MessagesSent)).To(Equal(1))
			msg := <-mockStomp.MessagesSent

			expectedMessage := mockstomp.MockStompMessage{
				Order:   0,
				Headers: []string{"persistent", "true", "destination", "/queue/queueName", "content-type", "text/plain;charset=UTF-8"},
				Message: logLine,
			}

			Expect(msg).To(Equal(expectedMessage))
		})
		g.It("should send disconnect properly", func() {
			hostname, port, _ := net.SplitHostPort(server.URL[7:])
			username := "username"
			password := "password"
			queueName := "queueName"

			stompWriter, _ := New(hostname, port, username, password, queueName)

			stompWriter.netCon = nil
			// verify disconnect not called when netCon is nil
			Expect(mockStomp.DisconnectCalled).NotTo(Equal(true))
			stompWriter.Disconnect()
			Expect(mockStomp.DisconnectCalled).NotTo(Equal(true))

			stompWriter, _ = New(hostname, port, username, password, queueName)
			// with non-nil netCon, should have been called
			Expect(mockStomp.DisconnectCalled).NotTo(Equal(true))
			stompWriter.Disconnect()
			Expect(mockStomp.DisconnectCalled).To(Equal(true))
		})
	})
}
