package stream

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"
)

type stream struct {
	conn net.Conn
}

var (
	ErrInvalidFormat  = errors.New("invalid format")
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
	if err := s.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer func() {
		_ = s.conn.SetReadDeadline(time.Time{})
	}()

	sizeBuff := make([]byte, 2)
	n, err := s.conn.Read(sizeBuff)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return nil, ErrExpired
		}
		if errors.Is(err, io.EOF) {
			return nil, ErrClosed
		}
		return nil, err
	}
	if n != 2 {
		return nil, ErrInvalidFormat
	}

	size := binary.LittleEndian.Uint16(sizeBuff)

	if maxLen > 0 && int(size) > maxLen {
		return nil, ErrMaxLenExceeded
	}

	buff := make([]byte, size)
	var result bytes.Buffer

	n, err = s.conn.Read(buff)
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return nil, ErrExpired
		}
		if errors.Is(err, io.EOF) {
			return nil, ErrClosed
		}
		return nil, err
	}
	if n != int(size) {
		return nil, ErrInvalidFormat
	}

	result.Write(buff[:n])

	return result.Bytes(), nil
}

func (s *stream) Write(data []byte, timeout time.Duration) error {
	if err := s.conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	_, err := s.conn.Write(data)
	return err
}
