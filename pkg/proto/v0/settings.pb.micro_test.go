package proto_test

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/owncloud/ocis-pkg/v2/service/grpc"
	"github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	svc "github.com/owncloud/ocis-settings/pkg/service/v0"
	"github.com/stretchr/testify/assert"
)

var service = grpc.Service{}

var (
	dummySettings = []*proto.Setting{
		{
			DisplayName: "dummy setting",
			SettingKey:  "dummy-setting",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.IntSetting{
					Default: 42,
				},
			},
			Description: "dummy setting",
		},
	}
)

func init() {
	service = grpc.NewService(
		grpc.Namespace("com.owncloud.api"),
		grpc.Name("settings"),
		grpc.Address("localhost:9992"),
	)

	cfg := config.New()
	err := proto.RegisterBundleServiceHandler(service.Server(), svc.NewService(cfg))
	if err != nil {
		log.Fatalf("could not register BundleServiceHandler: %v", err)
	}
	err = proto.RegisterValueServiceHandler(service.Server(), svc.NewService(cfg))
	if err != nil {
		log.Fatalf("could not register ValueServiceHandler: %v", err)
	}
	_ = service.Server().Start()
}

type CustomError struct {
	ID     string
	Code   int
	Detail string
	Status string
}

/**
testing that saving a settings bundle and retrieving it again works correctly
using various setting bundle properties
*/
func TestSettingsBundleProperties(t *testing.T) {
	type TestStruct struct {
		testDataName  string
		BundleKey     string
		SettingKey    string
		DisplayName   string
		Extension     string
		UUID          string
		expectedError CustomError
	}

	var tests = []TestStruct{
		{
			"ASCII",
			"simple-bundle-key",
			"simple-key",
			"simple-display-name",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{},
		},
		{
			"UTF disallowed on keys",
			"सिम्प्ले-bundle-key",
			"सिम्प्ले-key",
			"सिम्प्ले-display-name",
			"सिम्प्ले-extension-name",
			"सिम्प्ले",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "bundle_key: must be in a valid format; extension: must be in a valid format.",
				Status: "Internal Server Error",
			},
		},
		{
			"UTF allowed on display name",
			"simple-bundle-key",
			"simple-key",
			"सिम्प्ले-display-name",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{},
		},
		{
			"bundle key with ../ in the name",
			"../file-a-level-higher-up",
			"simple-key",
			"simple-display-name",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "bundle_key: must be in a valid format.",
				Status: "Internal Server Error",
			},
		},
		{
			"bundle key in the root directory",
			"/tmp/file",
			"simple-key",
			"simple-display-name",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "bundle_key: must be in a valid format.",
				Status: "Internal Server Error",
			},
		},
		{
			"extension name with ../ in the name",
			"simple-bundle-key",
			"simple-key",
			"simple-display-name",
			"../folder-a-level-higher-up",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "extension: must be in a valid format.",
				Status: "Internal Server Error",
			},
		},
		{
			"extension name with \\ in the name",
			"simple-bundle-key",
			"simple-key",
			"simple-display-name",
			"\\",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "extension: must be in a valid format.",
				Status: "Internal Server Error",
			},
		},
		{
			"bundle key with \\ as the name",
			"\\",
			"simple-key",
			"simple-display-name",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "bundle_key: must be in a valid format.",
				Status: "Internal Server Error",
			},
		},
		{
			"spaces in values",
			"simple bundle key",
			"simple key",
			"simple display name",
			"simple extension name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "bundle_key: must be in a valid format; extension: must be in a valid format.",
				Status: "Internal Server Error",
			},
		},
		{
			"bundle key missing",
			"",
			"simple-bundle-key",
			"simple-display-name",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "bundle_key: cannot be blank.",
				Status: "Internal Server Error",
			},
		},
		{
			"extension missing",
			"simple-bundle-key",
			"simple-key",
			"simple-display-name",
			"",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "extension: cannot be blank.",
				Status: "Internal Server Error",
			},
		},
		{
			"setting key missing (omitted on bundles)",
			"simple-bundle-key",
			"",
			"simple-display-name",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{},
		},
		{
			"display name missing",
			"simple-bundle-key",
			"simple-key",
			"",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "display_name: cannot be blank.",
				Status: "Internal Server Error",
			},
		},
		{
			"UUID missing (omitted on bundles)",
			"simple-bundle-key",
			"simple-key",
			"simple-display-name",
			"simple-extension-name",
			"",
			CustomError{},
		},
	}
	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.testDataName, func(t *testing.T) {
			identifier := proto.Identifier{
				Extension:   testCase.Extension,
				BundleKey:   testCase.BundleKey,
				SettingKey:  testCase.SettingKey,
				AccountUuid: testCase.UUID,
			}
			bundle := proto.SettingsBundle{
				Identifier:  &identifier,
				DisplayName: testCase.DisplayName,
				Settings:    dummySettings,
			}
			createRequest := proto.SaveSettingsBundleRequest{
				SettingsBundle: &bundle,
			}

			client := service.Client()
			cl := proto.NewBundleService("com.owncloud.api.settings", client)

			cresponse, err := cl.SaveSettingsBundle(context.Background(), &createRequest)
			if err != nil || (CustomError{} != testCase.expectedError) {
				assert.Error(t, err)
				var errorData CustomError
				_ = json.Unmarshal([]byte(err.Error()), &errorData)
				assert.Equal(t, testCase.expectedError.ID, errorData.ID)
				assert.Equal(t, testCase.expectedError.Code, errorData.Code)
				assert.Equal(t, testCase.expectedError.Detail, errorData.Detail)
				assert.Equal(t, testCase.expectedError.Status, errorData.Status)
			} else {
				assert.Equal(t, testCase.Extension, cresponse.SettingsBundle.Identifier.Extension)
				assert.Equal(t, testCase.BundleKey, cresponse.SettingsBundle.Identifier.BundleKey)
				assert.Equal(t, testCase.SettingKey, cresponse.SettingsBundle.Identifier.SettingKey)
				assert.Equal(t, testCase.UUID, cresponse.SettingsBundle.Identifier.AccountUuid)
				assert.Equal(t, testCase.DisplayName, cresponse.SettingsBundle.DisplayName)

				getRequest := proto.GetSettingsBundleRequest{Identifier: &identifier}
				getResponse, err := cl.GetSettingsBundle(context.Background(), &getRequest)
				assert.NoError(t, err)
				assert.Equal(t, testCase.Extension, getResponse.SettingsBundle.Identifier.Extension)
				assert.Equal(t, testCase.BundleKey, getResponse.SettingsBundle.Identifier.BundleKey)
				assert.Equal(t, testCase.SettingKey, getResponse.SettingsBundle.Identifier.SettingKey)
				assert.Equal(t, testCase.UUID, getResponse.SettingsBundle.Identifier.AccountUuid)
				assert.Equal(t, testCase.DisplayName, getResponse.SettingsBundle.DisplayName)
			}
			_ = os.RemoveAll("ocis-settings-store")
		})
	}
}

