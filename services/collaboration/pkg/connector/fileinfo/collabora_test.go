package fileinfo

import (
	"encoding/json"
	"testing"
)

func TestCollaboraSetPropertiesVersion(t *testing.T) {
	tests := []struct {
		name           string
		version        string
		expectedResult string
	}{
		{
			name:           "Version set correctly",
			version:        "v42",
			expectedResult: "v42",
		},
		{
			name:           "Empty version",
			version:        "",
			expectedResult: "",
		},
		{
			name:           "Complex version string",
			version:        "2024-07-16T12:34:56Z",
			expectedResult: "2024-07-16T12:34:56Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{}
			cinfo.SetProperties(map[string]interface{}{
				KeyVersion: tt.version,
			})

			if cinfo.Version != tt.expectedResult {
				t.Errorf("SetProperties Version: got %q, want %q", cinfo.Version, tt.expectedResult)
			}
		})
	}
}

func TestCollaboraJSONMarshallingIncludesVersion(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		shouldInclude bool
	}{
		{
			name:          "Version with value should be included",
			version:       "abc123",
			shouldInclude: true,
		},
		{
			name:          "Empty version should be omitted",
			version:       "",
			shouldInclude: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{
				Version: tt.version,
			}

			jsonBytes, err := json.Marshal(cinfo)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			var jsonMap map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			_, hasVersion := jsonMap["Version"]
			if hasVersion != tt.shouldInclude {
				if tt.shouldInclude {
					t.Errorf("Expected Version field to be in JSON, but it was not")
				} else {
					t.Errorf("Expected Version field to be omitted from JSON due to omitempty, but it was present")
				}
			}

			if tt.shouldInclude && tt.version != "" {
				if val, ok := jsonMap["Version"].(string); !ok || val != tt.version {
					t.Errorf("Version field: got %v, want %q", jsonMap["Version"], tt.version)
				}
			}
		})
	}
}

func TestCollaboraSetPropertiesLastModifiedTime(t *testing.T) {
	tests := []struct {
		name             string
		lastModifiedTime string
		expectedResult   string
	}{
		{
			name:             "LastModifiedTime set correctly",
			lastModifiedTime: "2025-04-13T10:00:00.0000000Z",
			expectedResult:   "2025-04-13T10:00:00.0000000Z",
		},
		{
			name:             "Empty LastModifiedTime",
			lastModifiedTime: "",
			expectedResult:   "",
		},
		{
			name:             "Complex LastModifiedTime string",
			lastModifiedTime: "2024-07-16T12:34:56.1234567Z",
			expectedResult:   "2024-07-16T12:34:56.1234567Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{}
			cinfo.SetProperties(map[string]interface{}{
				KeyLastModifiedTime: tt.lastModifiedTime,
			})

			if cinfo.LastModifiedTime != tt.expectedResult {
				t.Errorf("SetProperties LastModifiedTime: got %q, want %q", cinfo.LastModifiedTime, tt.expectedResult)
			}
		})
	}
}

func TestCollaboraJSONMarshallingIncludesLastModifiedTime(t *testing.T) {
	tests := []struct {
		name             string
		lastModifiedTime string
		shouldInclude    bool
	}{
		{
			name:             "LastModifiedTime with value should be included",
			lastModifiedTime: "2025-04-13T10:00:00.0000000Z",
			shouldInclude:    true,
		},
		{
			name:             "Empty LastModifiedTime should be omitted",
			lastModifiedTime: "",
			shouldInclude:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{
				LastModifiedTime: tt.lastModifiedTime,
			}

			jsonBytes, err := json.Marshal(cinfo)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			var jsonMap map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			_, hasLastModifiedTime := jsonMap["LastModifiedTime"]
			if hasLastModifiedTime != tt.shouldInclude {
				if tt.shouldInclude {
					t.Errorf("Expected LastModifiedTime field to be in JSON, but it was not")
				} else {
					t.Errorf("Expected LastModifiedTime field to be omitted from JSON due to omitempty, but it was present")
				}
			}

			if tt.shouldInclude && tt.lastModifiedTime != "" {
				if val, ok := jsonMap["LastModifiedTime"].(string); !ok || val != tt.lastModifiedTime {
					t.Errorf("LastModifiedTime field: got %v, want %q", jsonMap["LastModifiedTime"], tt.lastModifiedTime)
				}
			}
		})
	}
}

