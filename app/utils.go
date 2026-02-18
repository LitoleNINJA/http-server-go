package main

import (
	"compress/gzip"
	"bytes"
)

func HasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func CutPrefix(s, prefix string) (after string, found bool) {
	if !HasPrefix(s, prefix) {
		return s, false
	}
	return s[len(prefix):], true
}

func encodeGZIP(msg string) string {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte(msg))
	gz.Close()

	return string(buf.Bytes())
}