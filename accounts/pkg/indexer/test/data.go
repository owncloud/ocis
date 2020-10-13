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
	UID int
}

// Pet is a pet.
type Pet struct {
	ID, Kind, Color, Name string
	UID int
}

// Data mock data.
var Data = map[string][]interface{}{
	"users": {
		User{ID: "abcdefg-123", UserName: "mikey", Email: "mikey@example.com", UID: -1},
		User{ID: "hijklmn-456", UserName: "frank", Email: "frank@example.com", UID: -1},
		User{ID: "ewf4ofk-555", UserName: "jacky", Email: "jacky@example.com", UID: -1},
		User{ID: "rulan54-777", UserName: "jones", Email: "jones@example.com", UID: -1},
	},
	"pets": {
		Pet{ID: "rebef-123", Kind: "Dog", Color: "Brown", Name: "Waldo", UID: -1},
		Pet{ID: "wefwe-456", Kind: "Cat", Color: "White", Name: "Snowy", UID: -1},
		Pet{ID: "goefe-789", Kind: "Hog", Color: "Green", Name: "Dicky", UID: -1},
		Pet{ID: "xadaf-189", Kind: "Hog", Color: "Green", Name: "Ricky", UID: -1},
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

// WriteIndexBenchmarkDataCS3 writes more data to disk.
func WriteIndexBenchmarkDataCS3(b *testing.B, m map[string][]interface{}, pk string) string {
	rootDir := "/var/tmp/ocis/storage/users/data"
	for dirName := range m {
		fileTypePath := path.Join(rootDir, dirName)

		if err := os.MkdirAll(fileTypePath, 0777); err != nil {
			b.Fatal(err)
		}
		for _, u := range m[dirName] {
			data, err := json.Marshal(u)
			if err != nil {
				b.Fatal(err)
			}

			pkVal := ValueOf(u, pk)
			if err := ioutil.WriteFile(path.Join(fileTypePath, pkVal), data, 0777); err != nil {
				b.Fatal(err)
			}
		}
	}

	return rootDir
}