func TestSettingsBundleWithoutSettings(t *testing.T) {
	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)

	createRequest := proto.SaveSettingsBundleRequest{
		SettingsBundle: &proto.SettingsBundle{
			Identifier: &proto.Identifier{
				Extension: "great-extension",
				BundleKey: "alices-bundle",
			},
			DisplayName: "Alice's Bundle",
		},
	}
	response, err := cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.Error(t, err)
	assert.Nil(t, response)
	var errorData CustomError
	_ = json.Unmarshal([]byte(err.Error()), &errorData)
	assert.Equal(t, "go.micro.client", errorData.ID)
	assert.Equal(t, 500, errorData.Code)
	assert.Equal(t, "settings: cannot be blank.", errorData.Detail)
	assert.Equal(t, "Internal Server Error", errorData.Status)
	_ = os.RemoveAll("ocis-settings-store")
}

/**
testing that setting getting and listing a settings bundle works correctly with a set of setting definitions
*/
func TestSaveGetListSettingsBundle(t *testing.T) {
	identifier := proto.Identifier{
		Extension:   "my-extension",
		BundleKey:   "simple-bundle-with-setting",
		SettingKey:  "simple-key",
		AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
	}
	var settings []*proto.Setting

	intSetting := proto.IntSetting{
		Default:     1,
		Min:         1,
		Max:         124,
		Step:        1,
		Placeholder: "Int value",
	}
	settings = append(settings, &proto.Setting{
		SettingKey:  "int",
		DisplayName: "an integer value",
		Description: "with some description",
		Value: &proto.Setting_IntValue{
			IntValue: &intSetting,
		},
	})

	stringSetting := proto.StringSetting{
		Default:     "the default value",
		Required:    false,
		MinLength:   2,
		MaxLength:   255,
		Placeholder: "a string value",
	}
	settings = append(settings, &proto.Setting{
		SettingKey:  "string",
		DisplayName: "a string value",
		Description: "with some description",
		Value: &proto.Setting_StringValue{
			StringValue: &stringSetting,
		},
	})

	boolSetting := proto.BoolSetting{
		Default: false,
		Label:   "bool setting",
	}
	settings = append(settings, &proto.Setting{
		SettingKey:  "bool",
		DisplayName: "a bool value",
		Description: "with some description",
		Value: &proto.Setting_BoolValue{
			BoolValue: &boolSetting,
		},
	})

	//options for choice list settings
	var options []*proto.ListOption
	options = append(options, &proto.ListOption{
		Value: &proto.ListOptionValue{
			Option: &proto.ListOptionValue_StringValue{StringValue: "list option string value"},
		},
		Default:      true,
		DisplayValue: "a string value",
	})

	options = append(options, &proto.ListOption{
		Value: &proto.ListOptionValue{
			Option: &proto.ListOptionValue_IntValue{IntValue: 123},
		},
		Default:      true,
		DisplayValue: "a int value",
	})

	//MultiChoiceListSetting
	multipleChoiceSetting := proto.MultiChoiceListSetting{
		Options: options,
	}

	settings = append(settings, &proto.Setting{
		SettingKey:  "multiple choice",
		DisplayName: "a multiple choice setting",
		Description: "with some description",
		Value: &proto.Setting_MultiChoiceValue{
			MultiChoiceValue: &multipleChoiceSetting,
		},
	})

	//SingleChoiceListSetting
	singleChoiceSetting := proto.SingleChoiceListSetting{
		Options: options,
	}

	settings = append(settings, &proto.Setting{
		SettingKey:  "single choice",
		DisplayName: "a single choice setting",
		Description: "with some description",
		Value: &proto.Setting_SingleChoiceValue{
			SingleChoiceValue: &singleChoiceSetting,
		},
	})

	bundle := proto.SettingsBundle{
		Identifier:  &identifier,
		DisplayName: "bundle display name",
		Settings:    settings,
	}
	saveRequest := proto.SaveSettingsBundleRequest{
		SettingsBundle: &bundle,
	}

	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)

	//assert that SaveSettingsBundle returns the same bundle as we have sent there
	saveResponse, err := cl.SaveSettingsBundle(context.Background(), &saveRequest)
	assert.NoError(t, err)
	receivedBundle, _ := json.Marshal(saveResponse.SettingsBundle)
	expectedBundle, _ := json.Marshal(&bundle)
	assert.Equal(t, expectedBundle, receivedBundle)

	//assert that GetSettingsBundle returns the same bundle as saved
	getRequest := proto.GetSettingsBundleRequest{Identifier: &identifier}
	getResponse, err := cl.GetSettingsBundle(context.Background(), &getRequest)
	assert.NoError(t, err)
	receivedBundle, _ = json.Marshal(getResponse.SettingsBundle)
	assert.Equal(t, expectedBundle, receivedBundle)

	//assert that ListSettingsBundles returns the same bundle as saved
	listRequest := proto.ListSettingsBundlesRequest{Identifier: &identifier}
	listResponse, err := cl.ListSettingsBundles(context.Background(), &listRequest)
	assert.NoError(t, err)
	receivedBundle, _ = json.Marshal(listResponse.SettingsBundles[0])
	assert.Equal(t, expectedBundle, receivedBundle)

	_ = os.RemoveAll("ocis-settings-store")
}

