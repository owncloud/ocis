package test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

// User is a user.
type User struct {
	ID, UserName, Email string
}

// Pet is a pet.
type Pet struct {
	ID, Kind, Color, Name string
}

// Data mock data.
var Data = map[string][]interface{}{
	"users": {
		User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com"},
		User{ID: "hijklmn-456", UserName: "frank", Email: "frank@example.com"},
		User{ID: "ewf4ofk-555", UserName: "jacky", Email: "jacky@example.com"},
		User{ID: "rulan54-777", UserName: "jones", Email: "jones@example.com"},
	},
	"pets": {
		Pet{ID: "rebef-123", Kind: "Dog", Color: "Brown", Name: "Waldo"},
		Pet{ID: "wefwe-456", Kind: "Cat", Color: "White", Name: "Snowy"},
		Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky"},
		Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky"},
	},
}

// WriteIndexTestData writes mock data to disk.
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

// WriteIndexTestDataCS3 writes more data to disk.
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
