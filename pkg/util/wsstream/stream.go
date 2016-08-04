/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package wsstream

import (
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
	"k8s.io/kubernetes/pkg/util/runtime"
)

// The WebSocket subprotocol "binary.k8s.io" will only send messages to the
// client and ignore messages sent to the server. The received messages are
// the exact bytes written to the stream. Zero byte messages are possible.
const binaryWebSocketProtocol = "binary.k8s.io"

// The WebSocket subprotocol "base64.binary.k8s.io" will only send messages to the
// client and ignore messages sent to the server. The received messages are
// a base64 version of the bytes written to the stream. Zero byte messages are
// possible.
const base64BinaryWebSocketProtocol = "base64.binary.k8s.io"

// Reader supports returning an arbitrary byte stream over a websocket channel.
// Supports the "binary.k8s.io" and "base64.binary.k8s.io" subprotocols.
type Reader struct {
	err                chan error
	r                  io.Reader
	ping               bool
	timeout            time.Duration
	supportedProtocols []string
	protocolPrefix     string
}

// NewReader creates a WebSocket pipe that will copy the contents of r to a provided
// WebSocket connection. If ping is true, a zero length message will be sent to the client
// before the stream begins reading.
//
// One of the base64 ("base64.channel.k8s.io") and raw ("channel.k8s.io") protocol is negotiated
// during initial handshake. Given protocolPrefixes (being DNS 1123 subdomains or empty string)
// are provided all combinations of those and base64+raw are nogatiated, e.g. "", "v1", "v2" leads
// to the protocols "base64.channel.k8s.io", "channel.k8s.io", "v1.base64.channel.k8s.io",
// "v1.channel.k8s.io", "v2.base64.channel.k8s.io", "v2.channel.k8s.io".
func NewReader(r io.Reader, ping bool, protocolPrefixes []string) *Reader {
	return &Reader{
		r:                  r,
		err:                make(chan error),
		ping:               ping,
		supportedProtocols: generateProtocols(protocolPrefixes, []string{binaryWebSocketProtocol, base64BinaryWebSocketProtocol}),
	}
}

// SetIdleTimeout sets the interval for both reads and writes before timeout. If not specified,
// there is no timeout on the reader.
func (r *Reader) SetIdleTimeout(duration time.Duration) {
	r.timeout = duration
}

func (r *Reader) handshake(config *websocket.Config, req *http.Request) error {
	return handshake(config, req, r.supportedProtocols)
}

// Copy the reader to the response. The created WebSocket is closed after this
// method completes.
func (r *Reader) Copy(w http.ResponseWriter, req *http.Request) error {
	go func() {
		defer runtime.HandleCrash()
		websocket.Server{Handshake: r.handshake, Handler: r.handle}.ServeHTTP(w, req)
	}()
	return <-r.err
}

// ProtocolPrefix returns the prefix of the negotiated protocol.
func (r *Reader) ProtocolPrefix() string {
	return r.protocolPrefix
}

// handle implements a WebSocket handler.
func (r *Reader) handle(ws *websocket.Conn) {
	encode := false
	if len(ws.Config().Protocol) > 0 {
		prefix, suffix := splitProtocol(ws.Config().Protocol[0], base64BinaryWebSocketProtocol, binaryWebSocketProtocol)
		r.protocolPrefix = prefix
		encode = suffix == base64BinaryWebSocketProtocol
	}
	defer close(r.err)
	defer ws.Close()
	go IgnoreReceives(ws, r.timeout)
	r.err <- messageCopy(ws, r.r, encode, r.ping, r.timeout)
}

func resetTimeout(ws *websocket.Conn, timeout time.Duration) {
	if timeout > 0 {
		ws.SetDeadline(time.Now().Add(timeout))
	}
}

func messageCopy(ws *websocket.Conn, r io.Reader, base64Encode, ping bool, timeout time.Duration) error {
	buf := make([]byte, 2048)
	if ping {
		resetTimeout(ws, timeout)
		if err := websocket.Message.Send(ws, []byte{}); err != nil {
			return err
		}
	}
	for {
		resetTimeout(ws, timeout)
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if n > 0 {
			if base64Encode {
				if err := websocket.Message.Send(ws, base64.StdEncoding.EncodeToString(buf[:n])); err != nil {
					return err
				}
			} else {
				if err := websocket.Message.Send(ws, buf[:n]); err != nil {
					return err
				}
			}
		}
	}
}
