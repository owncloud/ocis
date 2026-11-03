package rolemapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePair(t *testing.T) {
	tests := []struct {
		pair     string
		value    string
		roleName string
		wantErr  bool
	}{
		{pair: "lehrer:user-light-vault", value: "lehrer", roleName: "user-light-vault"},
		{pair: " lehrer : user-light-vault ", value: "lehrer", roleName: "user-light-vault"},
		// split at the last colon, the value may contain colons
		{pair: "urn:ocis:lehrer:user-light-vault", value: "urn:ocis:lehrer", roleName: "user-light-vault"},
		{pair: "lehrer", wantErr: true},
		{pair: ":user-light-vault", wantErr: true},
		{pair: "lehrer:", wantErr: true},
		{pair: " :user-light-vault", wantErr: true},
		{pair: "lehrer: ", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.pair, func(t *testing.T) {
			value, roleName, err := ParsePair(tt.pair)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.roleName, roleName)
		})
	}
}

func TestMatcher(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		v       string
		matches bool
	}{
		{name: "exact value", value: "lehrer", v: "lehrer", matches: true},
		{name: "other value", value: "lehrer", v: "schueler", matches: false},
		{name: "case matters", value: "lehrer", v: "Lehrer", matches: false},
		{name: "regex", value: "lehr.*", v: "lehrer", matches: true},
		{name: "regex alternatives", value: "lehrer|Nutzerverwalter", v: "Nutzerverwalter", matches: true},
		{name: "a regex must match the whole value", value: "lehrer", v: "nichtlehrer", matches: false},
		{name: "invalid regex matches its exact value", value: "lehrer[", v: "lehrer[", matches: true},
		{name: "invalid regex matches nothing else", value: "lehrer[", v: "lehrer", matches: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.matches, NewMatcher(tt.value).Matches(tt.v))
		})
	}
}
