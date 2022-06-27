package godata

import (
	"bytes"
	"strconv"
)

// A response is a dictionary of keys to their corresponding fields. This will
// be converted into a JSON dictionary in the response to the web client.
type GoDataResponse struct {
	Fields map[string]*GoDataResponseField
}

// Serialize the result as JSON for sending to the client. If an error
// occurs during the serialization, it will be returned.
func (r *GoDataResponse) Json() ([]byte, error) {
	result, err := prepareJsonDict(r.Fields)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// A response that is a primitive JSON type or a list or a dictionary. When
// writing to JSON, it is automatically mapped from the Go type to a suitable
// JSON data type. Any type can be used, but if the data type is not supported
// for serialization, then an error is thrown.
type GoDataResponseField struct {
	Value interface{}
}

// Convert the response field to a JSON serialized form. If the type is not
// string, []byte, int, float64, map[string]*GoDataResponseField, or
// []*GoDataResponseField, then an error will be thrown.
func (f *GoDataResponseField) Json() ([]byte, error) {
	switch f.Value.(type) {
	case string:
		return prepareJsonString([]byte(f.Value.(string)))
	case []byte:
		return prepareJsonString(f.Value.([]byte))
	case int:
		return []byte(strconv.Itoa(f.Value.(int))), nil
	case float64:
		return []byte(strconv.FormatFloat(f.Value.(float64), 'f', -1, 64)), nil
	case map[string]*GoDataResponseField:
		return prepareJsonDict(f.Value.(map[string]*GoDataResponseField))
	case []*GoDataResponseField:
		return prepareJsonList(f.Value.([]*GoDataResponseField))
	default:
		return nil, InternalServerError("Response field type not recognized.")
	}
}

func prepareJsonString(s []byte) ([]byte, error) {
	// escape double quotes
	s = bytes.Replace(s, []byte("\""), []byte("\\\""), -1)
	var buf bytes.Buffer
	buf.WriteByte('"')
	buf.Write(s)
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

func prepareJsonDict(d map[string]*GoDataResponseField) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	count := 0
	for k, v := range d {
		buf.WriteByte('"')
		buf.Write([]byte(k))
		buf.WriteByte('"')
		buf.WriteByte(':')
		field, err := v.Json()
		if err != nil {
			return nil, err
		}
		buf.Write(field)
		count++
		if count < len(d) {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func prepareJsonList(l []*GoDataResponseField) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	count := 0
	for _, v := range l {
		field, err := v.Json()
		if err != nil {
			return nil, err
		}
		buf.Write(field)
		count++
		if count < len(l) {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}
