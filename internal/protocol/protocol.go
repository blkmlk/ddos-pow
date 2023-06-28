package protocol

import (
	"bytes"
	"errors"
	"io"
	"net"
	"time"
)

type instance struct {
	conn net.Conn
}

var (
	ErrMaxLenExceeded = errors.New("max length exceeded")
)

func New(conn net.Conn) *instance {
	return &instance{
		conn: conn,
	}
}

func (i *instance) Read(maxLen int) ([]byte, error) {
	const chunkSize = 1024
	buff := make([]byte, chunkSize)
	var result bytes.Buffer

	read := 0
	for {
		i.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
		n, err := i.conn.Read(buff)
		if errors.Is(err, io.EOF) || n == 0 {
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

func (i *instance) Write(data []byte) error {
	_, err := i.conn.Write(data)
	return err
}