// https://github.com/owncloud/ocis-settings/issues/18
func TestSaveSettingsBundleWithInvalidSettingValues(t *testing.T) {
	var tests = []proto.Setting{
		{
			SettingKey: "intValue default is out of range",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.IntSetting{
					Default: 30,
					Min:     10,
					Max:     20,
				},
			},
		},
		{
			SettingKey: "intValue min > max",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.IntSetting{
					Default: 100,
					Min:     100,
					Max:     20,
				},
			},
		},
		{
			SettingKey: "intValue step > max-min",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.IntSetting{
					Min:  10,
					Max:  20,
					Step: 100,
				},
			},
		},
		{
			SettingKey: "intValue step = 0",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.IntSetting{
					Min:  10,
					Max:  20,
					Step: 0,
				},
			},
		},
		{
			SettingKey: "intValue step < 0",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.IntSetting{
					Min:  10,
					Max:  20,
					Step: -10,
				},
			},
		},
		{
			SettingKey: "stringValue MinLength > MaxLength",
			Value: &proto.Setting_StringValue{
				StringValue: &proto.StringSetting{
					MinLength: 255,
					MaxLength: 1,
				},
			},
		},
		{
			SettingKey: "stringValue MaxLength = 0",
			Value: &proto.Setting_StringValue{
				StringValue: &proto.StringSetting{
					MaxLength: 0,
				},
			},
		},
		{
			SettingKey: "stringValue MinLength < 0",
			Value: &proto.Setting_StringValue{
				StringValue: &proto.StringSetting{
					MinLength: -1,
				},
			},
		},
		{
			SettingKey: "stringValue MaxLength < 0",
			Value: &proto.Setting_StringValue{
				StringValue: &proto.StringSetting{
					MaxLength: -1,
				},
			},
		},
		{
			SettingKey: "multiChoice multiple options are default",
			Value: &proto.Setting_MultiChoiceValue{
				MultiChoiceValue: &proto.MultiChoiceListSetting{
					Options: []*proto.ListOption{
						{
							Value: &proto.ListOptionValue{
								Option: &proto.ListOptionValue_IntValue{IntValue: 1},
							},
							Default: true,
						},
						{
							Value: &proto.ListOptionValue{
								Option: &proto.ListOptionValue_IntValue{IntValue: 2},
							},
							Default: true,
						},
					},
				},
			},
		},
		{
			SettingKey: "singleChoice multiple options are default",
			Value: &proto.Setting_SingleChoiceValue{
				SingleChoiceValue: &proto.SingleChoiceListSetting{
					Options: []*proto.ListOption{
						{
							Value: &proto.ListOptionValue{
								Option: &proto.ListOptionValue_IntValue{IntValue: 1},
							},
							Default: true,
						},
						{
							Value: &proto.ListOptionValue{
								Option: &proto.ListOptionValue_IntValue{IntValue: 2},
							},
							Default: true,
						},
					},
				},
			},
		},
	}

	identifier := proto.Identifier{
		Extension:   "my-extension",
		BundleKey:   "bundle-with-invalid-settings",
		SettingKey:  "simple-key",
		AccountUuid: "123e4567-d89b-12e3-a656-426652340000",
	}

	for index := range tests {
		index := index
		t.Run(tests[index].SettingKey, func(t *testing.T) {

			var settings []*proto.Setting

			settings = append(settings, &tests[index])

			bundle := proto.SettingsBundle{
				Identifier:  &identifier,
				DisplayName: "bundle display name",
				Settings:    settings,
			}
			saveRequest := proto.SaveSettingsBundleRequest{
				SettingsBundle: &bundle,
			}

			client := service.Client()
			cl := proto.NewBundleService("com.owncloud.api.settings", client)

			//assert that SaveSettingsBundle returns the same bundle as we have sent there
			saveResponse, err := cl.SaveSettingsBundle(context.Background(), &saveRequest)
			assert.NoError(t, err)
			receivedBundle, _ := json.Marshal(saveResponse.SettingsBundle)
			expectedBundle, _ := json.Marshal(&bundle)
			assert.Equal(t, expectedBundle, receivedBundle)
			_ = os.RemoveAll("ocis-settings-store")
		})
	}
}

