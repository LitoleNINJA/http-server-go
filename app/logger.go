package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Logger provides structured, detailed logging
type Logger struct{}

// Log is the global logger instance
var Log = &Logger{}

func (l *Logger) timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (l *Logger) separator() string {
	return strings.Repeat("=", 70)
}

func (l *Logger) subSeparator() string {
	return strings.Repeat("-", 70)
}

// Info logs informational messages
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.log("INFO ", msg, keyvals...)
}

// Warn logs warning messages
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.log("WARN ", msg, keyvals...)
}

// Error logs error messages
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.log("ERROR", msg, keyvals...)
}

// Debug logs debug messages
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.log("DEBUG", msg, keyvals...)
}

// Fatal logs error and exits
func (l *Logger) Fatal(msg string, keyvals ...interface{}) {
	l.log("FATAL", msg, keyvals...)
	os.Exit(1)
}

func (l *Logger) log(level, msg string, keyvals ...interface{}) {
	fmt.Println(l.separator())
	fmt.Printf("[%s] [%s]\n", l.timestamp(), level)
	fmt.Printf("  Message: %s\n", msg)

	// Add key-value pairs on separate lines
	for i := 0; i < len(keyvals)-1; i += 2 {
		key := keyvals[i]
		value := keyvals[i+1]
		fmt.Printf("  %v: %v\n", key, value)
	}
}

// Request logs an incoming HTTP request with full details
func (l *Logger) Request(client string, req httpReq, rawSize int) {
	fmt.Println(l.separator())
	fmt.Printf("[%s] --> INCOMING REQUEST\n", l.timestamp())
	fmt.Println(l.subSeparator())

	// Basic info
	fmt.Printf("  Client:       %s\n", client)
	fmt.Printf("  Request:      %s %s %s\n", req.method, req.path, req.version)
	fmt.Printf("  Size:         %d bytes\n", rawSize)

	// Headers
	fmt.Println(l.subSeparator())
	fmt.Println("  HEADERS:")
	if len(req.headers) == 0 {
		fmt.Println("    (none)")
	} else {
		for key, value := range req.headers {
			fmt.Printf("    %s: %s\n", key, value)
		}
	}

	// Body
	fmt.Println(l.subSeparator())
	fmt.Println("  BODY:")
	if req.body == "" {
		fmt.Println("    (empty)")
	} else {
		bodyPreview := req.body
		if len(bodyPreview) > 200 {
			bodyPreview = bodyPreview[:200] + "... (truncated)"
		}
		fmt.Printf("    %s\n", bodyPreview)
	}
}

// Response logs an outgoing HTTP response with full details
func (l *Logger) Response(client string, res httpRes) {
	fmt.Println(l.separator())
	fmt.Printf("[%s] <-- OUTGOING RESPONSE\n", l.timestamp())
	fmt.Println(l.subSeparator())

	// Basic info
	fmt.Printf("  Client:      %s\n", client)
	fmt.Printf("  Status:      %s\n", res.status)
	fmt.Printf("  HTTP Version: %s\n", res.version)

	// Headers
	fmt.Println(l.subSeparator())
	fmt.Println("  HEADERS:")
	if len(res.headers) == 0 {
		fmt.Println("    (none)")
	} else {
		for key, value := range res.headers {
			fmt.Printf("    %s: %s\n", key, value)
		}
	}

	// Body
	fmt.Println(l.subSeparator())
	fmt.Println("  BODY:")
	if res.body == "" {
		fmt.Println("    (empty)")
	} else {
		bodyPreview := res.body
		if len(bodyPreview) > 200 {
			bodyPreview = bodyPreview[:200] + "... (truncated)"
		}
		fmt.Printf("    %s\n", bodyPreview)
		fmt.Printf("    Size: %d bytes\n", len(res.body))
	}
}
