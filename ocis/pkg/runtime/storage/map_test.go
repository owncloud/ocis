package storage

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/owncloud/ocis/ocis/pkg/runtime/process"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := loadStore(); err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}

var (
	store = NewMapStorage()
)

func loadStore() error {
	for i := 0; i < 20; i++ {
		if err := store.Store(process.ProcEntry{
			Pid:       rand.Int(), //nolint:gosec
			Extension: fmt.Sprintf("extension-%s", strconv.Itoa(i)),
		}); err != nil {
			return err
		}
	}

	return nil
}

func TestLoadAll(t *testing.T) {
	all := store.LoadAll()
	assert.NotNil(t, all["extension-1"])
}

func TestDelete(t *testing.T) {
	err := store.Delete(process.ProcEntry{
		Extension: "extension-1",
	})
	assert.Nil(t, err)
	all := store.LoadAll()
	assert.Zero(t, all["extension-1"])
}