//https://github.com/owncloud/ocis-settings/issues/19
func TestGetSettingsBundleCreatesFolder(t *testing.T) {
	identifier := proto.Identifier{
		Extension: "not-existing-extension",
		BundleKey: "not-existing-bundle",
	}

	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)
	getRequest := proto.GetSettingsBundleRequest{Identifier: &identifier}

	_, _ = cl.GetSettingsBundle(context.Background(), &getRequest)
	assert.NoDirExists(t, "ocis-settings-store/bundles/not-existing-extension")
	assert.NoFileExists(t, "ocis-settings-store/bundles/not-existing-extension/not-existing-bundle.json")
	_ = os.RemoveAll("ocis-settings-store")
}

func TestGetSettingsBundleAccessOtherBundle(t *testing.T) {
	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)

	aliceBundle := proto.SettingsBundle{
		Identifier: &proto.Identifier{
			Extension: "alice-extension",
			BundleKey: "alice-bundle",
		},
		DisplayName: "alice settings bundle",
		Settings: dummySettings,
	}
	createRequest := proto.SaveSettingsBundleRequest{
		SettingsBundle: &aliceBundle,
	}
	_, err := cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.NoError(t, err)

	bobIdentifier := proto.Identifier{
		Extension: "../bundles/alice-extension/",
		BundleKey: "alice-bundle",
	}

	getRequest := proto.GetSettingsBundleRequest{Identifier: &bobIdentifier}

	response, err := cl.GetSettingsBundle(context.Background(), &getRequest)
	assert.Error(t, err)
	assert.Nil(t, response)
	var errorData CustomError
	_ = json.Unmarshal([]byte(err.Error()), &errorData)
	assert.Equal(t, "go.micro.client", errorData.ID)
	assert.Equal(t, 500, errorData.Code)
	assert.Equal(t, "extension: must be in a valid format.", errorData.Detail)
	assert.Equal(t, "Internal Server Error", errorData.Status)
	_ = os.RemoveAll("ocis-settings-store")
}

