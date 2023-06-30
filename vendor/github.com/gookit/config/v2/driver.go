package config

import (
	"encoding/json"

	"github.com/gookit/goutil/jsonutil"
)

// Driver interface.
// TODO refactor: rename GetDecoder() to Decode(), rename GetEncoder() to Encode()
type Driver interface {
	Name() string
	Aliases() []string // alias format names, use for resolve format name
	GetDecoder() Decoder
	GetEncoder() Encoder
}

// DriverV2 interface.
type DriverV2 interface {
	Name() string      // driver name, also is format name.
	Aliases() []string // alias format names, use for resolve format name
	Decode(blob []byte, v any) (err error)
	Encode(v any) (out []byte, err error)
}

// Decoder for decode yml,json,toml format content
type Decoder func(blob []byte, v any) (err error)

// Encoder for decode yml,json,toml format content
type Encoder func(v any) (out []byte, err error)

// StdDriver struct
type StdDriver struct {
	name    string
	aliases []string
	decoder Decoder
	encoder Encoder
}

// NewDriver new std driver instance.
func NewDriver(name string, dec Decoder, enc Encoder) *StdDriver {
	return &StdDriver{name: name, decoder: dec, encoder: enc}
}

// WithAliases set aliases for driver
func (d *StdDriver) WithAliases(aliases ...string) *StdDriver {
	d.aliases = aliases
	return d
}

// WithAlias add alias for driver
func (d *StdDriver) WithAlias(alias string) *StdDriver {
	d.aliases = append(d.aliases, alias)
	return d
}

// Name of driver
func (d *StdDriver) Name() string { return d.name }

// Aliases format name of driver
func (d *StdDriver) Aliases() []string {
	return d.aliases
}

// Decode of driver
func (d *StdDriver) Decode(blob []byte, v any) (err error) {
	return d.decoder(blob, v)
}

// Encode of driver
func (d *StdDriver) Encode(v any) ([]byte, error) {
	return d.encoder(v)
}

// GetDecoder of driver
func (d *StdDriver) GetDecoder() Decoder {
	return d.decoder
}

// GetEncoder of driver
func (d *StdDriver) GetEncoder() Encoder {
	return d.encoder
}

/*************************************************************
 * JSON driver
 *************************************************************/

var (
	// JSONAllowComments support write comments on json file.
	JSONAllowComments = true

	// JSONMarshalIndent if not empty, will use json.MarshalIndent for encode data.
	//
	// Deprecated: please use JSONDriver.MarshalIndent
	JSONMarshalIndent string
)

// JSONDecoder for json decode
var JSONDecoder Decoder = func(data []byte, v any) (err error) {
	JSONDriver.ClearComments = JSONAllowComments
	return JSONDriver.Decode(data, v)
}

// JSONEncoder for json encode
var JSONEncoder Encoder = func(v any) (out []byte, err error) {
	JSONDriver.MarshalIndent = JSONMarshalIndent
	return JSONDriver.Encode(v)
}

// JSONDriver instance fot json
var JSONDriver = &jsonDriver{
	driverName:    JSON,
	ClearComments: JSONAllowComments,
	MarshalIndent: JSONMarshalIndent,
}

// jsonDriver for json format content
type jsonDriver struct {
	driverName string
	// ClearComments before parse JSON string.
	ClearComments bool
	// MarshalIndent if not empty, will use json.MarshalIndent for encode data.
	MarshalIndent string
}

// Name of the driver
func (d *jsonDriver) Name() string {
	return d.driverName
}

// Aliases of the driver
func (d *jsonDriver) Aliases() []string {
	return nil
}

// Decode for the driver
func (d *jsonDriver) Decode(data []byte, v any) error {
	if d.ClearComments {
		str := jsonutil.StripComments(string(data))
		return json.Unmarshal([]byte(str), v)
	}
	return json.Unmarshal(data, v)
}

// GetDecoder for the driver
func (d *jsonDriver) GetDecoder() Decoder {
	return d.Decode
}

// Encode for the driver
func (d *jsonDriver) Encode(v any) (out []byte, err error) {
	if len(d.MarshalIndent) > 0 {
		return json.MarshalIndent(v, "", d.MarshalIndent)
	}
	return json.Marshal(v)
}

// GetEncoder for the driver
func (d *jsonDriver) GetEncoder() Encoder {
	return d.Encode
}
