package storage_test

import (
	"context"
	"github.com/blkmlk/ddos-pow/services/storage"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestAll(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite
	storage storage.Storage
}

func (t *testSuite) SetupTest() {
	t.storage = storage.NewMapStorage()
}

func (t *testSuite) TestCreateAndGet() {
	ctx := context.Background()

	challenge := storage.NewChallenge("puzzle", 8, 2, 1, 16)
	err := t.storage.CreateChallenge(ctx, &challenge)
	t.Require().NoError(err)

	foundChallenge, err := t.storage.GetChallenge(ctx, challenge.ID)
	t.Require().NoError(err)
	t.Require().NotNil(foundChallenge)
	t.Require().Equal(challenge, *foundChallenge)

	foundChallenge, err = t.storage.GetChallenge(ctx, "12345")
	t.Require().ErrorIs(err, storage.ErrNotFound)
	t.Require().Nil(foundChallenge)
}
