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

const (
	NetworkDelay = time.Millisecond * 500
)

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
		if err := s.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
			return nil, err
		}

		n, err := s.conn.Read(buff)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				break
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
	}
	return result.Bytes(), nil
}

func (s *stream) ReadUntil(maxLen int, timeout time.Duration) ([]byte, error) {
	startedAt := time.Now()
	for {
		data, err := s.Read(maxLen, NetworkDelay)
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
