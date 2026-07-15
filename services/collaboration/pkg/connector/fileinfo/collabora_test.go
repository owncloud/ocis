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
