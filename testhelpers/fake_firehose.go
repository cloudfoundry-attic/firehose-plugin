package testhelpers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

type FakeFirehose struct {
	server *httptest.Server
	lock   sync.Mutex

	validToken string

	lastAuthorization string
	requested         bool

	events       []events.Envelope
	closeMessage []byte
	stayAlive    bool
	wg           sync.WaitGroup
}

func NewFakeFirehose(validToken string) *FakeFirehose {
	return &FakeFirehose{
		validToken:   validToken,
		closeMessage: websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
	}
}

func (f *FakeFirehose) Start() {
	f.server = httptest.NewUnstartedServer(f)
	f.server.Start()
}

func (f *FakeFirehose) Close() {
	f.server.Close()
}

func (f *FakeFirehose) URL() string {
	return fmt.Sprintf("ws://%s", f.server.Listener.Addr().String())
}

func (f *FakeFirehose) LastAuthorization() string {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.lastAuthorization
}

func (f *FakeFirehose) Requested() bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.requested
}

func (f *FakeFirehose) SendLog(logMessage string) {
	f.addEvent(events.Envelope{
		Origin:    proto.String("origin"),
		Timestamp: proto.Int64(1000000000),
		EventType: events.Envelope_LogMessage.Enum(),
		LogMessage: &events.LogMessage{
			Message:     []byte(logMessage),
			MessageType: events.LogMessage_OUT.Enum(),
			Timestamp:   proto.Int64(1000000000),
		},
		Deployment: proto.String("deployment-name"),
		Job:        proto.String("doppler"),
	})
}

func (f *FakeFirehose) addEvent(event events.Envelope) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.events = append(f.events, event)
}

func (f *FakeFirehose) SetCloseMessage(message []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.closeMessage = make([]byte, len(message))
	copy(f.closeMessage, message)
}

func (f *FakeFirehose) KeepConnectionAlive() {
	f.wg.Add(1)
}

func (f *FakeFirehose) CloseAliveConnection() {
	f.wg.Done()
}

func (f *FakeFirehose) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.lastAuthorization = r.Header.Get("Authorization")
	f.requested = true

	if f.lastAuthorization != f.validToken {
		log.Printf("Bad token passed to firehose: %s", f.lastAuthorization)
		rw.WriteHeader(403)
		r.Body.Close()
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}

	ws, _ := upgrader.Upgrade(rw, r, nil)

	defer ws.Close()
	defer ws.WriteControl(websocket.CloseMessage, f.closeMessage, time.Time{})

	for _, envelope := range f.events {
		buffer, _ := proto.Marshal(&envelope)
		err := ws.WriteMessage(websocket.BinaryMessage, buffer)
		if err != nil {
			panic(err)
		}
	}
	f.wg.Wait()
}
