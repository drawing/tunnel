package tools

import (
	// "bytes"
	"code.google.com/p/go.net/websocket"
	// "encoding/base64"
	// "encoding/binary"
	"errors"
	// "io"
	"log"
	"net"
	"net/http"
)

type WebSocketListener struct {
	accept chan *WebSocketConn
}

type WebSocketConn struct {
	*websocket.Conn
	finish chan byte
	buffer []byte
}

/*
func (w *WebSocketConn) Write(b []byte) (n int, err error) {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, uint16(len(b)))

	header := []byte(base64.StdEncoding.EncodeToString(buffer.Bytes()))
	body := []byte(base64.StdEncoding.EncodeToString(b))

	_, err = w.Conn.Write(header)
	if err != nil {
		return 0, err
	}

	_, err = w.Conn.Write(body)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

func (w *WebSocketConn) Read(b []byte) (n int, err error) {
	if len(w.buffer) == 0 {
		w.buffer = make([]byte, MaxBufferSize*2)

		_, err = io.ReadFull(w.Conn, w.buffer[:4])
		if err != nil {
			return 0, err
		}

		header, err := base64.StdEncoding.DecodeString(string(w.buffer[:4]))
		if err != nil {
			return 0, err
		}

		var length uint16
		binary.Read(bytes.NewBuffer(header[:2]), binary.BigEndian, &length)

		if int(length) >= len(w.buffer) {
			return 0, errors.New("package error")
		}

		_, err = io.ReadFull(w.Conn, w.buffer[0:length])
		if err != nil {
			return 0, err
		}
		w.buffer = w.buffer[0:length]
	}

	length := 0
	for length < len(b) && length < len(w.buffer) {
		b[length] = w.buffer[length]
		length++
	}

	if length < len(w.buffer) {
		w.buffer = w.buffer[length:]
	} else {
		w.buffer = []byte{}
	}

	return length, nil
}
*/
func WebSocketDial(network, address string) (net.Conn, error) {
	if network != "ws" && network != "wss" {
		return nil, errors.New("only support ws and wss")
	}

	origin := "http://golang.org/"
	url := network + "://" + address + "/websocket_secret"

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func WebSocketListen(net, laddr string) (net.Listener, error) {
	if net != "ws" && net != "wss" {
		return nil, errors.New("only support ws")
	}

	acc := make(chan *WebSocketConn, 10)
	handler := &WebSocketHandler{acc}

	http.Handle("/websocket_secret", websocket.Handler(handler.ServeHTTP))

	go http.ListenAndServe(laddr, nil)

	return &WebSocketListener{acc}, nil
}

func (l *WebSocketListener) Accept() (c net.Conn, err error) {
	ws, ok := <-l.accept
	if !ok {
		return nil, errors.New("EOF")
	}

	return ws, nil
}

func (l *WebSocketListener) Close() error {
	close(l.accept)
	return nil
}

func (l *WebSocketListener) Addr() net.Addr {
	log.Fatalln("HTTPListener Addr not impement")
	return nil
}

/// http handler
type WebSocketHandler struct {
	accept chan *WebSocketConn
}

func (l *WebSocketHandler) ServeHTTP(ws *websocket.Conn) {
	ch := make(chan byte)
	sock := &WebSocketConn{ws, ch, []byte{}}
	l.accept <- sock
	<-ch
}
