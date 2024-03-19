package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type httpReq struct {
	method  string
	path    string
	version string
	headers map[string]string
	body    string
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

var dir = flag.String("directory", "", "Directory to serve files from")

func main() {
	flag.Parse()

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
	switch req.method {
	case "GET":
		res, err = handleGet(req)
		if err != nil {
			fmt.Println("Error handling GET: ", err.Error())
			return
		}
	case "POST":
		res, err = handlePost(req)
		if err != nil {
			fmt.Println("Error handling POST: ", err.Error())
			return
		}
	default:
		res.status = "405 Method Not Allowed"
	}

	_, err = conn.Write(res.encode())
	if err != nil {
		fmt.Println("Error writing: ", err.Error())
		return
	}
	fmt.Println("Sent: ", res)
}

func parseRequest(req string) (httpReq, error) {
	a := strings.Split(req, "\r\n")

	method := strings.TrimSpace(strings.Split(a[0], " ")[0])
	path := strings.TrimSpace(strings.Split(a[0], " ")[1])
	version := strings.TrimSpace(strings.Split(a[0], " ")[2])
	headers := make(map[string]string)
	body := ""
	isBody := false
	for i := 1; i < len(a); i++ {
		if a[i] == "" {
			isBody = true
			continue
		}
		if isBody {
			body += a[i]
		} else {
			h := strings.Split(a[i], ":")
			headers[strings.TrimSpace(h[0])] = strings.TrimSpace(h[1])
		}
	}

	// fmt.Printf("Method: %s\nPath: %s\nVersion: %s\nHeaders: %v\n Body: %s\n", method, path, version, headers, body)

	return httpReq{
		method:  method,
		path:    path,
		version: version,
		headers: headers,
		body:    body,
	}, nil
}

func handleGet(req httpReq) (httpRes, error) {
	var res httpRes
	res.version = req.version

	if req.path == "/" {
		res.status = "200 OK"
	} else if pre, ok := CutPrefix(req.path, "/echo/"); ok {
		res.status = "200 OK"
		res.headers = map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(pre)),
		}
		res.body = pre
	} else if strings.HasPrefix(req.path, "/user-agent") {
		res.status = "200 OK"
		res.body = req.headers["User-Agent"]
		res.headers = map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": fmt.Sprintf("%d", len(res.body)),
		}
	} else if filename, ok := CutPrefix(req.path, "/files/"); ok {
		file, err := os.Open(*dir + "/" + filename)
		if err != nil {
			res.status = "404 Not Found"
		} else {
			data, err := io.ReadAll(file)
			if err != nil {
				fmt.Println("Error reading file: ", err.Error())
				return httpRes{}, err
			}
			res.status = "200 OK"
			res.body = string(data)
			res.headers = map[string]string{
				"Content-Type":   "application/octet-stream",
				"Content-Length": fmt.Sprintf("%d", len(res.body)),
			}
		}
		file.Close()
	} else {
		res.status = "404 Not Found"
	}

	return res, nil
}

func handlePost(req httpReq) (httpRes, error) {
	var res httpRes
	res.version = req.version
	res.status = "201 Created"

	filename, ok := CutPrefix(req.path, "/files/")
	if !ok {
		res.status = "400 Bad Request"
	} else {
		file, err := os.Create(*dir + "/" + filename)
		if err != nil {
			fmt.Println("Error creating file: ", err.Error())
			return httpRes{}, err
		}
		defer file.Close()

		data := req.body
		len, _ := strconv.Atoi(req.headers["Content-Length"])
		_, err = file.WriteString(data[:len])
		if err != nil {
			fmt.Println("Error writing to file: ", err.Error())
			return httpRes{}, err
		}
	}
	return res, nil
}
