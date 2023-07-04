package stream

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
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

	// get the size of the message first
	sizeBuff := make([]byte, 2)
	n, err := s.conn.Read(sizeBuff)
	if err != nil {
		return nil, handleError(err)
	}
	if n != 2 {
		return nil, ErrInvalidFormat
	}

	size := binary.LittleEndian.Uint16(sizeBuff)

	if maxLen > 0 && int(size) > maxLen {
		return nil, ErrMaxLenExceeded
	}

	// waiting for the message
	buff := make([]byte, size)
	var result bytes.Buffer

	n, err = s.conn.Read(buff)
	if err != nil {
		return nil, handleError(err)
	}
	if n != int(size) {
		return nil, ErrInvalidFormat
	}

	result.Write(buff[:n])

	return result.Bytes(), nil
}

func (s *stream) Write(data []byte, timeout time.Duration) error {
	if len(data) > math.MaxUint16 {
		return ErrInvalidFormat
	}
	if err := s.conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer func() {
		_ = s.conn.SetReadDeadline(time.Time{})
	}()

	sizeBuff := make([]byte, 2)
	binary.LittleEndian.PutUint16(sizeBuff, uint16(len(data)))
	if _, err := s.conn.Write(sizeBuff); err != nil {
		return handleError(err)
	}

	if _, err := s.conn.Write(data); err != nil {
		return handleError(err)
	}
	return nil
}

func handleError(err error) error {
	if e, ok := err.(net.Error); ok && e.Timeout() {
		return ErrExpired
	}
	if errors.Is(err, io.EOF) {
		return ErrClosed
	}
	return err
}
