package main

import (
	"fmt"
	"strings"
)

const (
	StatusOK               = "200 OK"
	StatusCreated          = "201 Created"
	StatusBadRequest       = "400 Bad Request"
	StatusNotFound         = "404 Not Found"
	StatusMethodNotAllowed = "405 Method Not Allowed"
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

func newResponse(version, status string) httpRes {
	return httpRes{
		version: version,
		status:  status,
		headers: make(map[string]string),
	}
}

func (r *httpRes) setHeader(key string, val string) {
	r.headers[key] = val
}

func (r *httpRes) setTextBody(body string) {
	r.body = body
	r.headers["Content-Type"] = "text/plain"
	r.headers["Content-Length"] = fmt.Sprintf("%d", len(body))
}

func (r *httpRes) setFileBody(body string) {
	r.body = body
	r.headers["Content-Type"] = "application/octet-stream"
	r.headers["Content-Length"] = fmt.Sprintf("%d", len(body))
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

func parseRequest(req string) (httpReq, error) {
	lines := strings.Split(req, "\r\n")

	requestLine := strings.Split(lines[0], " ")
	method := strings.TrimSpace(requestLine[0])
	path := strings.TrimSpace(requestLine[1])
	version := strings.TrimSpace(requestLine[2])

	headers := make(map[string]string)
	body := ""
	isBody := false

	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			isBody = true
			continue
		}
		if isBody {
			body += lines[i]
		} else {
			headerParts := strings.Split(lines[i], ":")
			headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
		}
	}

	return httpReq{
		method:  method,
		path:    path,
		version: version,
		headers: headers,
		body:    body,
	}, nil
}
