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
// so both fields can have a value from a different scope.
//
// The more specific scope must win. Inside one scope the content wins.
const (
	credsSuffix     = "Creds"
	credsFileSuffix = "CredsFile"
)

// resolveCredsScope clears the creds content of a pair when the creds file of the
// pair comes from a more specific scope. Example: OCIS_NATS_CREDS is set globally
// and GRAPH_CACHE_CREDS_FILE is set for the service. Then the file is used.
func resolveCredsScope(s reflect.Value) {
	t := s.Type()
	for i := 0; i < t.NumField(); i++ {
		credsField := t.Field(i)
		if !credsField.IsExported() || credsField.Type.Kind() != reflect.String || !strings.HasSuffix(credsField.Name, credsSuffix) {
			continue
		}
		prefix := strings.TrimSuffix(credsField.Name, credsSuffix)
		fileField, ok := t.FieldByName(prefix + credsFileSuffix)
		if !ok || fileField.Type.Kind() != reflect.String {
			continue
		}

		creds := s.FieldByIndex(credsField.Index)
		file := s.FieldByIndex(fileField.Index)
		if creds.String() == "" || file.String() == "" {
			continue
		}
		if envScope(fileField.Tag.Get("env")) > envScope(credsField.Tag.Get("env")) {
			creds.SetString("")
		}
	}
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
