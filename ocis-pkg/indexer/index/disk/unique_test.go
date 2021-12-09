package disk

import (
	"os"
	"path"
	"testing"

	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/indexer/errors"
	"github.com/owncloud/ocis/ocis-pkg/indexer/index"
	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
	. "github.com/owncloud/ocis/ocis-pkg/indexer/test"
	"github.com/stretchr/testify/assert"
)

func TestUniqueLookupSingleEntry(t *testing.T) {
	uniq, dataDir := getUniqueIdxSut(t, "Email", User{})
	filesDir := path.Join(dataDir, "users")

	t.Log("existing lookup")
	resultPath, err := uniq.Lookup("mikey@example.com")
	assert.NoError(t, err)

	assert.Equal(t, []string{path.Join(filesDir, "abcdefg-123")}, resultPath)

	t.Log("non-existing lookup")
	resultPath, err = uniq.Lookup("doesnotExists@example.com")
	assert.Error(t, err)
	assert.IsType(t, &errors.NotFoundErr{}, err)
	assert.Empty(t, resultPath)

	_ = os.RemoveAll(dataDir)

}

func TestUniqueUniqueConstraint(t *testing.T) {
	uniq, dataDir := getUniqueIdxSut(t, "Email", User{})

	_, err := uniq.Add("abcdefg-123", "mikey@example.com")
	assert.Error(t, err)
	assert.IsType(t, &errors.AlreadyExistsErr{}, err)

	_ = os.RemoveAll(dataDir)
}

func TestUniqueRemove(t *testing.T) {
	uniq, dataDir := getUniqueIdxSut(t, "Email", User{})

	err := uniq.Remove("", "mikey@example.com")
	assert.NoError(t, err)

	_, err = uniq.Lookup("mikey@example.com")
	assert.Error(t, err)
	assert.IsType(t, &errors.NotFoundErr{}, err)

	_ = os.RemoveAll(dataDir)
}

func TestUniqueUpdate(t *testing.T) {
	uniq, dataDir := getUniqueIdxSut(t, "Email", User{})

	t.Log("successful update")
	err := uniq.Update("", "mikey@example.com", "mikey2@example.com")
	assert.NoError(t, err)

	t.Log("failed update because already exists")
	err = uniq.Update("", "mikey2@example.com", "mikey2@example.com")
	assert.Error(t, err)
	assert.IsType(t, &errors.AlreadyExistsErr{}, err)

	t.Log("failed update because not found")
	err = uniq.Update("", "nonexisting@example.com", "something2@example.com")
	assert.Error(t, err)
	assert.IsType(t, &errors.NotFoundErr{}, err)

	_ = os.RemoveAll(dataDir)
}

func TestUniqueIndexSearch(t *testing.T) {
	sut, dataDir := getUniqueIdxSut(t, "Email", User{})

	res, err := sut.Search("j*@example.com")

	assert.NoError(t, err)
	assert.Len(t, res, 2)

	assert.Equal(t, "ewf4ofk-555", path.Base(res[0]))
	assert.Equal(t, "rulan54-777", path.Base(res[1]))

	_, err = sut.Search("does-not-exist@example.com")
	assert.Error(t, err)
	assert.IsType(t, &errors.NotFoundErr{}, err)

	_ = os.RemoveAll(dataDir)
}

func TestErrors(t *testing.T) {
	assert.True(t, errors.IsAlreadyExistsErr(&errors.AlreadyExistsErr{}))
	assert.True(t, errors.IsNotFoundErr(&errors.NotFoundErr{}))
}

func getUniqueIdxSut(t *testing.T, indexBy string, entityType interface{}) (index.Index, string) {
	dataPath, _ := WriteIndexTestData(Data, "ID", "")
	cfg := config.Config{
		Repo: config.Repo{
			Backend: "disk",
			Disk: config.Disk{
				Path: dataPath,
			},
		},
	}

	sut := NewUniqueIndexWithOptions(
		option.WithTypeName(GetTypeFQN(entityType)),
		option.WithIndexBy(indexBy),
		option.WithFilesDir(path.Join(cfg.Repo.Disk.Path, "users")),
		option.WithDataDir(cfg.Repo.Disk.Path),
	)
	err := sut.Init()
	if err != nil {
		t.Fatal(err)
	}

	for _, u := range Data["users"] {
		pkVal := ValueOf(u, "ID")
		idxByVal := ValueOf(u, "Email")
		_, err := sut.Add(pkVal, idxByVal)
		if err != nil {
			t.Fatal(err)
		}
	}

	return sut, dataPath
}