/**
  test read settings bundles with identifiers that should be invalid, e.g. try to read other bundles
*/
func TestGetSettingsBundleWithInvalidIdentifier(t *testing.T) {
	type TestStruct struct {
		testDataName  string
		BundleKey     string
		SettingKey    string
		Extension     string
		UUID          string
		expectedError CustomError
	}

	var tests = []TestStruct{
		{
			"not existing",
			"this-key-should-not-exist",
			"this-key-should-not-exist",
			"this-extension-should-not-exist",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "open ocis-settings-store/bundles/this-extension-should-not-exist/this-key-should-not-exist.json: no such file or directory",
				Status: "Internal Server Error",
			},
		},
		{
			"bundle key in the root directory",
			"/tmp/file",
			"simple-key",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "bundle_key: must be in a valid format.",
				Status: "Internal Server Error",
			},
		},
		{
			"bundle key missing",
			"",
			"simple-bundle-key",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "bundle_key: cannot be blank.",
				Status: "Internal Server Error",
			},
		},
		{
			"extension missing",
			"simple-bundle-key",
			"simple-key",
			"",
			"123e4567-e89b-12d3-a456-426652340000",
			CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "extension: cannot be blank.",
				Status: "Internal Server Error",
			},
		},
	}
	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.testDataName, func(t *testing.T) {
			identifier := proto.Identifier{
				Extension:   testCase.Extension,
				BundleKey:   testCase.BundleKey,
				SettingKey:  testCase.SettingKey,
				AccountUuid: testCase.UUID,
			}

			client := service.Client()
			cl := proto.NewBundleService("com.owncloud.api.settings", client)
			getRequest := proto.GetSettingsBundleRequest{Identifier: &identifier}

			getResponse, err := cl.GetSettingsBundle(context.Background(), &getRequest)
			if err != nil || (CustomError{} != testCase.expectedError) {
				var errorData CustomError
				assert.Error(t, err)
				assert.Empty(t, getResponse)
				_ = json.Unmarshal([]byte(err.Error()), &errorData)
				assert.Equal(t, testCase.expectedError.ID, errorData.ID)
				assert.Equal(t, testCase.expectedError.Code, errorData.Code)
				assert.Equal(t, testCase.expectedError.Detail, errorData.Detail)
				assert.Equal(t, testCase.expectedError.Status, errorData.Status)
			} else {
				assert.NoError(t, err)
			}
			_ = os.RemoveAll("ocis-settings-store")
		})
	}
}

