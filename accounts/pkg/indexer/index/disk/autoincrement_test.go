package disk

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/owncloud/ocis/accounts/pkg/indexer/option"
	//. "github.com/owncloud/ocis/accounts/pkg/indexer/test"
	"github.com/stretchr/testify/assert"
)

func TestIsValidKind(t *testing.T) {
	scenarios := []struct {
		panics  bool
		name    string
		indexBy string
		entity  struct {
			Number      int
			Name        string
			NumberFloat float32
		}
	}{
		{
			name:    "valid autoincrement index",
			panics:  false,
			indexBy: "Number",
			entity: struct {
				Number      int
				Name        string
				NumberFloat float32
			}{
				Name: "tesy-mc-testace",
			},
		},
		{
			name:    "create autoincrement index on a non-existing field",
			panics:  true,
			indexBy: "Age",
			entity: struct {
				Number      int
				Name        string
				NumberFloat float32
			}{
				Name: "tesy-mc-testace",
			},
		},
		{
			name:    "attempt to create an autoincrement index with no entity",
			panics:  true,
			indexBy: "Age",
		},
		{
			name:    "create autoincrement index on a non-numeric field",
			panics:  true,
			indexBy: "Name",
			entity: struct {
				Number      int
				Name        string
				NumberFloat float32
			}{
				Name: "tesy-mc-testace",
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			if scenario.panics {
				assert.Panics(t, func() {
					_ = NewAutoincrementIndex(
						option.WithEntity(scenario.entity),
						option.WithIndexBy(scenario.indexBy),
					)
				})
			} else {
				assert.NotPanics(t, func() {
					_ = NewAutoincrementIndex(
						option.WithEntity(scenario.entity),
						option.WithIndexBy(scenario.indexBy),
					)
				})
			}
		})
	}
}

func TestNext(t *testing.T) {
	scenarios := []struct {
		name     string
		expected int
		indexBy  string
		entity   interface{}
	}{
		{
			name:     "get next value",
			expected: 0,
			indexBy:  "Number",
			entity: struct {
				Number      int
				Name        string
				NumberFloat float32
			}{
				Name: "tesy-mc-testace",
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			tmpDir, err := createTmpDirStr()
			assert.NoError(t, err)

			err = os.MkdirAll(filepath.Join(tmpDir, "data"), 0777)
			assert.NoError(t, err)

			i := NewAutoincrementIndex(
				option.WithDataDir(tmpDir),
				option.WithFilesDir(filepath.Join(tmpDir, "data")),
				option.WithEntity(scenario.entity),
				option.WithTypeName("LambdaType"),
				option.WithIndexBy(scenario.indexBy),
			)

			err = i.Init()
			assert.NoError(t, err)

			tmpFile, err := os.Create(filepath.Join(tmpDir, "data", "test-example"))
			assert.NoError(t, err)
			assert.NoError(t, tmpFile.Close())

			oldName, err := i.Add("test-example", "")
			assert.NoError(t, err)
			assert.Equal(t, filepath.Base(oldName), 0)

			oldName, err = i.Add("test-example", "")
			assert.NoError(t, err)
			assert.Equal(t, filepath.Base(oldName), 1)

			oldName, err = i.Add("test-example", "")
			assert.NoError(t, err)
			assert.Equal(t, filepath.Base(oldName), 2)
			t.Log(oldName)

			_ = os.RemoveAll(tmpDir)
		})
	}
}

func BenchmarkAdd(b *testing.B) {
	tmpDir, err := createTmpDirStr()
	assert.NoError(b, err)

	err = os.MkdirAll(filepath.Join(tmpDir, "data"), 0777)
	assert.NoError(b, err)

	tmpFile, err := os.Create(filepath.Join(tmpDir, "data", "test-example"))
	assert.NoError(b, err)
	assert.NoError(b, tmpFile.Close())

	i := NewAutoincrementIndex(
		option.WithDataDir(tmpDir),
		option.WithFilesDir(filepath.Join(tmpDir, "data")),
		option.WithEntity(struct {
			Number      int
			Name        string
			NumberFloat float32
		}{}),
		option.WithTypeName("LambdaType"),
		option.WithIndexBy("Number"),
	)

	err = i.Init()
	assert.NoError(b, err)

	for n := 0; n < b.N; n++ {
		_, err := i.Add("test-example", "")
		if err != nil {
			b.Error(err)
		}
		assert.NoError(b, err)
	}
}

func createTmpDirStr() (string, error) {
	name, err := ioutil.TempDir("/var/tmp", "testfiles-*")
	if err != nil {
		return "", err
	}

	return name, nil
}
