package rest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

type testStreamer struct {
}

func (t *testStreamer) Setup(r *http.Request) error {
	return nil
}

func (t *testStreamer) Cleanup() {

}

func (t *testStreamer) Input() <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch <- "foo"
		close(ch)
	}()
	return ch
}

func (t *testStreamer) KeepaliveMessage() interface{} {
	return struct {
		Type string `json:"type"`
	}{"keepAlive"}
}

func TestWebsocketHandler(t *testing.T) {
	server := httptest.NewServer(WebsocketHandler(&testStreamer{}))

	defer server.Close()

	conn, err := websocket.Dial(strings.Replace(server.URL, "http", "ws", -1), "ws", server.URL)
	if err != nil {
		t.Fatal("Got error requesting web socket: ", err)
	}
	defer conn.Close()
	for err == nil {
		var buf [1024]byte
		_, err = conn.Read(buf[:])
	}
}
