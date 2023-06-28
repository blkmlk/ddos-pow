package stream

import (
	"bytes"
	"errors"
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
	buff := make([]byte, 256)
	var result bytes.Buffer

	read := 0
	for {
		_ = s.conn.SetReadDeadline(time.Now().Add(timeout))

		n, err := s.conn.Read(buff)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
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
		buff = buff[:0]
	}
	return result.Bytes(), nil
}

func (s *stream) ReadUntil(maxLen int, maxTimeout time.Duration) ([]byte, error) {
	startedAt := time.Now()
	for {
		data, err := s.Read(maxLen, time.Millisecond*50)
		if err != nil {
			return nil, err
		}

		if len(data) > 0 {
			return data, nil
		}

		if time.Since(startedAt) > maxTimeout {
			return nil, ErrExpired
		}
	}
}

func (s *stream) Write(data []byte) error {
	_, err := s.conn.Write(data)
	return err
}
