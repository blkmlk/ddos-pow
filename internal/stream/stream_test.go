package stream

import (
	"crypto/rand"
	"github.com/stretchr/testify/require"
	"net"
	"sync"
	"testing"
	"time"
)

func TestStream(t *testing.T) {
	server, client := net.Pipe()
	streamServer := New(server)
	streamClient := New(client)

	toSend := make([]byte, 100)
	_, err := rand.Read(toSend)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := streamServer.Write(toSend, time.Millisecond*100)
		require.NoError(t, err)
	}()

	received, err := streamClient.Read(100, time.Millisecond*100)
	require.NoError(t, err)
	require.Equal(t, toSend, received)

	wg.Wait()
}

func TestStreamClosed(t *testing.T) {
	server, client := net.Pipe()
	streamServer := New(server)
	streamClient := New(client)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := streamServer.conn.Close()
		require.NoError(t, err)
	}()

	_, err := streamClient.Read(100, time.Millisecond*100)
	require.ErrorIs(t, err, ErrClosed)

	wg.Wait()
}

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
