package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type httpReq struct {
	method  string
	path    string
	version string
	headers map[string]string
}

type httpRes struct {
	version string
	status  string
	headers map[string]string
	body    string
}

func (r httpRes) encode() []byte {
	var res string
	res += r.version + " " + r.status + "\r\n"
	for k, v := range r.headers {
		res += k + ": " + v + "\r\n"
	}
	res += "\r\n"
	res += r.body
	return []byte(res)
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on: ", l.Addr().String())

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			return
		}
		fmt.Printf("\nRecieved: %s Bytes: %d\n", string(buf[:n]), n)

		req, err := parseRequest(string(buf))
		if err != nil {
			fmt.Println("Error parsing request: ", err.Error())
			return
		}

		var res httpRes
		if req.path == "/" {
			res.status = "200 OK"
		} else if pre, ok := CutPrefix(req.path, "/echo/"); ok {
			res = httpRes{
				version: req.version,
				status:  "200 OK",
				headers: map[string]string{
					"Content-Type":   "text/plain",
					"Content-Length": fmt.Sprintf("%d", len(pre)),
				},
				body: pre,
			}
		} else {
			res.status = "404 Not Found"
		}

		_, err = conn.Write(res.encode())
		if err != nil {
			fmt.Println("Error writing: ", err.Error())
			return
		}
		fmt.Println("Sent: ", res)
	}
}

func parseRequest(req string) (httpReq, error) {
	a := strings.Split(req, "\r\n")

	method := strings.TrimSpace(strings.Split(a[0], " ")[0])
	path := strings.TrimSpace(strings.Split(a[0], " ")[1])
	version := strings.TrimSpace(strings.Split(a[0], " ")[2])
	headers := make(map[string]string)
	for i := 1; i < len(a); i++ {
		if a[i] == "" {
			break
		}
		h := strings.Split(a[i], ":")
		headers[strings.TrimSpace(h[0])] = strings.TrimSpace(h[1])
	}
	// fmt.Printf("Method: %s\nPath: %s\nVersion: %s\nHeaders: %v\n", method, path, version, headers)

	return httpReq{
		method:  method,
		path:    path,
		version: version,
		headers: headers,
	}, nil
}