func TestCollaboraLastModifiedTimeOmitEmpty(t *testing.T) {
	tests := []struct {
		name             string
		lastModifiedTime string
		expected         string
	}{
		{
			name:             "LastModifiedTime should be in JSON",
			lastModifiedTime: "2025-04-13T10:00:00.0000000Z",
			expected:         `"LastModifiedTime":"2025-04-13T10:00:00.0000000Z"`,
		},
		{
			name:             "Empty LastModifiedTime should not appear in JSON",
			lastModifiedTime: "",
			expected:         ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{LastModifiedTime: tt.lastModifiedTime}
			jsonBytes, err := json.Marshal(cinfo)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			jsonStr := string(jsonBytes)
			if tt.lastModifiedTime != "" {
				if !json.Valid(jsonBytes) {
					t.Errorf("Invalid JSON produced: %s", jsonStr)
				}

				var result map[string]interface{}
				if err := json.Unmarshal(jsonBytes, &result); err != nil {
					t.Errorf("Failed to unmarshal JSON: %v", err)
				}

				if v, ok := result["LastModifiedTime"]; !ok || v != tt.lastModifiedTime {
					t.Errorf("LastModifiedTime not correctly included in JSON, got: %v", v)
				}
			} else {
				// For empty LastModifiedTime, just verify it doesn't appear in the JSON string
				if json.Valid(jsonBytes) {
					var result map[string]interface{}
					if err := json.Unmarshal(jsonBytes, &result); err == nil {
						if _, hasLastModifiedTime := result["LastModifiedTime"]; hasLastModifiedTime {
							t.Errorf("LastModifiedTime should be omitted from JSON when empty, but was present")
						}
					}
				}
			}
		})
	}
}

func TestCollaboraSetPropertiesReadOnly(t *testing.T) {
	tests := []struct {
		name           string
		readOnly       bool
		expectedResult bool
	}{
		{
			name:           "ReadOnly set to true",
			readOnly:       true,
			expectedResult: true,
		},
		{
			name:           "ReadOnly set to false",
			readOnly:       false,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{}
			cinfo.SetProperties(map[string]interface{}{
				KeyReadOnly: tt.readOnly,
			})

			if cinfo.ReadOnly != tt.expectedResult {
				t.Errorf("SetProperties ReadOnly: got %v, want %v", cinfo.ReadOnly, tt.expectedResult)
			}
		})
	}
}

func TestCollaboraJSONMarshallingIncludesReadOnly(t *testing.T) {
	// ReadOnly has no `omitempty`, so it must always be present in the JSON
	// output, including when false.
	tests := []struct {
		name     string
		readOnly bool
	}{
		{
			name:     "ReadOnly true is included",
			readOnly: true,
		},
		{
			name:     "ReadOnly false is included",
			readOnly: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{
				ReadOnly: tt.readOnly,
			}

			jsonBytes, err := json.Marshal(cinfo)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			var jsonMap map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			val, hasReadOnly := jsonMap["ReadOnly"]
			if !hasReadOnly {
				t.Fatalf("Expected ReadOnly field to always be present in JSON (no omitempty), but it was missing: %s", string(jsonBytes))
			}

			boolVal, ok := val.(bool)
			if !ok || boolVal != tt.readOnly {
				t.Errorf("ReadOnly field: got %v, want %v", val, tt.readOnly)
			}
		})
	}
}

func TestCollaboraVersionOmitEmpty(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "Version abc123 should be in JSON",
			version:  "abc123",
			expected: `"Version":"abc123"`,
		},
		{
			name:     "Empty version should not appear in JSON",
			version:  "",
			expected: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{Version: tt.version}
			jsonBytes, err := json.Marshal(cinfo)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			jsonStr := string(jsonBytes)
			if tt.version != "" {
				if !json.Valid(jsonBytes) {
					t.Errorf("Invalid JSON produced: %s", jsonStr)
				}

				var result map[string]interface{}
				if err := json.Unmarshal(jsonBytes, &result); err != nil {
					t.Errorf("Failed to unmarshal JSON: %v", err)
				}

				if v, ok := result["Version"]; !ok || v != tt.version {
					t.Errorf("Version not correctly included in JSON, got: %v", v)
				}
			} else {
				// For empty version, just verify it doesn't appear in the JSON string
				if json.Valid(jsonBytes) {
					var result map[string]interface{}
					if err := json.Unmarshal(jsonBytes, &result); err == nil {
						if _, hasVersion := result["Version"]; hasVersion {
							t.Errorf("Version should be omitted from JSON when empty, but was present")
						}
					}
				}
			}
		})
	}
}

