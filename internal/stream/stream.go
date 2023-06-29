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

	if err := s.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	read := 0
	for {
		n, err := s.conn.Read(buff)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				return nil, ErrExpired
			}
			if errors.Is(err, io.EOF) && read > 0 {
				break
			}
			return nil, err
		}
		if n == 0 {
			break
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

func (s *stream) Write(data []byte, timeout time.Duration) error {
	if err := s.conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	_, err := s.conn.Write(data)
	return err
}
