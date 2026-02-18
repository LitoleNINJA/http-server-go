package main

import (
	"flag"
	"net"
)

var dir = flag.String("directory", "", "Directory to serve files from")

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		Log.Fatal("Server startup failed", "port", 4221, "error", err.Error())
	}
	defer listener.Close()
	Log.Info("Server started", "address", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			Log.Fatal("Failed to accept connection", "error", err.Error())
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		clientAddr := conn.RemoteAddr().String()

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			Log.Error("Failed to read request", "client", clientAddr, "error", err.Error())
			return
		}

		req, err := parseRequest(string(buf))
		if err != nil {
			Log.Error("Failed to parse request", "client", clientAddr, "error", err.Error())
			return
		}
		Log.Request(clientAddr, req, n)

		if val, ok := req.headers["Connection"]; ok && val == "close" {
			res := newResponse(req.version, StatusOK)
			res.setHeader("Connection", "close")
			res.setTextBody("")

			_, err = conn.Write(res.encode())
			if err != nil {
				Log.Error("Failed to send close response", "client", clientAddr, "error", err.Error())
				return
			}

			Log.Response(clientAddr, res)
			break
		}

		var res httpRes
		switch req.method {
		case "GET":
			res, err = handleGet(req)
			if err != nil {
				Log.Error("GET handler failed", "client", clientAddr, "path", req.path, "error", err.Error())
				return
			}
		case "POST":
			res, err = handlePost(req)
			if err != nil {
				Log.Error("POST handler failed", "client", clientAddr, "path", req.path, "error", err.Error())
				return
			}
		default:
			Log.Warn("Unsupported method", "client", clientAddr, "method", req.method)
			res.version = req.version
			res.status = StatusMethodNotAllowed
		}

		_, err = conn.Write(res.encode())
		if err != nil {
			Log.Error("Failed to send response", "client", clientAddr, "error", err.Error())
			return
		}
		Log.Response(clientAddr, res)
	}
}
