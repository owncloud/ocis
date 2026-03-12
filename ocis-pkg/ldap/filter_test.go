package ldap

import (
	"testing"
)

func TestEnhanceFilterWithMasterID(t *testing.T) {
	tests := []struct {
		name           string
		filter         string
		masterID       string
		memberAttr     string
		guestAttr      string
		expectedFilter string
	}{
		{
			name:           "empty masterID",
			filter:         "(ownCloudMemberOf=instance-1)",
			masterID:       "",
			memberAttr:     "ownCloudMemberOf",
			guestAttr:      "ownCloudGuestOf",
			expectedFilter: "(ownCloudMemberOf=instance-1)",
		},
		{
			name:           "no attributes configured",
			filter:         "(uid=testuser)",
			masterID:       "0000-0000-000-0000",
			memberAttr:     "",
			guestAttr:      "",
			expectedFilter: "(uid=testuser)",
		},
		{
			name:           "both member and guest attributes",
			filter:         "(ownCloudMemberOf=instance-1)",
			masterID:       "0000-0000-000-0000",
			memberAttr:     "ownCloudMemberOf",
			guestAttr:      "ownCloudGuestOf",
			expectedFilter: "(|(ownCloudMemberOf=instance-1)(|(ownCloudMemberOf=0000-0000-000-0000)(ownCloudGuestOf=0000-0000-000-0000)))",
		},
		{
			name:           "only member attribute",
			filter:         "(ownCloudMemberOf=instance-1)",
			masterID:       "0000-0000-000-0000",
			memberAttr:     "ownCloudMemberOf",
			guestAttr:      "",
			expectedFilter: "(|(ownCloudMemberOf=instance-1)(ownCloudMemberOf=0000-0000-000-0000))",
		},
		{
			name:           "only guest attribute",
			filter:         "(ownCloudGuestOf=instance-1)",
			masterID:       "0000-0000-000-0000",
			memberAttr:     "",
			guestAttr:      "ownCloudGuestOf",
			expectedFilter: "(|(ownCloudGuestOf=instance-1)(ownCloudGuestOf=0000-0000-000-0000))",
		},
		{
			name:           "empty existing filter",
			filter:         "",
			masterID:       "0000-0000-000-0000",
			memberAttr:     "ownCloudMemberOf",
			guestAttr:      "ownCloudGuestOf",
			expectedFilter: "(|(ownCloudMemberOf=0000-0000-000-0000)(ownCloudGuestOf=0000-0000-000-0000))",
		},
		{
			name:           "LDAP injection protection",
			filter:         "(ownCloudMemberOf=instance-1)",
			masterID:       "0000*)(uid=*",
			memberAttr:     "ownCloudMemberOf",
			guestAttr:      "",
			expectedFilter: "(|(ownCloudMemberOf=instance-1)(ownCloudMemberOf=0000\\2a\\29\\28uid=\\2a))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnhanceFilterWithMasterID(tt.filter, tt.masterID, tt.memberAttr, tt.guestAttr)
			if result != tt.expectedFilter {
				t.Errorf("EnhanceFilterWithMasterID() = %q, want %q", result, tt.expectedFilter)
			}
		})
	}
}
