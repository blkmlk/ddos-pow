package stream

import (
	"bytes"
	"errors"
	"io"
	"net"
	"time"
)

type stream struct {
	conn net.Conn
}

var (
	ErrClosed         = errors.New("closed")
	ErrExpired        = errors.New("expired")
	ErrMaxLenExceeded = errors.New("max length exceeded")
)

func New(conn net.Conn) *stream {
	return &stream{
		conn: conn,
	}
}

func (s *stream) Read(maxLen int, timeout time.Duration) ([]byte, error) {
	const chunkSize = 256
	buff := make([]byte, chunkSize)
	var result bytes.Buffer

	read := 0
	for {
		if err := s.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			return nil, err
		}

		n, err := s.conn.Read(buff)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				break
			}
			if errors.Is(err, io.EOF) {
				return nil, ErrClosed
			}
			return nil, err
		}

		read += n
		if maxLen > 0 && read > maxLen {
			return nil, ErrMaxLenExceeded
		}
		result.Write(buff[:n])

		if n < chunkSize {
			break
		}
	}
	return result.Bytes(), nil
}

func (s *stream) ReadUntil(maxLen int, timeout time.Duration) ([]byte, error) {
	startedAt := time.Now()
	for {
		data, err := s.Read(maxLen, time.Millisecond*50)
		if err != nil {
			return nil, err
		}

		if len(data) > 0 {
			return data, nil
		}

		if time.Since(startedAt) > timeout {
			return nil, ErrExpired
		}
	}
}

func (s *stream) Write(data []byte, timeout time.Duration) error {
	if err := s.conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	_, err := s.conn.Write(data)
	return err
}