func TestCollaboraSetPropertiesSupportsUpdate(t *testing.T) {
	tests := []struct {
		name           string
		supportsUpdate bool
		expectedResult bool
	}{
		{
			name:           "SupportsUpdate set to true",
			supportsUpdate: true,
			expectedResult: true,
		},
		{
			name:           "SupportsUpdate set to false",
			supportsUpdate: false,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{}
			cinfo.SetProperties(map[string]interface{}{
				KeySupportsUpdate: tt.supportsUpdate,
			})

			if cinfo.SupportsUpdate != tt.expectedResult {
				t.Errorf("SetProperties SupportsUpdate: got %v, want %v", cinfo.SupportsUpdate, tt.expectedResult)
			}
		})
	}
}

func TestCollaboraJSONMarshallingIncludesSupportsUpdate(t *testing.T) {
	// SupportsUpdate has no `omitempty`, so it must always be present in the JSON
	// output, including when false.
	tests := []struct {
		name           string
		supportsUpdate bool
	}{
		{
			name:           "SupportsUpdate true is included",
			supportsUpdate: true,
		},
		{
			name:           "SupportsUpdate false is included",
			supportsUpdate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{
				SupportsUpdate: tt.supportsUpdate,
			}

			jsonBytes, err := json.Marshal(cinfo)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			var jsonMap map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			val, hasSupportsUpdate := jsonMap["SupportsUpdate"]
			if !hasSupportsUpdate {
				t.Fatalf("Expected SupportsUpdate field to always be present in JSON (no omitempty), but it was missing: %s", string(jsonBytes))
			}

			boolVal, ok := val.(bool)
			if !ok || boolVal != tt.supportsUpdate {
				t.Errorf("SupportsUpdate field: got %v, want %v", val, tt.supportsUpdate)
			}
		})
	}
}

func TestCollaboraSetPropertiesIsAnonymousUser(t *testing.T) {
	tests := []struct {
		name           string
		isAnonymous    bool
		expectedResult bool
	}{
		{
			name:           "IsAnonymousUser set to true",
			isAnonymous:    true,
			expectedResult: true,
		},
		{
			name:           "IsAnonymousUser set to false",
			isAnonymous:    false,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cinfo := &Collabora{}
			cinfo.SetProperties(map[string]interface{}{
				KeyIsAnonymousUser: tt.isAnonymous,
			})

			if cinfo.IsAnonymousUser != tt.expectedResult {
				t.Errorf("SetProperties IsAnonymousUser: got %v, want %v", cinfo.IsAnonymousUser, tt.expectedResult)
			}
		})
	}
}

func TestCollaboraJSONOmitsIsAnonymousUserWhenFalse(t *testing.T) {
	// IsAnonymousUser has `omitempty`, so it should be omitted from JSON when false
	cinfo := &Collabora{
		IsAnonymousUser: false,
	}

	jsonBytes, err := json.Marshal(cinfo)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	_, hasIsAnonymousUser := jsonMap["IsAnonymousUser"]
	if hasIsAnonymousUser {
		t.Errorf("Expected IsAnonymousUser field to be omitted from JSON when false, but it was present: %s", string(jsonBytes))
	}
}

func TestCollaboraJSONIncludesIsAnonymousUserWhenTrue(t *testing.T) {
	// IsAnonymousUser has `omitempty`, so it should be included in JSON when true
	cinfo := &Collabora{
		IsAnonymousUser: true,
	}

	jsonBytes, err := json.Marshal(cinfo)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	val, hasIsAnonymousUser := jsonMap["IsAnonymousUser"]
	if !hasIsAnonymousUser {
		t.Errorf("Expected IsAnonymousUser field to be in JSON when true, but it was not: %s", string(jsonBytes))
	}

	boolVal, ok := val.(bool)
	if !ok || !boolVal {
		t.Errorf("IsAnonymousUser field: got %v, want true", val)
	}
}

