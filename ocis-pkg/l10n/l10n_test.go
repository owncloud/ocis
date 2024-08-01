package l10n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateStruct(t *testing.T) {

	type InnreStruct struct {
		Description string
		DisplayName *string
	}

	type TopLevelStruct struct {
		Description string
		DisplayName *string
		SubStruct   *InnreStruct
	}

	type WrapperStruct struct {
		StructList []*InnreStruct
	}

	toStrPointer := func(str string) *string {
		return &str
	}

	type args struct {
		structPtr interface{}
		request   []any
	}
	tests := []struct {
		name     string
		args     args
		expected any
		wantErr  bool
	}{
		{
			name: "top level slice of struct",
			args: args{
				structPtr: []*InnreStruct{
					{
						Description: "inner 1",
						DisplayName: toStrPointer("innerDisplayName 1"),
					},
					{
						Description: "inner 2",
						DisplayName: toStrPointer("innerDisplayName 2"),
					},
				},
				request: []any{
					TranslateField("Description"),
					TranslateField("DisplayName")},
			},
			expected: []*InnreStruct{
				{
					Description: "new Inner 1",
					DisplayName: toStrPointer("new InnerDisplayName 1"),
				},
				{
					Description: "new Inner 2",
					DisplayName: toStrPointer("new InnerDisplayName 2"),
				},
			},
		},
		{
			name: "wrapped struct full",
			args: args{
				structPtr: &WrapperStruct{
					StructList: []*InnreStruct{
						{
							Description: "inner 1",
							DisplayName: toStrPointer("innerDisplayName 1"),
						},
						{
							Description: "inner 2",
							DisplayName: toStrPointer("innerDisplayName 2"),
						},
					},
				},
				request: []any{
					TranslateField("StructList",
						TranslateField("Description"),
						TranslateField("DisplayName"))},
			},
			expected: &WrapperStruct{
				StructList: []*InnreStruct{
					{
						Description: "new Inner 1",
						DisplayName: toStrPointer("new InnerDisplayName 1"),
					},
					{
						Description: "new Inner 2",
						DisplayName: toStrPointer("new InnerDisplayName 2"),
					},
				},
			},
		},
		{
			name: "empty struct, NotExistingSubStructName",
			args: args{
				structPtr: &TopLevelStruct{},
				request: []any{
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateField("NotExistingSubStructName",
						TranslateField("Description"),
						TranslateField("DisplayName")),
				},
			},
			expected: &TopLevelStruct{},
		},
		{
			name: "empty struct",
			args: args{
				structPtr: &TopLevelStruct{},
				request: []any{
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateField("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"))},
			},
			expected: &TopLevelStruct{},
		},
		{
			name: "empty struct, not existing field",
			args: args{
				structPtr: &TopLevelStruct{
					Description: "description",
					DisplayName: toStrPointer("displayName"),
				},
				request: []any{
					TranslateField("NotExistingFieldName"),
					TranslateField("SubStruct",
						TranslateField("NotExistingFieldName"))},
			},
			expected: &TopLevelStruct{
				Description: "description",
				DisplayName: toStrPointer("displayName"),
			},
		},
		{
			name: "inner struct DisplayName empy",
			args: args{
				structPtr: &TopLevelStruct{
					Description: "description",
					DisplayName: toStrPointer("displayName"),
				},
				request: []any{TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateField("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"))},
			},
			expected: &TopLevelStruct{
				Description: "new Description",
				DisplayName: toStrPointer("new DisplayName"),
			},
		},
		{
			name: "inner struct full",
			args: args{
				structPtr: &TopLevelStruct{
					Description: "description",
					DisplayName: toStrPointer("displayName"),
				},
				request: []any{TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateField("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"))},
			},
			expected: &TopLevelStruct{
				Description: "new Description",
				DisplayName: toStrPointer("new DisplayName"),
			},
		},
		{
			name: "full struct",
			args: args{
				structPtr: &TopLevelStruct{
					Description: "description",
					DisplayName: toStrPointer("displayName"),
					SubStruct: &InnreStruct{
						Description: "inner",
						DisplayName: toStrPointer("innerDisplayName"),
					},
				},
				request: []any{
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateField("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"))},
			},
			expected: &TopLevelStruct{
				Description: "new Description",
				DisplayName: toStrPointer("new DisplayName"),
				SubStruct: &InnreStruct{
					Description: "new Inner",
					DisplayName: toStrPointer("new InnerDisplayName"),
				},
			},
		},
		{
			name: "nil",
			args: args{
				structPtr: nil,
				request:   []any{TranslateField("Description")},
			},
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TranslateEntity(tt.args.structPtr, mock(), tt.args.request...)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranslateEntity() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, tt.args.structPtr)
		})
	}
}

func mock() func(string, ...interface{}) string {
	return func(s string, i ...interface{}) string {
		switch s {
		case "description":
			return "new Description"
		case "displayName":
			return "new DisplayName"
		case "inner":
			return "new Inner"
		case "innerDisplayName":
			return "new InnerDisplayName"
		case "inner 1":
			return "new Inner 1"
		case "innerDisplayName 1":
			return "new InnerDisplayName 1"
		case "inner 2":
			return "new Inner 2"
		case "innerDisplayName 2":
			return "new InnerDisplayName 2"
		}
		return s
	}
}
