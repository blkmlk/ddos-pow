package stream

import (
	"github.com/stretchr/testify/require"
	"net"
	"sync"
	"testing"
	"time"
)

func TestStreamErrors(t *testing.T) {
	server, client := net.Pipe()
	streamServer := New(server)
	streamClient := New(client)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		buff := make([]byte, 100)
		err := streamServer.Write(buff, time.Millisecond*100)
		require.ErrorIs(t, err, ErrExpired)
	}()

	_, err := streamClient.Read(50, time.Millisecond*100)
	require.ErrorIs(t, err, ErrMaxLenExceeded)
	wg.Wait()
}
