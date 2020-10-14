package indexer

import (
	"github.com/owncloud/ocis/accounts/pkg/config"
	_ "github.com/owncloud/ocis/accounts/pkg/indexer/index/cs3"
	_ "github.com/owncloud/ocis/accounts/pkg/indexer/index/disk"
	. "github.com/owncloud/ocis/accounts/pkg/indexer/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestIndexer_AddWithUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, Data, "ID")
	indexer := CreateIndexer(&config.Config{
		Repo: config.Repo{
			Disk: config.Disk{
				Path: dataDir,
			},
		},
	})

	err := indexer.AddIndex(&User{}, "UserName", "ID", "users", "unique")
	assert.NoError(t, err)

	u := &User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	_, err = indexer.Add(u)
	assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_AddWithUniqueIndexCS3(t *testing.T) {
	dir := WriteIndexTestDataCS3(t, Data, "ID")
	indexer := CreateIndexer(&config.Config{
		Repo: config.Repo{
			CS3: config.CS3{
				ProviderAddr: "0.0.0.0:9215",
				DataURL:      "http://localhost:9216",
				DataPrefix:   "data",
				JWTSecret:    "Pive-Fumkiu4",
			},
		},
	})

	err := indexer.AddIndex(&User{}, "UserName", "ID", "users", "unique")
	assert.NoError(t, err)

	u := &User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	_, err = indexer.Add(u)
	assert.NoError(t, err)

	_ = os.RemoveAll(dir)
}

func TestIndexer_AddWithNonUniqueIndexCS3(t *testing.T) {
	dataDir := WriteIndexTestDataCS3(t, Data, "ID")
	indexer := CreateIndexer(&config.Config{
		Repo: config.Repo{
			CS3: config.CS3{
				ProviderAddr: "0.0.0.0:9215",
				DataURL:      "http://localhost:9216",
				DataPrefix:   "data",
				JWTSecret:    "Pive-Fumkiu4",
			},
		},
	})

	err := indexer.AddIndex(&User{}, "UserName", "ID", "users", "non_unique")
	assert.NoError(t, err)

	u := &User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	_, err = indexer.Add(u)
	assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_FindByWithUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, Data, "ID")
	indexer := CreateIndexer(&config.Config{
		Repo: config.Repo{
			Disk: config.Disk{
				Path: dataDir,
			},
		},
	})

	err := indexer.AddIndex(&User{}, "UserName", "ID", "users", "unique")
	assert.NoError(t, err)

	u := &User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	_, err = indexer.Add(u)
	assert.NoError(t, err)

	res, err := indexer.FindBy(User{}, "UserName", "mikey")
	assert.NoError(t, err)
	t.Log(res)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_AddWithNonUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, Data, "ID")
	indexer := CreateIndexer(&config.Config{
		Repo: config.Repo{
			Disk: config.Disk{
				Path: dataDir,
			},
		},
	})

	err := indexer.AddIndex(&Pet{}, "Kind", "ID", "pets", "non_unique")
	assert.NoError(t, err)

	pet1 := Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	_, err = indexer.Add(pet1)
	assert.NoError(t, err)

	_, err = indexer.Add(pet2)
	assert.NoError(t, err)

	res, err := indexer.FindBy(Pet{}, "Kind", "Hog")
	assert.NoError(t, err)

	t.Log(res)
}

func TestIndexer_DeleteWithNonUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, Data, "ID")
	indexer := CreateIndexer(&config.Config{
		Repo: config.Repo{
			Disk: config.Disk{
				Path: dataDir,
			},
		},
	})

	err := indexer.AddIndex(&Pet{}, "Kind", "ID", "pets", "non_unique")
	assert.NoError(t, err)

	pet1 := Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	_, err = indexer.Add(pet1)
	assert.NoError(t, err)

	_, err = indexer.Add(pet2)
	assert.NoError(t, err)

	err = indexer.Delete(pet2)
	assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_SearchWithNonUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, Data, "ID")
	indexer := CreateIndexer(&config.Config{
		Repo: config.Repo{
			Disk: config.Disk{
				Path: dataDir,
			},
		},
	})

	err := indexer.AddIndex(&Pet{}, "Name", "ID", "pets", "non_unique")
	assert.NoError(t, err)

	pet1 := Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	_, err = indexer.Add(pet1)
	assert.NoError(t, err)

	_, err = indexer.Add(pet2)
	assert.NoError(t, err)

	res, err := indexer.FindByPartial(pet2, "Name", "*ky")
	assert.NoError(t, err)

	t.Log(res)
	_ = os.RemoveAll(dataDir)
}

func TestIndexer_UpdateWithUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, Data, "ID")
	indexer := CreateIndexer(&config.Config{
		Repo: config.Repo{
			Disk: config.Disk{
				Path: dataDir,
			},
		},
	})

	err := indexer.AddIndex(&User{}, "UserName", "ID", "users", "unique")
	assert.NoError(t, err)

	err = indexer.AddIndex(&User{}, "Email", "ID", "users", "unique")
	assert.NoError(t, err)

	user1 := &User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	user2 := &User{ID: "hijklmn-456", UserName: "frank", Email: "frank@example.com"}

	_, err = indexer.Add(user1)
	assert.NoError(t, err)

	_, err = indexer.Add(user2)
	assert.NoError(t, err)

	err = indexer.Update(user1, &User{
		ID:       "abcdefg-123",
		UserName: "mikey-new",
		Email:    "mikey@example.com",
	})
	assert.NoError(t, err)
	v, err1 := indexer.FindBy(&User{}, "UserName", "mikey-new")
	assert.NoError(t, err1)
	assert.Len(t, v, 1)
	v, err2 := indexer.FindBy(&User{}, "UserName", "mikey")
	assert.NoError(t, err2)
	assert.Len(t, v, 0)

	err1 = indexer.Update(&User{
		ID:       "abcdefg-123",
		UserName: "mikey-new",
		Email:    "mikey@example.com",
	}, &User{
		ID:       "abcdefg-123",
		UserName: "mikey-newest",
		Email:    "mikey-new@example.com",
	})
	assert.NoError(t, err1)
	fbUserName, err2 := indexer.FindBy(&User{}, "UserName", "mikey-newest")
	assert.NoError(t, err2)
	assert.Len(t, fbUserName, 1)
	fbEmail, err3 := indexer.FindBy(&User{}, "Email", "mikey-new@example.com")
	assert.NoError(t, err3)
	assert.Len(t, fbEmail, 1)

	_ = os.RemoveAll(dataDir)
}

func TestIndexer_UpdateWithNonUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, Data, "ID")
	indexer := CreateIndexer(&config.Config{
		Repo: config.Repo{
			Disk: config.Disk{
				Path: dataDir,
			},
		},
	})

	err := indexer.AddIndex(&Pet{}, "Name", "ID", "pets", "non_unique")
	assert.NoError(t, err)

	pet1 := Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	_, err = indexer.Add(pet1)
	assert.NoError(t, err)

	_, err = indexer.Add(pet2)
	assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}
