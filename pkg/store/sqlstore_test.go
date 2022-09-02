package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func helperConnect(t *testing.T) (store *SQLStore, cleanup func()) {
	t.Helper()
	assert := require.New(t)

	var err error
	store, err = NewSQLStore()
	assert.NoError(err)

	cleanup = func() {
		err := store.Cleanup()
		assert.NoError(err)
	}
	cleanup()

	return
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
