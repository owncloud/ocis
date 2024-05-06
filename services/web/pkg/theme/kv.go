package theme

import (
	"bytes"
	"encoding/json"

	"dario.cat/mergo"
	"github.com/spf13/afero"
	"github.com/tidwall/sjson"

	"github.com/owncloud/ocis/v2/ocis-pkg/x/io/fsx"
)

// KV is a generic key-value map.
type KV map[string]any

// MergeKV merges the given key-value maps.
func MergeKV(values ...KV) (KV, error) {
	var kv KV

	for _, v := range values {
		err := mergo.Merge(&kv, v, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
	}

	return kv, nil
}

// PatchKV injects the given values into to v.
func PatchKV(v any, values KV) error {
	bv, err := json.Marshal(v)
	if err != nil {
		return err
	}

	nv := string(bv)

	for k, val := range values {
		var err error
		switch val {
		// if the value is nil, we delete the key
		case nil:
			nv, err = sjson.Delete(nv, k)
		default:
			nv, err = sjson.Set(nv, k, val)
		}

		if err != nil {
			return err
		}
	}

	return json.Unmarshal([]byte(nv), v)
}

// LoadKV loads a key-value map from the given file system.
func LoadKV(fsys fsx.FS, p string) (KV, error) {
	f, err := fsys.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var kv KV
	err = json.NewDecoder(f).Decode(&kv)
	if err != nil {
		return nil, err
	}

	return kv, nil
}

// WriteKV writes the given key-value map to the file system.
func WriteKV(fsys fsx.FS, p string, kv KV) error {
	data, err := json.Marshal(kv)
	if err != nil {
		return err
	}

	return afero.WriteReader(fsys, p, bytes.NewReader(data))
}

// UpdateKV updates the key-value map at the given path with the given values.
func UpdateKV(fsys fsx.FS, p string, values KV) error {
	var kv KV

	existing, err := LoadKV(fsys, p)
	if err == nil {
		kv = existing
	}

	err = PatchKV(&kv, values)
	if err != nil {
		return err
	}

	return WriteKV(fsys, p, kv)
}
