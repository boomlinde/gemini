package client

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Header encodes a Gemini header
type Header struct {
	Code int
	Meta string
}

// GetHeader will read a Gemini header from the input Reader
func GetHeader(r io.Reader) (*Header, error) {
	headerbytes := make([]byte, 0, 2048)
	buf := make([]byte, 1)
	for {
		_, err := r.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("read header failed: %w", err)
		}
		headerbytes = append(headerbytes, buf[0])
		if headerbytes[len(headerbytes)-1] == '\n' {
			break
		}
		if len(headerbytes) == 2048 {
			return nil, errors.New("too long header")
		}
	}

	fields := strings.SplitN(strings.TrimSpace(string(headerbytes)), " ", 2)
	if len(fields) == 1 {
		fields = append(fields, "")
	}
	if len(fields) != 2 {
		return nil, errors.New("wrong header format")
	}

	code, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("code is not an integer: %w", err)
	}

	// Empty mime type should default to "text/gemini; charset=utf-8"
	if code == 20 && fields[1] == "" {
		fields[1] = "text/gemini; charset=utf-8"
	}

	return &Header{
		Code: code,
		Meta: fields[1],
	}, nil
}
