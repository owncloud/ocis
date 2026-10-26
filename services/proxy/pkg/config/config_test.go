package config

import (
	"reflect"
	"testing"
)

func TestGuestRoleMappings(t *testing.T) {
	tests := []struct {
		name    string
		pairs   []string
		want    []RoleMapping
		wantErr bool
	}{
		{
			name:  "empty",
			pairs: nil,
			want:  []RoleMapping{},
		},
		{
			name:  "pairs keep their order",
			pairs: []string{"teacher:guest-teacher", "admin:guest-admin"},
			want: []RoleMapping{
				{ClaimValue: "teacher", RoleName: "guest-teacher"},
				{ClaimValue: "admin", RoleName: "guest-admin"},
			},
		},
		{
			name:  "claim value may contain colons",
			pairs: []string{"urn:school:teacher:guest-teacher"},
			want:  []RoleMapping{{ClaimValue: "urn:school:teacher", RoleName: "guest-teacher"}},
		},
		{name: "missing colon", pairs: []string{"teacher"}, wantErr: true},
		{name: "missing claim value", pairs: []string{":guest-teacher"}, wantErr: true},
		{name: "missing role name", pairs: []string{"teacher:"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MultiInstanceConfig{GuestRoleMapping: tt.pairs}.GuestRoleMappings()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error for %v", tt.pairs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
