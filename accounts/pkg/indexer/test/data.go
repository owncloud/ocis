package test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type TestUser struct {
	Id, UserName, Email string
}

type TestPet struct {
	Id, Kind, Color, Name string
}

var TestData = map[string][]interface{}{
	"users": {
		TestUser{Id: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"},
		TestUser{Id: "hijklmn-456", UserName: "frank", Email: "frank@example.com"},
		TestUser{Id: "ewf4ofk-555", UserName: "jacky", Email: "jacky@example.com"},
		TestUser{Id: "rulan54-777", UserName: "jones", Email: "jones@example.com"},
	},
	"pets": {
		TestPet{Id: "rebef-123", Kind: "Dog", Color: "Brown", Name: "Waldo"},
		TestPet{Id: "wefwe-456", Kind: "Cat", Color: "White", Name: "Snowy"},
		TestPet{Id: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"},
		TestPet{Id: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"},
	},
}

func WriteIndexTestData(t *testing.T, m map[string][]interface{}, pk string) string {
	rootDir := CreateTmpDir(t)
	for dirName := range m {
		fileTypePath := path.Join(rootDir, dirName)

		if err := os.MkdirAll(fileTypePath, 0777); err != nil {
			t.Fatal(err)
		}
		for _, u := range m[dirName] {
			data, err := json.Marshal(u)
			if err != nil {
				t.Fatal(err)
			}

			pkVal := ValueOf(u, pk)
			if err := ioutil.WriteFile(path.Join(fileTypePath, pkVal), data, 0777); err != nil {
				t.Fatal(err)
			}
		}
	}

	return rootDir
}

func WriteIndexTestDataCS3(t *testing.T, m map[string][]interface{}, pk string) string {
	rootDir := "/var/tmp/ocis/storage/users/data"
	for dirName := range m {
		fileTypePath := path.Join(rootDir, dirName)

		if err := os.MkdirAll(fileTypePath, 0777); err != nil {
			t.Fatal(err)
		}
		for _, u := range m[dirName] {
			data, err := json.Marshal(u)
			if err != nil {
				t.Fatal(err)
			}

			pkVal := ValueOf(u, pk)
			if err := ioutil.WriteFile(path.Join(fileTypePath, pkVal), data, 0777); err != nil {
				t.Fatal(err)
			}
		}
	}

	return rootDir
}
