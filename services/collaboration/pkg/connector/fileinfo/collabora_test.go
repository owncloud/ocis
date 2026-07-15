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
