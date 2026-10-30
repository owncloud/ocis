package envdecode

import (
	"os"
	"reflect"
	"strings"
)

// NATS JWT credentials use two fields: the path of a .creds file (<prefix>CredsFile)
// and the content of a .creds file (<prefix>Creds). Each field has its own list of
// environment variables, for example "OCIS_NATS_CREDS_FILE;GRAPH_CACHE_CREDS_FILE"
// and "OCIS_NATS_CREDS;GRAPH_CACHE_CREDS". The decoder sets each field on its own,
// so both fields can have a value from a different scope. The same applies to the
// username and the password (<prefix>AuthUsername and <prefix>AuthPassword).
//
// The more specific scope must win. Inside one scope the content wins over the file,
// and the JWT credentials win over the username and the password.
const (
	credsSuffix        = "Creds"
	credsFileSuffix    = "CredsFile"
	authUsernameSuffix = "AuthUsername"
	authPasswordSuffix = "AuthPassword"
)

// resolveCredsScope applies the scope rules to each creds pair of the struct s.
//
//   - It clears the creds content when the creds file comes from a more specific scope.
//     Example: OCIS_NATS_CREDS is set globally and GRAPH_CACHE_CREDS_FILE is set for
//     the service. Then the service uses the file.
//   - It clears the creds content and the creds file when the username and the password
//     both have a value and one of them comes from a more specific environment variable.
//     Example: OCIS_NATS_CREDS_FILE is set globally, and GRAPH_CACHE_AUTH_USERNAME and
//     GRAPH_CACHE_AUTH_PASSWORD are set for the service. Then the connection uses the
//     username and the password. The function keeps credentials from a config file or
//     a default, because these credentials have no environment scope to compare.
func resolveCredsScope(s reflect.Value) {
	t := s.Type()
	for i := 0; i < t.NumField(); i++ {
		credsField := t.Field(i)
		if !credsField.IsExported() || credsField.Type.Kind() != reflect.String || !strings.HasSuffix(credsField.Name, credsSuffix) {
			continue
		}
		prefix := strings.TrimSuffix(credsField.Name, credsSuffix)
		fileField, ok := stringField(t, prefix+credsFileSuffix)
		if !ok {
			continue
		}

		creds := s.FieldByIndex(credsField.Index)
		file := s.FieldByIndex(fileField.Index)
		if creds.String() != "" && file.String() != "" &&
			envScope(fileField.Tag.Get("env")) > envScope(credsField.Tag.Get("env")) {
			creds.SetString("")
		}

		credsScope := max(valueScope(creds, credsField), valueScope(file, fileField))
		if credsScope >= 0 && authScope(s, prefix) > credsScope {
			creds.SetString("")
			file.SetString("")
		}
	}
}

// authScope returns the scope of the username and the password with the given prefix.
// It returns -1 when the username or the password is empty. A username or a password
// alone cannot authenticate, so it does not replace the JWT credentials. Otherwise
// it returns the more specific scope of the two values.
func authScope(s reflect.Value, prefix string) int {
	scope := -1
	for _, suffix := range []string{authUsernameSuffix, authPasswordSuffix} {
		f, ok := stringField(s.Type(), prefix+suffix)
		if !ok {
			return -1
		}
		v := s.FieldByIndex(f.Index)
		if v.String() == "" {
			return -1
		}
		scope = max(scope, valueScope(v, f))
	}
	return scope
}

// stringField returns the field of t with the given name when it is a string field.
func stringField(t reflect.Type, name string) (reflect.StructField, bool) {
	f, ok := t.FieldByName(name)
	if !ok || f.Type.Kind() != reflect.String {
		return reflect.StructField{}, false
	}
	return f, true
}

// valueScope returns the scope of the value v of the field f. It returns -1 when v is empty.
func valueScope(v reflect.Value, f reflect.StructField) int {
	if v.String() == "" {
		return -1
	}
	return envScope(f.Tag.Get("env"))
}

// envScope returns the position of the last environment variable of the tag that is set.
// A later variable is more specific. It returns -1 when no variable is set, so a value
// from a config file or a default has the lowest scope.
func envScope(tag string) int {
	names := strings.Split(strings.Split(tag, ",")[0], ";")
	scope := -1
	for i, name := range names {
		if name == "" {
			continue
		}
		if _, set := os.LookupEnv(name); set {
			scope = i
		}
	}
	return scope
}
