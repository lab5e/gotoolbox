package rest

//
//Copyright 2018 Telenor Digital AS
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ExploratoryEngineering/logging"
	"golang.org/x/net/websocket"
)

// Timeout for keepalive messages
const timeoutSeconds = 30

// WebsocketStreamer is an adapter type for websocket streams. They all follow
// a similar pattern: Set up request, get a channel for (more or less) realtime
// data and push data to the client whenever the channel gets something.
// A keepalive message is sent at regular intervals if there's no data.
type WebsocketStreamer interface {
	Setup(r *http.Request) error
	Cleanup()
	Input() <-chan interface{}
	KeepaliveMessage() interface{}
}

// DefaultKeepAliveMessage is the default KeepAlive-message
func DefaultKeepAliveMessage() interface{} {
	return struct {
		KeepAlive bool `json:"keepAlive"`
	}{true}
}

// WebsocketHandler generates a http.HandlerFunc from the websocket adapter
func WebsocketHandler(streamer WebsocketStreamer) http.HandlerFunc {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		if err := streamer.Setup(ws.Request()); err != nil {
			// write error and return
			logging.Warning("Setup error: %v", err)
			return
		}
		defer streamer.Cleanup()
		ch := streamer.Input()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				buf, err := json.Marshal(msg)
				if err != nil {
					logging.Warning("Got error marshalling message %+v: %v", msg, err)
					return
				}
				_, err = ws.Write(buf)
				if err != nil {
					logging.Warning("Error writing. Exiting: %v", err)
					return
				}
			case <-time.After(timeoutSeconds * time.Second):
				msg := streamer.KeepaliveMessage()
				if msg == nil {
					continue
				}
				buf, err := json.Marshal(&msg)
				if err != nil {
					return
				}
				_, err = ws.Write(buf)
				if err != nil {
					logging.Warning("Error writing. Exiting: %v", err)
					return
				}
			}
		}
	}).ServeHTTP
}
