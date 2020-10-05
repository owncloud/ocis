package cs3

import (
	. "github.com/owncloud/ocis/accounts/pkg/indexer/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCS3NonUniqueIndex_FakeSymlink(t *testing.T) {
	dataDir := WriteIndexTestDataCS3(t, TestData, "Id")
	sut := NewUniqueIndex("User", "Name", "/meta", "index.cs3", &Config{
		ProviderAddr:    "0.0.0.0:9215",
		DataURL:         "http://localhost:9216",
		DataPrefix:      "data",
		JWTSecret:       "Pive-Fumkiu4",
		ServiceUserName: "",
		ServiceUserUUID: "",
	})

	err := sut.Init()
	assert.NoError(t, err)

	res, err := sut.Add("abcdefg-123", "mikey")
	assert.NoError(t, err)
	t.Log(res)
	//
	//resLookup, err := sut.Lookup("mikey")
	//assert.NoError(t, err)
	//t.Log(resLookup)
	//
	//err = sut.Update("abcdefg-123", "mikey", "mickeyX")
	//assert.NoError(t, err)
	//
	//_, err = sut.Search("mi*")
	//assert.NoError(t, err)
	//
	//err = sut.Remove("", "mikey")
	//assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}

//func TestCS3UniqueIndexSearch(t *testing.T) {
//	dataDir := WriteIndexTestDataCS3(t, TestData, "Id")
//	sut := NewUniqueIndex("User", "UserName", "/meta", "index.cs3", &Config{
//		ProviderAddr:    "0.0.0.0:9215",
//		DataURL:         "http://localhost:9216",
//		DataPrefix:      "data",
//		JWTSecret:       "Pive-Fumkiu4",
//		ServiceUserName: "",
//		ServiceUserUUID: "",
//	})
//
//	err := sut.Init()
//	assert.NoError(t, err)
//
//	_, err = sut.Add("hijklmn-456", "mikey")
//	assert.NoError(t, err)
//
//	_, err = sut.Add("ewf4ofk-555", "jacky")
//	assert.NoError(t, err)
//
//	res, err := sut.Search("*y")
//	assert.NoError(t, err)
//	t.Log(res)
//
//	_ = os.RemoveAll(dataDir)
//
//}