func TestListMultipleSettingsBundlesOfSameExtension(t *testing.T) {
	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)

	createRequest := proto.SaveSettingsBundleRequest{
		SettingsBundle: &proto.SettingsBundle{
			Identifier: &proto.Identifier{
				Extension: "great-extension",
				BundleKey: "alices-bundle",
			},
			DisplayName: "Alice's Bundle",
			Settings:    dummySettings,
		},
	}
	_, err := cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.NoError(t, err)

	createRequest = proto.SaveSettingsBundleRequest{
		SettingsBundle: &proto.SettingsBundle{
			Identifier: &proto.Identifier{
				Extension: "great-extension",
				BundleKey: "bobs-bundle",
			},
			DisplayName: "Bob's Bundle",
			Settings:    dummySettings,
		},
	}
	_, err = cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.NoError(t, err)

	createRequest = proto.SaveSettingsBundleRequest{
		SettingsBundle: &proto.SettingsBundle{
			Identifier: &proto.Identifier{
				Extension: "an-other-extension",
				BundleKey: "bobs-bundle",
			},
			DisplayName: "Bob's Bundle",
			Settings:    dummySettings,
		},
	}
	_, err = cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.NoError(t, err)

	listRequest := proto.ListSettingsBundlesRequest{Identifier: &proto.Identifier{Extension: "great-extension"}}

	response, err := cl.ListSettingsBundles(context.Background(), &listRequest)
	assert.NoError(t, err)
	assert.Equal(t, response.SettingsBundles[0].Identifier.Extension, "great-extension")
	assert.Equal(t, response.SettingsBundles[0].Identifier.BundleKey, "alices-bundle")

	assert.Equal(t, response.SettingsBundles[1].Identifier.Extension, "great-extension")
	assert.Equal(t, response.SettingsBundles[1].Identifier.BundleKey, "bobs-bundle")
	assert.Equal(t, 2, len(response.SettingsBundles))
	_ = os.RemoveAll("ocis-settings-store")
}

func TestListAllSettingsBundlesOfSameExtension(t *testing.T) {
	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)

	createRequest := proto.SaveSettingsBundleRequest{
		SettingsBundle: &proto.SettingsBundle{
			Identifier: &proto.Identifier{
				Extension: "great-extension",
				BundleKey: "alices-bundle",
			},
			DisplayName: "Alice's Bundle",
			Settings:    dummySettings,
		},
	}
	_, err := cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.NoError(t, err)

	createRequest = proto.SaveSettingsBundleRequest{
		SettingsBundle: &proto.SettingsBundle{
			Identifier: &proto.Identifier{
				Extension: "great-extension",
				BundleKey: "bobs-bundle",
			},
			DisplayName: "Bob's Bundle",
			Settings:    dummySettings,
		},
	}
	_, err = cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.NoError(t, err)

	createRequest = proto.SaveSettingsBundleRequest{
		SettingsBundle: &proto.SettingsBundle{
			Identifier: &proto.Identifier{
				Extension: "an-other-extension",
				BundleKey: "bobs-bundle",
			},
			DisplayName: "Bob's Bundle",
			Settings:    dummySettings,
		},
	}
	_, err = cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.NoError(t, err)

	listRequest := proto.ListSettingsBundlesRequest{Identifier: &proto.Identifier{Extension: ""}}

	response, err := cl.ListSettingsBundles(context.Background(), &listRequest)
	assert.NoError(t, err)
	assert.Equal(t, response.SettingsBundles[0].Identifier.Extension, "an-other-extension")
	assert.Equal(t, response.SettingsBundles[0].Identifier.BundleKey, "bobs-bundle")

	assert.Equal(t, response.SettingsBundles[1].Identifier.Extension, "great-extension")
	assert.Equal(t, response.SettingsBundles[1].Identifier.BundleKey, "alices-bundle")

	assert.Equal(t, response.SettingsBundles[2].Identifier.Extension, "great-extension")
	assert.Equal(t, response.SettingsBundles[2].Identifier.BundleKey, "bobs-bundle")
	assert.Equal(t, 3, len(response.SettingsBundles))
	_ = os.RemoveAll("ocis-settings-store")
}

func TestListSettingsBundlesOfNonExistingExtension(t *testing.T) {
	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)

	listRequest := proto.ListSettingsBundlesRequest{Identifier: &proto.Identifier{Extension: "does-not-exist"}}

	response, err := cl.ListSettingsBundles(context.Background(), &listRequest)
	assert.NoError(t, err)
	assert.Empty(t, response.String())
	assert.NoDirExists(t, "ocis-settings-store/bundles")
	assert.NoDirExists(t, "ocis-settings-store/bundles/does-not-exist")
}

