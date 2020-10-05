package indexer

import (
	"github.com/owncloud/ocis/accounts/pkg/indexer/errors"
	. "github.com/owncloud/ocis/accounts/pkg/indexer/test"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIndexer_AddWithUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, TestData, "Id")
	indexer := NewIndex(&Config{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	indexer.AddUniqueIndex(&User{}, "UserName", "Id", "users")

	u := &User{Id: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	err := indexer.Add(u)
	assert.NoError(t, err)

}

func TestIndexer_FindByWithUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, TestData, "Id")
	indexer := NewIndex(&Config{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	indexer.AddUniqueIndex(&User{}, "UserName", "Id", "users")

	u := &User{Id: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	err := indexer.Add(u)
	assert.NoError(t, err)

	res, err := indexer.FindBy(User{}, "UserName", "mikey")
	assert.NoError(t, err)
	t.Log(res)
}

func TestIndexer_AddWithNonUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, TestData, "Id")
	indexer := NewIndex(&Config{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	indexer.AddNonUniqueIndex(&TestPet{}, "Kind", "Id", "pets")

	pet1 := TestPet{Id: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := TestPet{Id: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	err := indexer.Add(pet1)
	assert.NoError(t, err)

	err = indexer.Add(pet2)
	assert.NoError(t, err)

	res, err := indexer.FindBy(TestPet{}, "Kind", "Hog")
	assert.NoError(t, err)

	t.Log(res)
}

func TestIndexer_DeleteWithNonUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, TestData, "Id")
	indexer := NewIndex(&Config{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	indexer.AddNonUniqueIndex(&TestPet{}, "Kind", "Id", "pets")

	pet1 := TestPet{Id: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := TestPet{Id: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	err := indexer.Add(pet1)
	assert.NoError(t, err)

	err = indexer.Add(pet2)
	assert.NoError(t, err)

	err = indexer.Delete(pet2)
	assert.NoError(t, err)
}

func TestIndexer_SearchWithNonUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, TestData, "Id")
	indexer := NewIndex(&Config{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	indexer.AddNonUniqueIndex(&TestPet{}, "Name", "Id", "pets")

	pet1 := TestPet{Id: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := TestPet{Id: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	err := indexer.Add(pet1)
	assert.NoError(t, err)

	err = indexer.Add(pet2)
	assert.NoError(t, err)

	res, err := indexer.FindByPartial(pet2, "Name", "*ky")
	assert.NoError(t, err)

	t.Log(res)
}

func TestIndexer_UpdateWithUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, TestData, "Id")
	indexer := NewIndex(&Config{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	indexer.AddUniqueIndex(&User{}, "UserName", "Id", "users")

	user1 := &User{Id: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"}
	user2 := &User{Id: "hijklmn-456", UserName: "frank", Email: "frank@example.com"}

	err := indexer.Add(user1)
	assert.NoError(t, err)

	err = indexer.Add(user2)
	assert.NoError(t, err)

	// Update to non existing value
	err = indexer.Update(user2, "UserName", "frank", "jane")
	assert.NoError(t, err)

	// Update to non existing value
	err = indexer.Update(user2, "UserName", "mikey", "jane")
	assert.Error(t, err)
	assert.IsType(t, &errors.AlreadyExistsErr{}, err)
}

func TestIndexer_UpdateWithNonUniqueIndex(t *testing.T) {
	dataDir := WriteIndexTestData(t, TestData, "Id")
	indexer := NewIndex(&Config{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	indexer.AddNonUniqueIndex(&TestPet{}, "Name", "Id", "pets")

	pet1 := TestPet{Id: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"}
	pet2 := TestPet{Id: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"}

	err := indexer.Add(pet1)
	assert.NoError(t, err)

	err = indexer.Add(pet2)
	assert.NoError(t, err)

	err = indexer.Update(pet2, "Name", "Ricky", "Jonny")
	assert.NoError(t, err)
}

/*
func TestManagerQueryMultipleIndices(t *testing.T) {
	dataDir := writeIndexTestData(t, testData, "Id")
	man := NewIndex(&Config{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	err := man.AddUniqueIndex("User", "Email", "users")
	assert.NoError(t, err)

	err = man.AddUniqueIndex("User", "UserName", "users")
	assert.NoError(t, err)

	err = man.AddNormalIndex("TestPet", "Color", "pets")
	assert.NoError(t, err)

	err = man.AddUniqueIndex("TestPet", "Name", "pets")
	assert.NoError(t, err)

	for path := range testData {
		for _, entity := range testData[path] {
			err := man.Add(valueOf(entity, "Id"), entity)
			assert.NoError(t, err)
		}
	}

	type test struct {
		typeName, key, value, wantRes string
		wantErr                       error
	}

	tests := []test{
		{typeName: "User", key: "Email", value: "jacky@example.com", wantRes: "ewf4ofk-555"},
		{typeName: "User", key: "UserName", value: "jacky", wantRes: "ewf4ofk-555"},
		{typeName: "TestPet", key: "Color", value: "Brown", wantRes: "rebef-123"},
		{typeName: "TestPet", key: "Color", value: "Cyan", wantRes: "", wantErr: &notFoundErr{}},
	}

	for _, tc := range tests {
		name := fmt.Sprintf("Query%sBy%s=%s", tc.typeName, tc.key, tc.value)
		t.Run(name, func(t *testing.T) {
			pk, err := man.Find(tc.typeName, tc.key, tc.value)
			assert.Equal(t, tc.wantRes, pk)
			assert.IsType(t, tc.wantErr, err)
		})
	}

	_ = os.RemoveAll(dataDir)
}

*/

/*
func TestManagerDelete(t *testing.T) {
	dataDir := writeIndexTestData(t, testData, "Id")
	man := NewIndex(&Config{
		DataDir:          dataDir,
		IndexRootDirName: "index.disk",
		Log:              zerolog.Logger{},
	})

	err := man.AddUniqueIndex("User", "Email", "users")
	assert.NoError(t, err)

	err = man.AddUniqueIndex("User", "UserName", "users")
	assert.NoError(t, err)

	err = man.AddUniqueIndex("TestPet", "Name", "pets")
	assert.NoError(t, err)

	for path := range testData {
		for _, entity := range testData[path] {
			err := man.Add(valueOf(entity, "Id"), entity)
			assert.NoError(t, err)
		}
	}

	err = man.Delete("User", "hijklmn-456")
	_ = os.RemoveAll(dataDir)

}

*/
