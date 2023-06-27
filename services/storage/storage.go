package storage

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage interface {
	CreateChallenge(ctx context.Context, challenge *Challenge) error
	GetChallenge(ctx context.Context, id string) (*Challenge, error)
}

func NewMapStorage() Storage {
	return &mapStorage{
		challenges: make(map[string]Challenge),
	}
}

type mapStorage struct {
	locker     sync.RWMutex
	challenges map[string]Challenge
}

func (m *mapStorage) CreateChallenge(ctx context.Context, challenge *Challenge) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.challenges[challenge.ID] = *challenge
	return nil
}

func (m *mapStorage) GetChallenge(ctx context.Context, id string) (*Challenge, error) {
	challenge, ok := m.challenges[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &challenge, nil
}