func TestListSettingsBundlesInFoldersThatAreNotAccessible(t *testing.T) {
	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)
	createRequest := proto.SaveSettingsBundleRequest{
		SettingsBundle: &proto.SettingsBundle{
			Identifier: &proto.Identifier{
				Extension: "great-extension",
				BundleKey: "alices-bundle",
			},
			DisplayName: "Alice's Bundle",
			Settings:    dummySettings,
		},
	}
	_, err := cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.NoError(t, err)

	createRequest = proto.SaveSettingsBundleRequest{
		SettingsBundle: &proto.SettingsBundle{
			Identifier: &proto.Identifier{
				Extension: "great-extension",
				BundleKey: "bobs-bundle",
			},
			DisplayName: "Bob's Bundle",
			Settings:    dummySettings,
		},
	}
	_, err = cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.NoError(t, err)

	createRequest = proto.SaveSettingsBundleRequest{
		SettingsBundle: &proto.SettingsBundle{
			Identifier: &proto.Identifier{
				Extension: "an-other-extension",
				BundleKey: "bobs-bundle",
			},
			DisplayName: "Bob's Bundle",
			Settings:    dummySettings,
		},
	}
	_, err = cl.SaveSettingsBundle(context.Background(), &createRequest)
	assert.NoError(t, err)

	listRequest := proto.ListSettingsBundlesRequest{Identifier: &proto.Identifier{Extension: "../"}}

	response, err := cl.ListSettingsBundles(context.Background(), &listRequest)
	assert.Error(t, err)
	assert.Nil(t, response)
	var errorData CustomError
	_ = json.Unmarshal([]byte(err.Error()), &errorData)
	assert.Equal(t, "go.micro.client", errorData.ID)
	assert.Equal(t, 500, errorData.Code)
	assert.Equal(t, "extension: must be in a valid format.", errorData.Detail)
	assert.Equal(t, "Internal Server Error", errorData.Status)
	_ = os.RemoveAll("ocis-settings-store")
}

