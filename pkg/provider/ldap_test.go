package provider

import (
	"testing"

	"github.com/CiscoM31/godata"
	"github.com/owncloud/ocis-accounts/pkg/config"
)

var c *config.LDAPSchema

func init() {
	c = &config.LDAPSchema{
		AccountID:   "ownclouduuid",
		Username:    "uid",
		Mail:        "mail",
		DisplayName: "displayname",
	}
}

func TestEQ(t *testing.T) { testLDAPFilters(t, "accountid eq 'a-b-c-d'", "(ownclouduuid=a-b-c-d)") }
func TestNE(t *testing.T) { testLDAPFilters(t, "mail ne 'foo@bar.com'", "(!(mail=foo@bar.com))") }
func TestGE(t *testing.T) { testLDAPFilters(t, "displayname ge 'marie'", "(displayname>=marie)") }
func TestLE(t *testing.T) { testLDAPFilters(t, "username le 'marie'", "(uid<=marie)") }

//func TestHas(t *testing.T) { testLDAPFilters(t, "Style has Sales.Color'Yellow'", "(foo=*)") }
func TestAP(t *testing.T) {
	testLDAPFilters(t, "displayname ap 'einstein'", "(displayname~=einstein)")
}
func TestAND(t *testing.T) {
	testLDAPFilters(t, "accountid le 500000 and accountid ge 300000", "(&(ownclouduuid<=500000)(ownclouduuid>=300000))")
}
func TestOR(t *testing.T) {
	testLDAPFilters(t, "accountid le 700000 or accountid ge 900000", "(|(ownclouduuid<=700000)(ownclouduuid>=900000))")
}
func TestNOT(t *testing.T) {
	// not operator takes precedence over ap, so we need brackets
	testLDAPFilters(t, "not ( displayname ap 'einstein' )", "(!(displayname~=einstein))")
}
func TestContains(t *testing.T) {
	testLDAPFilters(t, "contains(username,'eins')", "(uid=*eins*)")
}
func TestStartsWith(t *testing.T) {
	testLDAPFilters(t, "startswith(username,'eins')", "(uid=eins*)")
}
func TestEndsWith(t *testing.T) {
	testLDAPFilters(t, "endswith(username,'eins')", "(uid=*eins)")
}
func TestEncoding(t *testing.T) {
	testLDAPFilters(t, "displayname eq 'eins(*)tein'", "(displayname=eins\\28\\2a\\29tein)")
}

func testLDAPFilters(t *testing.T, have string, want string) {

	var err error

	var q *godata.GoDataFilterQuery
	if q, err = godata.ParseFilterString(have); err != nil {
		t.Error(err)
	}
	var filter string
	if filter, err = BuildLDAPFilter(q, c); err != nil {
		t.Error(err)
	}
	if filter != want {
		t.Error("expected", want, "for", have, "but got", filter)
	}
}
