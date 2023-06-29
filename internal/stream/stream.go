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

func (s *stream) Write(data []byte, timeout time.Duration) error {
	if err := s.conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	_, err := s.conn.Write(data)
	return err
}
