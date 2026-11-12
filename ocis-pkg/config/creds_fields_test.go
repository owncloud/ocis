package config_test

import (
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
)

// The envdecode package finds the NATS credentials of a struct by the field names
// <prefix>Creds, <prefix>CredsFile, <prefix>AuthUsername and <prefix>AuthPassword.
// It skips a struct without an error when one name does not match. This test makes
// sure that the names match in each real config struct.
var _ = Describe("NATS credentials fields", func() {
	It("use the same prefix in each config struct", func() {
		var problems []string
		pairs := 0
		visited := map[reflect.Type]bool{}

		var walk func(t reflect.Type)
		walk = func(t reflect.Type) {
			for t.Kind() == reflect.Pointer || t.Kind() == reflect.Slice || t.Kind() == reflect.Map {
				t = t.Elem()
			}
			if t.Kind() != reflect.Struct || visited[t] || !strings.HasPrefix(t.PkgPath(), "github.com/owncloud/ocis/v2/") {
				return
			}
			visited[t] = true

			for i := 0; i < t.NumField(); i++ {
				f := t.Field(i)
				if !f.IsExported() {
					continue
				}
				walk(f.Type)

				var prefix string
				var want []string
				switch {
				case strings.HasSuffix(f.Name, "CredsFile"):
					prefix = strings.TrimSuffix(f.Name, "CredsFile")
					want = []string{"Creds"}
				case strings.HasSuffix(f.Name, "Creds"):
					prefix = strings.TrimSuffix(f.Name, "Creds")
					want = []string{"CredsFile", "AuthUsername", "AuthPassword"}
					pairs++
				default:
					continue
				}
				if f.Type.Kind() != reflect.String {
					problems = append(problems, t.String()+"."+f.Name+" is not a string")
				}
				for _, suffix := range want {
					other, ok := t.FieldByName(prefix + suffix)
					if !ok || other.Type.Kind() != reflect.String {
						problems = append(problems, t.String()+"."+f.Name+" has no string field "+prefix+suffix)
					}
				}
			}
		}
		walk(reflect.TypeFor[config.Config]())

		Expect(problems).To(BeEmpty())
		// The walk must find the credentials, or the test checks nothing.
		Expect(pairs).To(BeNumerically(">", 20))
	})
})
