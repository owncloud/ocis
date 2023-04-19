package iso6391

var (
	// Codes is a list of all ISO 639-1 two character codes.
	Codes []string

	// Names is a list of all ISO 639-1 english language names.
	Names []string

	// NativeNames is a list of all ISO 639-1 native language names.
	NativeNames []string
)

func init() {
	Codes = make([]string, 0, len(Languages))
	Names = make([]string, 0, len(Languages))
	NativeNames = make([]string, 0, len(Languages))

	for key, value := range Languages {
		Codes = append(Codes, key)
		Names = append(Names, value.Name)
		NativeNames = append(NativeNames, value.NativeName)
	}
}
