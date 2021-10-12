package store

import (
	"io/ioutil"
	"os"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Unmarshal file into record
func (s Store) parseRecordFromFile(record proto.Message, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if err := protojson.Unmarshal(b, record); err != nil {
		return err
	}
	return nil
}

// Marshal record into file
func (s Store) writeRecordToFile(record proto.Message, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if v, err := protojson.Marshal(record); err != nil {
		file.Write(v)
		return err
	}

	return nil
}
