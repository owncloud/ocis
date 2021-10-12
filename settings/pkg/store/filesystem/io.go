package store

import (
	"os"

	"github.com/gogo/protobuf/jsonpb"
	"google.golang.org/protobuf/proto"
)

// Unmarshal file into record
func (s Store) parseRecordFromFile(record proto.Message, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := jsonpb.Unmarshaler{}
	if err = decoder.Unmarshal(file, record); err != nil {
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

	encoder := jsonpb.Marshaler{}
	if err = encoder.Marshal(file, record); err != nil {
		return err
	}

	return nil
}
