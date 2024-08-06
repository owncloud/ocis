package l10n

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslateStruct(t *testing.T) {

	type InnerStruct struct {
		Description string
		DisplayName *string
	}

	type TopLevelStruct struct {
		Description string
		DisplayName *string
		SubStruct   *InnerStruct
	}

	type WrapperStruct struct {
		Description string
		StructList  []*InnerStruct
	}

	toStrPointer := func(str string) *string {
		return &str
	}

	type args struct {
		structPtr []TranslateOption
		//request   []any
	}
	tests := []struct {
		name     string
		entity   any
		args     args
		expected any
		wantErr  bool
	}{
		{
			name: "top level slice of struct",
			entity: []*InnerStruct{
				{
					Description: "inner 1",
					DisplayName: toStrPointer("innerDisplayName 1"),
				},
				{
					Description: "inner 2",
					DisplayName: toStrPointer("innerDisplayName 2"),
				},
			},
			args: args{
				structPtr: []TranslateOption{
					TranslateField("Description"),
					TranslateField("DisplayName"),
				},
			},
			expected: []*InnerStruct{
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
			name: "top level slice of string",
			entity: []string{
				"inner 1",
				"inner 2",
			},
			args: args{},
			expected: []string{
				"new Inner 1",
				"new Inner 2",
			},
		},
		{
			name: "top level slice of struct",
			entity: []*TopLevelStruct{
				{
					Description: "inner 1",
					DisplayName: toStrPointer("innerDisplayName 1"),
					SubStruct: &InnerStruct{
						Description: "inner",
						DisplayName: toStrPointer("innerDisplayName"),
					},
				},
				{
					Description: "inner 2",
					DisplayName: toStrPointer("innerDisplayName 2"),
				},
			},
			args: args{
				structPtr: []TranslateOption{
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"),
					),
				},
			},
			expected: []*TopLevelStruct{
				{
					Description: "new Inner 1",
					DisplayName: toStrPointer("new InnerDisplayName 1"),
					SubStruct: &InnerStruct{
						Description: "new Inner",
						DisplayName: toStrPointer("new InnerDisplayName"),
					},
				},
				{
					Description: "new Inner 2",
					DisplayName: toStrPointer("new InnerDisplayName 2"),
				},
			},
		},
		{
			name: "wrapped struct full",
			entity: &WrapperStruct{
				StructList: []*InnerStruct{
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
			args: args{
				structPtr: []TranslateOption{
					TranslateEach("StructList",
						TranslateField("Description"),
						TranslateField("DisplayName"),
					),
				},
			},
			expected: &WrapperStruct{
				StructList: []*InnerStruct{
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
			name:   "empty struct, NotExistingSubStructName",
			entity: &TopLevelStruct{},
			args: args{
				structPtr: []TranslateOption{
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("NotExistingSubStructName",
						TranslateField("Description"),
						TranslateField("DisplayName"),
					),
				},
			},
			expected: &TopLevelStruct{},
		},
		{
			name:   "empty struct",
			entity: &TopLevelStruct{},
			args: args{
				structPtr: []TranslateOption{
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"),
					),
				},
			},
			expected: &TopLevelStruct{},
		},
		{
			name: "empty struct, not existing field",
			entity: &TopLevelStruct{
				Description: "description",
				DisplayName: toStrPointer("displayName"),
			},
			args: args{
				structPtr: []TranslateOption{
					TranslateField("NotExistingFieldName"),
					TranslateStruct("SubStruct",
						TranslateField("NotExistingFieldName"),
					),
				},
			},
			expected: &TopLevelStruct{
				Description: "description",
				DisplayName: toStrPointer("displayName"),
			},
		},
		{
			name: "inner struct DisplayName empy",
			entity: &TopLevelStruct{
				Description: "description",
				DisplayName: toStrPointer("displayName"),
			},
			args: args{
				structPtr: []TranslateOption{
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"),
					),
				},
			},
			expected: &TopLevelStruct{
				Description: "new Description",
				DisplayName: toStrPointer("new DisplayName"),
			},
		},
		{
			name: "inner struct full",
			entity: &TopLevelStruct{
				Description: "description",
				DisplayName: toStrPointer("displayName"),
			},
			args: args{
				structPtr: []TranslateOption{
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"),
					),
				},
			},
			expected: &TopLevelStruct{
				Description: "new Description",
				DisplayName: toStrPointer("new DisplayName"),
			},
		},
		{
			name: "full struct",
			entity: &TopLevelStruct{
				Description: "description",
				DisplayName: toStrPointer("displayName"),
				SubStruct: &InnerStruct{
					Description: "inner",
					DisplayName: toStrPointer("innerDisplayName"),
				},
			},
			args: args{
				structPtr: []TranslateOption{
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"),
					),
				},
			},
			expected: &TopLevelStruct{
				Description: "new Description",
				DisplayName: toStrPointer("new DisplayName"),
				SubStruct: &InnerStruct{
					Description: "new Inner",
					DisplayName: toStrPointer("new InnerDisplayName"),
				},
			},
		},
		{
			name: "nil",
			args: args{
				structPtr: nil,
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			args: args{
				structPtr: []TranslateOption{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TranslateEntity(mock(), tt.entity, tt.args.structPtr...)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranslateEntity() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, tt.entity)
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
