package disk

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/owncloud/ocis/ocis-pkg/indexer/config"
	"github.com/owncloud/ocis/ocis-pkg/indexer/errors"
	"github.com/owncloud/ocis/ocis-pkg/indexer/index"
	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
	. "github.com/owncloud/ocis/ocis-pkg/indexer/test"
	"github.com/stretchr/testify/assert"
)

func TestNonUniqueIndexAdd(t *testing.T) {
	sut, dataPath := getNonUniqueIdxSut(t, Pet{}, "Color")

	ids, err := sut.Lookup("Green")
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"goefe-789", "xadaf-189"}, ids)

	ids, err = sut.Lookup("White")
	assert.NoError(t, err)
	assert.EqualValues(t, []string{"wefwe-456"}, ids)

	ids, err = sut.Lookup("Cyan")
	assert.Error(t, err)
	assert.EqualValues(t, []string{}, ids)

	_ = os.RemoveAll(dataPath)

}

func TestNonUniqueIndexUpdate(t *testing.T) {
	sut, dataPath := getNonUniqueIdxSut(t, Pet{}, "Color")

	err := sut.Update("goefe-789", "Green", "Black")
	assert.NoError(t, err)

	err = sut.Update("xadaf-189", "Green", "Black")
	assert.NoError(t, err)

	assert.DirExists(t, path.Join(dataPath, fmt.Sprintf("index.disk/non_unique.%v.Color/Black", GetTypeFQN(Pet{}))))
	assert.NoDirExists(t, path.Join(dataPath, fmt.Sprintf("index.disk/non_unique.%v.Color/Green", GetTypeFQN(Pet{}))))

	_ = os.RemoveAll(dataPath)
}

func TestNonUniqueIndexDelete(t *testing.T) {
	sut, dataPath := getNonUniqueIdxSut(t, Pet{}, "Color")
	assert.FileExists(t, path.Join(dataPath, fmt.Sprintf("index.disk/non_unique.%v.Color/Green/goefe-789", GetTypeFQN(Pet{}))))

	err := sut.Remove("goefe-789", "Green")
	assert.NoError(t, err)
	assert.NoFileExists(t, path.Join(dataPath, fmt.Sprintf("index.disk/non_unique.%v.Color/Green/goefe-789", GetTypeFQN(Pet{}))))
	assert.FileExists(t, path.Join(dataPath, fmt.Sprintf("index.disk/non_unique.%v.Color/Green/xadaf-189", GetTypeFQN(Pet{}))))

	_ = os.RemoveAll(dataPath)
}

func TestNonUniqueIndexSearch(t *testing.T) {
	sut, dataPath := getNonUniqueIdxSut(t, Pet{}, "Email")

	res, err := sut.Search("Gr*")

	assert.NoError(t, err)
	assert.Len(t, res, 2)

	assert.Equal(t, "goefe-789", path.Base(res[0]))
	assert.Equal(t, "xadaf-189", path.Base(res[1]))

	_, err = sut.Search("does-not-exist@example.com")
	assert.Error(t, err)
	assert.IsType(t, &errors.NotFoundErr{}, err)

	_ = os.RemoveAll(dataPath)
}

// entity: used to get the fully qualified name for the index root path.
func getNonUniqueIdxSut(t *testing.T, entity interface{}, indexBy string) (index.Index, string) {
	dataPath, _ := WriteIndexTestData(Data, "ID", "")
	cfg := config.Config{
		Repo: config.Repo{
			Backend: "disk",
			Disk: config.Disk{
				Path: dataPath,
			},
		},
	}

	sut := NewNonUniqueIndexWithOptions(
		option.WithTypeName(GetTypeFQN(entity)),
		option.WithIndexBy(indexBy),
		option.WithFilesDir(path.Join(cfg.Repo.Disk.Path, "pets")),
		option.WithDataDir(cfg.Repo.Disk.Path),
	)
	err := sut.Init()
	if err != nil {
		t.Fatal(err)
	}

	for _, u := range Data["pets"] {
		pkVal := ValueOf(u, "ID")
		idxByVal := ValueOf(u, "Color")
		_, err := sut.Add(pkVal, idxByVal)
		if err != nil {
			t.Fatal(err)
		}
	}

	return sut, dataPath
}
