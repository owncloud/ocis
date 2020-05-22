package store

import (
	"os"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// Unmarshal file into record
func (s Store) parseRecordFromFile(record proto.Message, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		s.Logger.Err(err).Msgf("error reading file %v: file not found", filePath)
		return err
	}
	defer file.Close()

	decoder := jsonpb.Unmarshaler{}
	if err = decoder.Unmarshal(file, record); err != nil {
		s.Logger.Err(err).Msgf("error reading file %v: unmarshalling failed", filePath)
		return err
	}
	return nil
}

// Marshal record into file
func (s Store) writeRecordToFile(record proto.Message, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		s.Logger.Err(err).Msgf("error writing file %v: opening failed", filePath)
		return err
	}
	defer file.Close()

	encoder := jsonpb.Marshaler{}
	if err = encoder.Marshal(file, record); err != nil {
		s.Logger.Err(err).Msgf("error writing file %v: marshalling failed", filePath)
		return err
	}

	return nil
}