func TestSaveGetListSettingsValues(t *testing.T) {
	type TestStruct struct {
		testDataName  string
		SettingsValue proto.SettingsValue
		expectedError CustomError
	}

	var tests = []TestStruct{
		{
			testDataName: "simple int",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "alices-bundle",
					AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "age",
				},
				Value: &proto.SettingsValue_IntValue{IntValue: 12},
			},
		},
		{
			testDataName: "simple string",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "alices-bundle",
					AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "location",
				},
				Value: &proto.SettingsValue_StringValue{StringValue: "पोखरा"},
			},
		},
		{
			testDataName: "simple bool",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "alices-bundle",
					AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "locked",
				},
				Value: &proto.SettingsValue_BoolValue{BoolValue: false},
			},
		},
		{
			testDataName: "string list",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "alices-bundle",
					AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "currencies",
				},
				Value: &proto.SettingsValue_ListValue{
					ListValue: &proto.ListValue{
						Values: []*proto.ListOptionValue{
							{Option: &proto.ListOptionValue_StringValue{StringValue: "NPR"}},
							{Option: &proto.ListOptionValue_StringValue{StringValue: "EUR"}},
							{Option: &proto.ListOptionValue_StringValue{StringValue: "USD"}},
						},
					},
				},
			},
		},
		{
			testDataName: "int list",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "alices-bundle",
					AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "font-size",
				},
				Value: &proto.SettingsValue_ListValue{
					ListValue: &proto.ListValue{
						Values: []*proto.ListOptionValue{
							{Option: &proto.ListOptionValue_IntValue{IntValue: 11}},
							{Option: &proto.ListOptionValue_IntValue{IntValue: 12}},
							{Option: &proto.ListOptionValue_IntValue{IntValue: 13}},
						},
					},
				},
			},
		},
		{
			//https://github.com/owncloud/ocis-settings/issues/20
			testDataName: "mixed type list",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "alices-bundle",
					AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "apple-and-peaches",
				},
				Value: &proto.SettingsValue_ListValue{
					ListValue: &proto.ListValue{
						Values: []*proto.ListOptionValue{
							{Option: &proto.ListOptionValue_StringValue{StringValue: "string"}},
							{Option: &proto.ListOptionValue_IntValue{IntValue: 123}},
						},
					},
				},
			},
		},
		{
			testDataName: "extension name missing",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "",
					BundleKey:   "alices-bundle",
					AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "locked",
				},
				Value: &proto.SettingsValue_BoolValue{BoolValue: false},
			},
			expectedError: CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "extension: cannot be blank.",
				Status: "Internal Server Error",
			},
		},
		{
			testDataName: "bundle key missing",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "",
					AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "locked",
				},
				Value: &proto.SettingsValue_BoolValue{BoolValue: false},
			},
			expectedError: CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "bundle_key: cannot be blank.",
				Status: "Internal Server Error",
			},
		},
		{
			testDataName: "account uuid missing",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "bobs-bundle",
					AccountUuid: "",
					SettingKey:  "locked",
				},
				Value: &proto.SettingsValue_BoolValue{BoolValue: false},
			},
			expectedError: CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "account_uuid: cannot be blank.",
				Status: "Internal Server Error",
			},
		},
		{
			testDataName: "settings key missing",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "bobs-bundle",
					AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "",
				},
				Value: &proto.SettingsValue_BoolValue{BoolValue: false},
			},
			expectedError: CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "setting_key: cannot be blank.",
				Status: "Internal Server Error",
			},
		},
		{
			testDataName: "../ in bundle key",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "../bobs-bundle",
					AccountUuid: "123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "should-not-be-possible",
				},
				Value: &proto.SettingsValue_BoolValue{BoolValue: false},
			},
			expectedError: CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "bundle_key: must be in a valid format.",
				Status: "Internal Server Error",
			},
		},
		{
			testDataName: "../ in account uuid",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "great-extension",
					BundleKey:   "bobs-bundle",
					AccountUuid: "../123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "should-not-be-possible",
				},
				Value: &proto.SettingsValue_BoolValue{BoolValue: false},
			},
			expectedError: CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "account_uuid: must be a valid UUID.",
				Status: "Internal Server Error",
			},
		},
		{
			testDataName: "\\ in fields that are used to create folder and file names",
			SettingsValue: proto.SettingsValue{
				Identifier: &proto.Identifier{
					Extension:   "\\-extension",
					BundleKey:   "\\-bundle",
					AccountUuid: "\\123e4567-e89b-12d3-a456-426652340000",
					SettingKey:  "should-not-be-possible",
				},
				Value: &proto.SettingsValue_BoolValue{BoolValue: false},
			},
			expectedError: CustomError{
				ID:     "go.micro.client",
				Code:   500,
				Detail: "account_uuid: must be a valid UUID; bundle_key: must be in a valid format; extension: must be in a valid format.",
				Status: "Internal Server Error",
			},
		},
	}
	client := service.Client()
	cl := proto.NewValueService("com.owncloud.api.settings", client)

	for index := range tests {
		index := index
		t.Run(tests[index].testDataName, func(t *testing.T) {
			createRequest := proto.SaveSettingsValueRequest{
				SettingsValue: &tests[index].SettingsValue,
			}
			saveResponse, err := cl.SaveSettingsValue(context.Background(), &createRequest)
			if err != nil || (CustomError{} != tests[index].expectedError) {
				assert.Error(t, err)
				var errorData CustomError
				_ = json.Unmarshal([]byte(err.Error()), &errorData)
				assert.Equal(t, tests[index].expectedError.ID, errorData.ID)
				assert.Equal(t, tests[index].expectedError.Code, errorData.Code)
				assert.Equal(t, tests[index].expectedError.Detail, errorData.Detail)
				assert.Equal(t, tests[index].expectedError.Status, errorData.Status)
			} else {
				expectedSetting, _ := json.Marshal(&tests[index].SettingsValue)

				assert.NoError(t, err)
				receivedSetting, _ := json.Marshal(saveResponse.SettingsValue)
				assert.Equal(t, expectedSetting, receivedSetting)

				getRequest := proto.GetSettingsValueRequest{
					Identifier: tests[index].SettingsValue.Identifier,
				}
				getResponse, err := cl.GetSettingsValue(context.Background(), &getRequest)
				assert.NoError(t, err)
				receivedSetting, _ = json.Marshal(getResponse.SettingsValue)
				assert.Equal(t, expectedSetting, receivedSetting)

				listRequest := proto.ListSettingsValuesRequest{
					Identifier: tests[index].SettingsValue.Identifier,
				}
				listResponse, err := cl.ListSettingsValues(context.Background(), &listRequest)
				assert.NoError(t, err)
				receivedSetting, _ = json.Marshal(listResponse.SettingsValues[0])
				assert.Equal(t, expectedSetting, receivedSetting)
			}

			_ = os.RemoveAll("ocis-settings-store")
		})
	}
}