// urlPropertyCases describes each of the 7 WOPI URL properties added to Collabora: the
// SetProperties key that populates it, the JSON tag it is expected to be marshalled
// under, a getter for its value, and a distinct sample value for it. Used to drive the
// table-driven tests below.
var urlPropertyCases = []struct {
	name    string
	key     string
	jsonTag string
	value   string
	getter  func(*Collabora) string
}{
	{"CloseURL", KeyCloseURL, "CloseUrl", "https://ocis.example.prv/close", func(c *Collabora) string { return c.CloseURL }},
	{"DownloadURL", KeyDownloadURL, "DownloadUrl", "https://ocis.example.prv/download", func(c *Collabora) string { return c.DownloadURL }},
	{"FileSharingURL", KeyFileSharingURL, "FileSharingUrl", "https://ocis.example.prv/share", func(c *Collabora) string { return c.FileSharingURL }},
	{"FileVersionURL", KeyFileVersionURL, "FileVersionUrl", "https://ocis.example.prv/versions", func(c *Collabora) string { return c.FileVersionURL }},
	{"HostEditURL", KeyHostEditURL, "HostEditUrl", "https://ocis.example.prv/edit", func(c *Collabora) string { return c.HostEditURL }},
	{"HostViewURL", KeyHostViewURL, "HostViewUrl", "https://ocis.example.prv/view", func(c *Collabora) string { return c.HostViewURL }},
	{"SignoutURL", KeySignoutURL, "SignoutUrl", "https://ocis.example.prv/signout", func(c *Collabora) string { return c.SignoutURL }},
}

// TestCollaboraSetPropertiesURLProperties covers brief item 1: for each of the 7 URL
// properties, SetProperties must correctly set the corresponding struct field.
func TestCollaboraSetPropertiesURLProperties(t *testing.T) {
	for _, prop := range urlPropertyCases {
		t.Run(prop.name, func(t *testing.T) {
			cinfo := &Collabora{}
			cinfo.SetProperties(map[string]interface{}{
				prop.key: prop.value,
			})

			if got := prop.getter(cinfo); got != prop.value {
				t.Errorf("SetProperties %s: got %q, want %q", prop.name, got, prop.value)
			}
		})
	}
}

// TestCollaboraJSONMarshallingIncludesURLProperties covers brief item 2: marshalling a
// Collabora struct with all 7 URL properties populated must include all of them in the
// resulting JSON, under their MS-WOPI-compatible tag names.
func TestCollaboraJSONMarshallingIncludesURLProperties(t *testing.T) {
	cinfo := &Collabora{}
	for _, prop := range urlPropertyCases {
		cinfo.SetProperties(map[string]interface{}{prop.key: prop.value})
	}

	jsonBytes, err := json.Marshal(cinfo)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	for _, prop := range urlPropertyCases {
		val, ok := jsonMap[prop.jsonTag]
		if !ok {
			t.Errorf("Expected %q field to be in JSON, but it was not: %s", prop.jsonTag, string(jsonBytes))
			continue
		}
		if val != prop.value {
			t.Errorf("%s field: got %v, want %q", prop.jsonTag, val, prop.value)
		}
	}
}

// TestCollaboraJSONOmitsEmptyURLProperties covers brief item 3: each of the 7 URL
// properties has `omitempty` and must be absent from the JSON output when empty.
func TestCollaboraJSONOmitsEmptyURLProperties(t *testing.T) {
	cinfo := &Collabora{}

	jsonBytes, err := json.Marshal(cinfo)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	for _, prop := range urlPropertyCases {
		if _, ok := jsonMap[prop.jsonTag]; ok {
			t.Errorf("Expected %q field to be omitted from JSON when empty due to omitempty, but it was present: %s", prop.jsonTag, string(jsonBytes))
		}
	}
}

// TestCollaboraSetPropertiesAllURLPropertiesSimultaneously covers brief item 4: setting
// all 7 URL properties in a single SetProperties call must set all of them correctly.
func TestCollaboraSetPropertiesAllURLPropertiesSimultaneously(t *testing.T) {
	props := make(map[string]interface{}, len(urlPropertyCases))
	for _, prop := range urlPropertyCases {
		props[prop.key] = prop.value
	}

	cinfo := &Collabora{}
	cinfo.SetProperties(props)

	for _, prop := range urlPropertyCases {
		if got := prop.getter(cinfo); got != prop.value {
			t.Errorf("SetProperties %s: got %q, want %q", prop.name, got, prop.value)
		}
	}
}
