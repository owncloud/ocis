package iso6391

// FromCode returns the language for given code.
func FromCode(code string) Language {
	return Languages[code]
}

// FromName returns the language for given name.
func FromName(name string) Language {
	for _, lang := range Languages {
		if lang.Name == name {
			return lang
		}
	}

	return Language{}
}

// FromNativeName returns the language for given native name.
func FromNativeName(name string) Language {
	for _, lang := range Languages {
		if lang.NativeName == name {
			return lang
		}
	}

	return Language{}
}

// Name returns the language name in english for given code.
func Name(code string) string {
	return FromCode(code).Name
}

// NativeName returns the language native name for given code.
func NativeName(code string) string {
	return FromCode(code).NativeName
}

// CodeForName returns the language code for given name.
func CodeForName(name string) string {
	for code, lang := range Languages {
		if lang.Name == name {
			return code
		}
	}

	return ""
}

// CodeForNativeName returns the language code for given native name.
func CodeForNativeName(name string) string {
	for code, lang := range Languages {
		if lang.NativeName == name {
			return code
		}
	}

	return ""
}

// ValidCode returns true if the given code is a valid ISO 639-1 language code.
// The code must be passed in lowercase.
func ValidCode(code string) bool {
	_, ok := Languages[code]
	return ok
}

// ValidName returns true if the given name is a valid ISO 639-1 language name.
// The name must use uppercase characters (e.g. English, Hiri Motu, ...).
func ValidName(name string) bool {
	for _, lang := range Languages {
		if lang.Name == name {
			return true
		}
	}

	return false
}

// ValidNativeName returns true if the given code is a valid ISO 639-1 language native name.
// The name must be passed in its native form (e.g. English, 中文, ...).
func ValidNativeName(name string) bool {
	for _, lang := range Languages {
		if lang.NativeName == name {
			return true
		}
	}

	return false
}
