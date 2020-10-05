package cs3

import (
	. "github.com/owncloud/ocis/accounts/pkg/indexer/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
)

func TestCS3NonUniqueIndex_FakeSymlink(t *testing.T) {
	dataDir := WriteIndexTestDataCS3(t, TestData, "Id")
	sut := NewNonUniqueIndex("User", "Name", "/meta", "index.cs3", &Config{
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

	_, err = sut.Add("abcdefg-234", "mikey")
	assert.NoError(t, err)
	t.Log(res)

	_, err = sut.Add("soasdioahsfash", "milo")
	assert.NoError(t, err)
	t.Log(res)

	_, err = sut.Add("asdasdsa", "jonas")
	assert.NoError(t, err)
	t.Log(res)

	lookupRes, err := sut.Lookup("mikey")
	assert.NoError(t, err)
	t.Log(lookupRes)

	searchRes, err := sut.Search("mi*")
	assert.NoError(t, err)
	t.Log(searchRes)

	err = sut.Update("abcdefg-234", "mikey", "jonas")
	assert.NoError(t, err)

	_ = os.RemoveAll(dataDir)
}
