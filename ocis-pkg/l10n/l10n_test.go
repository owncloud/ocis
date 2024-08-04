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
		structPtr interface{}
		//request   []any
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
				structPtr: TranslateEach(
					[]*InnerStruct{
						{
							Description: "inner 1",
							DisplayName: toStrPointer("innerDisplayName 1"),
						},
						{
							Description: "inner 2",
							DisplayName: toStrPointer("innerDisplayName 2"),
						},
					},
					TranslateField("Description"),
					TranslateField("DisplayName")),
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
			args: args{
				structPtr: TranslateEach(
					[]string{
						"inner 1",
						"inner 2",
					}),
			},
			expected: []string{
				"new Inner 1",
				"new Inner 2",
			},
		},
		{
			name: "top level slice of struct",
			args: args{
				structPtr: TranslateEach(
					[]*TopLevelStruct{
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
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"))),
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
			args: args{
				structPtr: TranslateStruct(
					&WrapperStruct{
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
					TranslateEach("StructList",
						TranslateField("Description"),
						TranslateField("DisplayName")),
				),
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
			name: "empty struct, NotExistingSubStructName",
			args: args{
				structPtr: TranslateStruct(
					&TopLevelStruct{},
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("NotExistingSubStructName",
						TranslateField("Description"),
						TranslateField("DisplayName"))),
			},
			expected: &TopLevelStruct{},
		},
		{
			name: "empty struct",
			args: args{
				structPtr: TranslateStruct(
					&TopLevelStruct{},
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"))),
			},
			expected: &TopLevelStruct{},
		},
		{
			name: "empty struct, not existing field",
			args: args{
				structPtr: TranslateStruct(
					&TopLevelStruct{
						Description: "description",
						DisplayName: toStrPointer("displayName"),
					}, TranslateField("NotExistingFieldName"),
					TranslateStruct("SubStruct",
						TranslateField("NotExistingFieldName"))),
			},
			expected: &TopLevelStruct{
				Description: "description",
				DisplayName: toStrPointer("displayName"),
			},
		},
		{
			name: "inner struct DisplayName empy",
			args: args{
				structPtr: TranslateStruct(
					&TopLevelStruct{
						Description: "description",
						DisplayName: toStrPointer("displayName"),
					},
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName"))),
			},
			expected: &TopLevelStruct{
				Description: "new Description",
				DisplayName: toStrPointer("new DisplayName"),
			},
		},
		{
			name: "inner struct full",
			args: args{
				structPtr: TranslateStruct(
					&TopLevelStruct{
						Description: "description",
						DisplayName: toStrPointer("displayName"),
					},
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName")),
				),
			},
			expected: &TopLevelStruct{
				Description: "new Description",
				DisplayName: toStrPointer("new DisplayName"),
			},
		},
		{
			name: "full struct",
			args: args{
				structPtr: TranslateStruct(
					&TopLevelStruct{
						Description: "description",
						DisplayName: toStrPointer("displayName"),
						SubStruct: &InnerStruct{
							Description: "inner",
							DisplayName: toStrPointer("innerDisplayName"),
						}},
					TranslateField("Description"),
					TranslateField("DisplayName"),
					TranslateStruct("SubStruct",
						TranslateField("Description"),
						TranslateField("DisplayName")),
				),
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
				structPtr: []any{TranslateEach([]string{})},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := TranslateEntity(mock(), tt.args.structPtr)
			if (err != nil) != tt.wantErr {
				t.Errorf("TranslateEntity() error = %v, wantErr %v", err, tt.wantErr)
			}
			switch a := tt.args.structPtr.(type) {
			case structs:
				assert.Equal(t, tt.expected, a()[0])
			case maps:
				assert.Equal(t, tt.expected, a()[0])
			case each:
				assert.Equal(t, tt.expected, a()[0])
			}
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
