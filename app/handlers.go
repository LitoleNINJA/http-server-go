package main

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
)

func handleGet(req httpReq) (httpRes, error) {
	res := newResponse(req.version, StatusOK)
	if val, ok := req.headers["Accept-Encoding"]; ok {
		switch val {
		case "gzip":
			res.setHeader("Content-Encoding", val)

		default:
			Log.Error("Encoding not supported.", "Client encoding : ", val)
		}
	}
	switch {
	case req.path == "/":
		break

	case HasPrefix(req.path, "/echo/"):
		message, _ := CutPrefix(req.path, "/echo/")
		res.setTextBody(message)

	case HasPrefix(req.path, "/user-agent"):
		userAgent := req.headers["User-Agent"]
		res.setTextBody(userAgent)

	case HasPrefix(req.path, "/files/"):
		return handleGetFile(req)
	default:
		res.status = StatusNotFound
	}

	return res, nil
}

// handles GET requests for "/files/{filename}"
func handleGetFile(req httpReq) (httpRes, error) {
	filename, _ := CutPrefix(req.path, "/files/")
	filePath := filepath.Join(*dir, filename)

	file, err := os.Open(filePath)
	if err != nil {
		Log.Warn("File not found", "path", filePath)
		return newResponse(req.version, StatusNotFound), nil
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		Log.Error("Failed to read file", "path", filePath, "error", err.Error())
		return httpRes{}, err
	}

	Log.Debug("Serving file", "path", filePath, "size", len(data))
	res := newResponse(req.version, StatusOK)
	res.setFileBody(string(data))
	return res, nil
}

func handlePost(req httpReq) (httpRes, error) {
	if !HasPrefix(req.path, "/files/") {
		return newResponse(req.version, StatusBadRequest), nil
	}
	return handlePostFile(req)
}

// handles POST requests for "/files/{filename}"
func handlePostFile(req httpReq) (httpRes, error) {
	filename, _ := CutPrefix(req.path, "/files/")
	filePath := filepath.Join(*dir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		Log.Error("Failed to create file", "path", filePath, "error", err.Error())
		return httpRes{}, err
	}
	defer file.Close()

	contentLen, _ := strconv.Atoi(req.headers["Content-Length"])
	content := req.body[:contentLen]

	if _, err = file.WriteString(content); err != nil {
		Log.Error("Failed to write to file", "path", filePath, "error", err.Error())
		return httpRes{}, err
	}

	Log.Debug("File created", "path", filePath, "size", contentLen)
	return newResponse(req.version, StatusCreated), nil
}
