package proto_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	ocislog "github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-pkg/v2/middleware"
	"github.com/owncloud/ocis-pkg/v2/service/grpc"
	"github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	svc "github.com/owncloud/ocis-settings/pkg/service/v0"
	store "github.com/owncloud/ocis-settings/pkg/store/filesystem"
	"github.com/stretchr/testify/assert"
)

var (
	service grpc.Service
	handler svc.Service
	bundleService proto.BundleService
	valueService proto.ValueService
	roleService proto.RoleService
	permissionService proto.PermissionService

	testAccountID = "e8a7f56b-10ce-4f67-b67f-eca40aa0ef26"

	settingsStub = []*proto.Setting{
		{
			Id:          "336c4db1-5062-4931-990f-d88e6b02cb02",
			DisplayName: "dummy setting",
			Name:        "dummy-setting",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Value: &proto.Setting_IntValue{
				IntValue: &proto.Int{
					Default: 42,
				},
			},
			Description: "dummy setting",
		},
	}

	//optionsForListStub for choice list settings
	optionsForListStub = []*proto.ListOption{
		{
			Value: &proto.ListOptionValue{
				Option: &proto.ListOptionValue_StringValue{StringValue: "list option string value"},
			},
			Default:      true,
			DisplayValue: "a string value",
		},
		{
			Value: &proto.ListOptionValue{
				Option: &proto.ListOptionValue_IntValue{IntValue: 123},
			},
			Default:      true,
			DisplayValue: "a int value",
		},
	}

	//MultiChoiceList
	multipleChoiceSettingStub = proto.MultiChoiceList{
		Options: optionsForListStub,
	}

	//SingleChoiceList
	singleChoiceSettingStub = proto.SingleChoiceList{
		Options: optionsForListStub,
	}

	complexSettingsStub = []*proto.Setting{
		{
			Name:        "int",
			Id:          "4e00633d-5373-4df4-9299-1c9ed9c3ebed",
			DisplayName: "an integer value",
			Description: "with some description",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.Int{
					Default:     1,
					Min:         1,
					Max:         124,
					Step:        1,
					Placeholder: "Int value",
				},
			},
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
		},
		{
			Name:        "string",
			Id:          "f792acb4-9f09-4fa8-92d3-4a0d0a6ca721",
			DisplayName: "a string value",
			Description: "with some description",
			Value: &proto.Setting_StringValue{
				StringValue: &proto.String{
					Default:     "thedefaultvalue",
					Required:    false,
					MinLength:   2,
					MaxLength:   255,
					Placeholder: "a string value",
				},
			},
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
		},
		{
			Name:        "bool",
			Id:          "6ef9268c-f0bd-48a7-a0a0-3ba3ee42b5cc",
			DisplayName: "a bool value",
			Description: "with some description",
			Value: &proto.Setting_BoolValue{
				BoolValue: &proto.Bool{
					Default: false,
					Label:   "bool setting",
				},
			},
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
		},
		{
			Name:        "multipleChoice",
			Id:          "905da88c-3be0-42c2-a8b2-e6bcd9976b2d",
			DisplayName: "a multiple choice setting",
			Description: "with some description",
			Value: &proto.Setting_MultiChoiceValue{
				MultiChoiceValue: &multipleChoiceSettingStub,
			},
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
		},
		{
			Name:        "singleChoice",
			Id:          "5bf4de47-57cc-4705-a456-4fcf39673994",
			DisplayName: "a single choice setting",
			Description: "with some description",
			Value: &proto.Setting_SingleChoiceValue{
				SingleChoiceValue: &singleChoiceSettingStub,
			},
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
		},
	}

	bundleStub = proto.Bundle{
		Name:        "test",
		Id:          "b1b8c9d0-fb3c-4e12-b868-5a8508218d2e",
		DisplayName: "bundleDisplayName",
		Extension:   "testExtension",
		Type:        proto.Bundle_TYPE_DEFAULT,
		Settings:    complexSettingsStub,
		Resource: &proto.Resource{
			Type: proto.Resource_TYPE_SYSTEM,
		},
	}
)

const dataPath = "/var/tmp/grpc-tests-ocis-settings"

func init() {
	service = grpc.NewService(
		grpc.Namespace("com.owncloud.api"),
		grpc.Name("settings"),
		grpc.Address("localhost:9992"),
	)

	cfg := config.New()
	cfg.Storage.DataPath = dataPath
	handler = svc.NewService(cfg, ocislog.NewLogger(ocislog.Color(true), ocislog.Pretty(true)))
	err := proto.RegisterBundleServiceHandler(service.Server(), handler)
	if err != nil {
		log.Fatalf("could not register BundleServiceHandler: %v", err)
	}
	err = proto.RegisterValueServiceHandler(service.Server(), handler)
	if err != nil {
		log.Fatalf("could not register ValueServiceHandler: %v", err)
	}
	err = proto.RegisterRoleServiceHandler(service.Server(), handler)
	if err != nil {
		log.Fatalf("could not register RoleServiceHandler: %v", err)
	}
	err = proto.RegisterPermissionServiceHandler(service.Server(), handler)
	if err != nil {
		log.Fatalf("could not register PermissionServiceHandler: %v", err)
	}

	if err = service.Server().Start(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	client := service.Client()
	bundleService = proto.NewBundleService("com.owncloud.api.settings", client)
	valueService = proto.NewValueService("com.owncloud.api.settings", client)
	roleService = proto.NewRoleService("com.owncloud.api.settings", client)
	permissionService = proto.NewPermissionService("com.owncloud.api.settings", client)
}

func setup() func() {
	handler.RegisterDefaultRoles()
	return func() {
		if err := os.RemoveAll(dataPath); err != nil {
			log.Printf("could not delete data root: %s", dataPath)
		} else {
			log.Println("data root deleted")
		}
	}
}

/**
testing that saving a settings bundle and retrieving it again works correctly
using various setting bundle properties
*/
func TestBundleInputValidation(t *testing.T) {
	var scenarios = []struct {
		name          string
		bundleName    string
		displayName   string
		extensionName string
		expectedError error
	}{
		{
			"ASCII",
			"bundle-name",
			"simple-bundle-key",
			"simple-extension-name",
			nil,
		},
		{
			"UTF validation on bundle name",
			"सिम्प्ले-bundle-name",
			"सिम्प्ले-display-name",
			"सिम्प्ले-extension-name",
			merrors.New("ocis-settings", "extension: must be in a valid format; name: must be in a valid format.", 400),
		},
		{
			"UTF validation on display name",
			"सिम्प्ले-bundle-name",
			"सिम्प्ले-display-name",
			"simple-extension-name",
			merrors.New("ocis-settings", "name: must be in a valid format.", 400),
		},
		{
			"extension name with ../ in the name",
			"bundle-name",
			"simple-display-name",
			"../folder-a-level-higher-up",
			merrors.New("ocis-settings", "extension: must be in a valid format.", 400),
		},
		{
			"extension name with \\ in the name",
			"bundle-name",
			"simple-display-name",
			"\\",
			merrors.New("ocis-settings", "extension: must be in a valid format.", 400),
		},
		{
			"spaces are disallowed in bundle names",
			"bundle name",
			"simple display name",
			"simple extension name",
			merrors.New("ocis-settings", "extension: must be in a valid format; name: must be in a valid format.", 400),
		},
		{
			"spaces are allowed in display names",
			"bundle-name",
			"simple display name",
			"simple-extension-name",
			nil,
		},
		{
			"extension missing",
			"bundle-name",
			"simple-display-name",
			"",
			merrors.New("ocis-settings", "extension: cannot be blank.", 400),
		},
		{
			"display name missing",
			"bundleName",
			"",
			"simple-extension-name",
			merrors.New("ocis-settings", "display_name: cannot be blank.", 400),
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			bundle := proto.Bundle{
				Name:        scenario.bundleName,
				Extension:   scenario.extensionName,
				DisplayName: scenario.displayName,
				Type:        proto.Bundle_TYPE_DEFAULT,
				Resource: &proto.Resource{
					Type: proto.Resource_TYPE_SYSTEM,
				},
				Settings: settingsStub,
			}
			createRequest := proto.SaveBundleRequest{
				Bundle: &bundle,
			}

			cresponse, err := bundleService.SaveBundle(context.Background(), &createRequest)
			if err != nil || scenario.expectedError != nil {
				t.Log(err)
				assert.Equal(t, scenario.expectedError, err)
			} else {
				assert.Equal(t, scenario.extensionName, cresponse.Bundle.Extension)
				assert.Equal(t, scenario.displayName, cresponse.Bundle.DisplayName)

				// we want to test input validation, so just allow the request permission-wise
				setFullReadWriteOnBundle(t, testAccountID, cresponse.Bundle.Id)

				ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
				getRequest := proto.GetBundleRequest{BundleId: cresponse.Bundle.Id}
				getResponse, err := bundleService.GetBundle(ctx, &getRequest)
				assert.NoError(t, err)
				if err == nil {
					assert.Equal(t, scenario.displayName, getResponse.Bundle.DisplayName)
				}
			}
		})
	}
}

func TestSaveBundleWithoutSettings(t *testing.T) {
	teardown := setup()
	defer teardown()

	createRequest := proto.SaveBundleRequest{
		Bundle: &proto.Bundle{
			DisplayName: "Alice's Bundle",
		},
	}
	response, err := bundleService.SaveBundle(context.Background(), &createRequest)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, merrors.New("ocis-settings", "extension: cannot be blank; name: cannot be blank; settings: cannot be blank.", 400), err)
}

/**
testing that setting getting and listing a settings bundle works correctly with a set of setting definitions
*/
func TestSaveAndGetBundle(t *testing.T) {
	teardown := setup()
	defer teardown()

	saveRequest := proto.SaveBundleRequest{
		Bundle: &bundleStub,
	}

	// assert that SaveBundle returns the same bundle as we have sent
	saveResponse, err := bundleService.SaveBundle(context.Background(), &saveRequest)
	assert.NoError(t, err)
	receivedBundle, _ := json.Marshal(saveResponse.Bundle.Settings)
	expectedBundle, _ := json.Marshal(&bundleStub.Settings)
	assert.Equal(t, receivedBundle, expectedBundle)

	// set full permissions for getting the created bundle
	setFullReadWriteOnBundle(t, testAccountID, saveResponse.Bundle.Id)

	//assert that GetBundle returns the same bundle as saved
	getRequest := proto.GetBundleRequest{BundleId: saveResponse.Bundle.Id}
	ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
	getResponse, err := bundleService.GetBundle(ctx, &getRequest)
	assert.NoError(t, err)
	if err == nil {
		receivedBundle, _ = json.Marshal(getResponse.Bundle.Settings)
		assert.Equal(t, expectedBundle, receivedBundle)
	}
}

/**
testing that saving a value works and can be retrieved again
*/
func TestSaveGetIntValue(t *testing.T) {
	tests := []struct {
		name  string
		value proto.Value_IntValue
	}{
		{
			name:  "simple int",
			value: proto.Value_IntValue{IntValue: 43},
		},
		{
			name:  "negative",
			value: proto.Value_IntValue{IntValue: -42},
			// https://github.com/owncloud/ocis-settings/issues/57
		},
		{
			name:  "less than Min",
			value: proto.Value_IntValue{IntValue: 0},
			// https://github.com/owncloud/ocis-settings/issues/57
		},
		{
			name:  "more than Max",
			value: proto.Value_IntValue{IntValue: 128},
			// https://github.com/owncloud/ocis-settings/issues/57
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			saveResponse, err := bundleService.SaveBundle(context.Background(), &proto.SaveBundleRequest{
				Bundle: &bundleStub,
			})
			assert.NoError(t, err)

			saveValueResponse, err := valueService.SaveValue(context.Background(), &proto.SaveValueRequest{
				Value: &proto.Value{
					BundleId:    saveResponse.Bundle.Id,
					SettingId:   "4e00633d-5373-4df4-9299-1c9ed9c3ebed", //setting id of the int setting
					AccountUuid: "047d31b0-219a-47a4-8ee5-c5fa3802a3c2",
					Value:       &tt.value,
					Resource: &proto.Resource{
						Type: 0,
						Id:   ".",
					},
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, tt.value.IntValue, saveValueResponse.Value.Value.GetIntValue())

			getValueResponse, err := valueService.GetValue(
				context.Background(), &proto.GetValueRequest{Id: saveValueResponse.Value.Value.Id},
			)
			assert.NoError(t, err)
			assert.Equal(t, tt.value.IntValue, getValueResponse.Value.Value.GetIntValue())
		})
	}
}

/**
try to save a wrong type of the value
https://github.com/owncloud/ocis-settings/issues/57
*/
func TestSaveGetIntValueIntoString(t *testing.T) {
	teardown := setup()
	defer teardown()

	saveResponse, err := bundleService.SaveBundle(context.Background(), &proto.SaveBundleRequest{
		Bundle: &bundleStub,
	})
	assert.NoError(t, err)

	saveValueResponse, err := valueService.SaveValue(context.Background(), &proto.SaveValueRequest{
		Value: &proto.Value{
			BundleId:    saveResponse.Bundle.Id,
			SettingId:   "f792acb4-9f09-4fa8-92d3-4a0d0a6ca721", //setting id of the string setting
			AccountUuid: "047d31b0-219a-47a4-8ee5-c5fa3802a3c2",
			Value:       &proto.Value_StringValue{StringValue: "forty two"},
			Resource: &proto.Resource{
				Type: 0,
				Id:   ".",
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "forty two", saveValueResponse.Value.Value.GetStringValue())

	getValueResponse, err := valueService.GetValue(
		context.Background(), &proto.GetValueRequest{Id: saveValueResponse.Value.Value.Id},
	)
	assert.NoError(t, err)
	assert.Equal(t, "forty two", getValueResponse.Value.Value.GetStringValue())
}

// https://github.com/owncloud/ocis-settings/issues/18
func TestSaveBundleWithInvalidSettings(t *testing.T) {
	var tests = []proto.Setting{
		{
			Name: "intValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "intValue default is out of range",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.Int{
					Default: 30,
					Min:     10,
					Max:     20,
				},
			},
		},
		{
			Name: "intValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "intValue min gt max",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.Int{
					Default: 100,
					Min:     100,
					Max:     20,
				},
			},
		},
		{
			Name: "intValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "intValue step gt max-min",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.Int{
					Min:  10,
					Max:  20,
					Step: 100,
				},
			},
		},
		{
			Name: "intValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "intValue step eq 0",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.Int{
					Min:  10,
					Max:  20,
					Step: 0,
				},
			},
		},
		{
			Name: "intValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "intValue step lt 0",
			Value: &proto.Setting_IntValue{
				IntValue: &proto.Int{
					Min:  10,
					Max:  20,
					Step: -10,
				},
			},
		},
		{
			Name: "stringValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "stringValue MinLength gt MaxLength",
			Value: &proto.Setting_StringValue{
				StringValue: &proto.String{
					MinLength: 255,
					MaxLength: 1,
				},
			},
		},
		{
			Name: "stringValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "stringValue MaxLength eq 0",
			Value: &proto.Setting_StringValue{
				StringValue: &proto.String{
					MaxLength: 0,
				},
			},
		},
		{
			Name: "stringValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "stringValue MinLength lt 0",
			Value: &proto.Setting_StringValue{
				StringValue: &proto.String{
					MinLength: -1,
				},
			},
		},
		{
			Name: "stringValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "stringValue MaxLength lt 0",
			Value: &proto.Setting_StringValue{
				StringValue: &proto.String{
					MaxLength: -1,
				},
			},
		},
		{
			Name: "multiChoiceValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "multiChoice multiple options are default",
			Value: &proto.Setting_MultiChoiceValue{
				MultiChoiceValue: &proto.MultiChoiceList{
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
			Name: "singleChoiceValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_USER,
			},
			Description: "singleChoice multiple options are default",
			Value: &proto.Setting_SingleChoiceValue{
				SingleChoiceValue: &proto.SingleChoiceList{
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

	for index := range tests {
		t.Run(tests[index].Name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			var settings []*proto.Setting

			settings = append(settings, &tests[index])
			bundle := proto.Bundle{
				Name:        "bundle",
				Extension:   "bundleExtension",
				DisplayName: "bundledisplayname",
				Type:        proto.Bundle_TYPE_DEFAULT,
				Resource: &proto.Resource{
					Type: proto.Resource_TYPE_SYSTEM,
				},
				Settings: settings,
			}
			saveRequest := proto.SaveBundleRequest{
				Bundle: &bundle,
			}

			//assert that SaveBundle returns the same bundle as we have sent there
			saveResponse, err := bundleService.SaveBundle(context.Background(), &saveRequest)
			assert.NoError(t, err)
			receivedBundle, _ := json.Marshal(saveResponse.Bundle.Settings)
			expectedBundle, _ := json.Marshal(&bundle.Settings)
			assert.Equal(t, expectedBundle, receivedBundle)
		})
	}
}

// https://github.com/owncloud/ocis-settings/issues/19
func TestGetBundleNoSideEffectsOnDisk(t *testing.T) {
	teardown := setup()
	defer teardown()

	getRequest := proto.GetBundleRequest{BundleId: "non-existing-bundle"}

	_, _ = bundleService.GetBundle(context.Background(), &getRequest)
	assert.NoDirExists(t, store.Name+"/bundles/non-existing-bundle")
	assert.NoFileExists(t, store.Name+"/bundles/non-existing-bundle/not-existing-bundle.json")
}

// TODO non-deterministic. Fix.
func TestCreateRoleAndAssign(t *testing.T) {
	teardown := setup()
	defer teardown()

	res, err := bundleService.SaveBundle(context.Background(), &proto.SaveBundleRequest{
		Bundle: &proto.Bundle{
			Type:        proto.Bundle_TYPE_ROLE,
			DisplayName: "test role - update",
			Name:        "TEST_ROLE",
			Extension:   "ocis-settings",
			Settings: []*proto.Setting{
				{
					Name: "settingName",
					Resource: &proto.Resource{
						Id:   settingsStub[0].Id,
						Type: proto.Resource_TYPE_SETTING,
					},
					Value: &proto.Setting_PermissionValue{
						PermissionValue: &proto.Permission{
							Operation:  proto.Permission_OPERATION_UPDATE,
							Constraint: proto.Permission_CONSTRAINT_OWN,
						},
					},
				},
			},
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_SYSTEM,
			},
		},
	})
	if err == nil {
		_, err = roleService.AssignRoleToUser(context.Background(), &proto.AssignRoleToUserRequest{
			AccountUuid: "4c510ada-c86b-4815-8820-42cdf82c3d51",
			RoleId:      res.Bundle.Id,
		})
		if err != nil {
			fmt.Println(err)
		}
		assert.NoError(t, err)
	}
}

// // TODO this tests are non-deterministic at least on my machine. Find a way to make them deterministic.
// func TestListBudlesOnAuthorizedUser(t *testing.T) {
// 	client := service.Client()
// 	client2 := service.Client()
// 	cl := proto.NewBundleService("com.owncloud.api.settings", client)
// 	rc := proto.NewRoleService("com.owncloud.api.settings", client2)

// 	_, err := cl.SaveBundle(context.Background(), &proto.SaveBundleRequest{
// 		Bundle: &proto.Bundle{
// 			DisplayName: "Alice's Bundle",
// 			Name:        "bundle1",
// 			Extension:   "extension1",
// 			Resource: &proto.Resource{
// 				Type: proto.Resource_TYPE_BUNDLE,
// 			},
// 			Type:     proto.Bundle_TYPE_DEFAULT,
// 			Settings: settingsStub,
// 		},
// 	})
// 	assert.NoError(t, err)

// res, err := cl.SaveBundle(context.Background(), &proto.SaveBundleRequest{
// 	Bundle: &proto.Bundle{
// 		// Id:          "f36db5e6-a03c-40df-8413-711c67e40b47", // bug: providing the ID ignores its value for the filename.
// 		Type:        proto.Bundle_TYPE_ROLE,
// 		DisplayName: "test role - update",
// 		Name:        "TEST_ROLE",
// 		Extension:   "ocis-settings",
// 		Settings: []*proto.Setting{
// 			{
// 				Name: "settingName",
// 				Resource: &proto.Resource{
// 					Id:   settingsStub[0].Id,
// 					Type: proto.Resource_TYPE_SETTING,
// 				},
// 				Value: &proto.Setting_PermissionValue{
// 					&proto.Permission{
// 						Operation:  proto.PermissionSetting_OPERATION_UPDATE,
// 						Constraint: proto.PermissionSetting_CONSTRAINT_OWN,
// 					},
// 				},
// 			},
// 		},
// 		Resource: &proto.Resource{
// 			Type: proto.Resource_TYPE_SYSTEM,
// 		},
// 	},
// })
// assert.NoError(t, err)

// _, err = rc.AssignRoleToUser(context.Background(), &proto.AssignRoleToUserRequest{
// 	AccountUuid: "4c510ada-c86b-4815-8820-42cdf82c3d51",
// 	RoleId:      res.Bundle.Id,
// })
// assert.NoError(t, err)

// 	time.Sleep(200 * time.Millisecond)
// 	listRequest := proto.ListSettingsBundlesRequest{AccountUuid: "4c510ada-c86b-4815-8820-42cdf82c3d51"}

// 	response, err := cl.ListSettingsBundles(context.Background(), &listRequest)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 1, len(response.Bundles))
// 	assert.Equal(t, response.Bundles[0].Name, "bundle1")
//

func TestListRolesAfterSavingBundle(t *testing.T) {
	type expectedBundle struct {
		displayName string
		name        string
	}

	tests := []struct {
		name            string
		bundles         []*proto.Bundle
		expectedBundles []expectedBundle
	}{
		{"don't add bundle",
			[]*proto.Bundle{},
			[]expectedBundle{
				{displayName: "Guest", name: "guest"},
				{displayName: "Admin", name: "admin"},
				{displayName: "User", name: "user"},
			},
		},
		{name: "one bundle",
			bundles: []*proto.Bundle{{
				Type:        proto.Bundle_TYPE_ROLE,
				DisplayName: "test role - update",
				Name:        "TEST_ROLE",
				Extension:   "ocis-settings",
				Settings: []*proto.Setting{
					{
						Name: "settingName",
						Resource: &proto.Resource{
							Id:   settingsStub[0].Id,
							Type: proto.Resource_TYPE_SETTING,
						},
						Value: &proto.Setting_PermissionValue{
							PermissionValue: &proto.Permission{
								Operation:  proto.Permission_OPERATION_UPDATE,
								Constraint: proto.Permission_CONSTRAINT_OWN,
							},
						},
					},
				},
				Resource: &proto.Resource{
					Type: proto.Resource_TYPE_SYSTEM,
				},
			}},
			expectedBundles: []expectedBundle{
				{displayName: "test role - update", name: "TEST_ROLE"},
				{displayName: "Guest", name: "guest"},
				{displayName: "Admin", name: "admin"},
				{displayName: "User", name: "user"},
			},
		},
		{name: "two added bundles",
			bundles: []*proto.Bundle{{
				Type:        proto.Bundle_TYPE_ROLE,
				DisplayName: "test role - update",
				Name:        "TEST_ROLE",
				Extension:   "ocis-settings",
				Settings: []*proto.Setting{
					{
						Name: "settingName",
						Resource: &proto.Resource{
							Id:   settingsStub[0].Id,
							Type: proto.Resource_TYPE_SETTING,
						},
						Value: &proto.Setting_PermissionValue{
							PermissionValue: &proto.Permission{
								Operation:  proto.Permission_OPERATION_UPDATE,
								Constraint: proto.Permission_CONSTRAINT_OWN,
							},
						},
					},
				},
				Resource: &proto.Resource{
					Type: proto.Resource_TYPE_SYSTEM,
				},
			},
				{
					Type:        proto.Bundle_TYPE_ROLE,
					DisplayName: "an other role",
					Name:        "AnOtherROLE",
					Extension:   "ocis-settings",
					Settings: []*proto.Setting{
						{
							Name: "settingName",
							Resource: &proto.Resource{
								Id:   settingsStub[0].Id,
								Type: proto.Resource_TYPE_SETTING,
							},
							Value: &proto.Setting_PermissionValue{
								PermissionValue: &proto.Permission{
									Operation:  proto.Permission_OPERATION_UPDATE,
									Constraint: proto.Permission_CONSTRAINT_OWN,
								},
							},
						},
					},
					Resource: &proto.Resource{
						Type: proto.Resource_TYPE_SYSTEM,
					},
				}},
			expectedBundles: []expectedBundle{
				{displayName: "test role - update", name: "TEST_ROLE"},
				{displayName: "an other role", name: "AnOtherROLE"},
				{displayName: "Guest", name: "guest"},
				{displayName: "Admin", name: "admin"},
				{displayName: "User", name: "user"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			for _, bundle := range tt.bundles {
				_, err := bundleService.SaveBundle(context.Background(), &proto.SaveBundleRequest{
					Bundle: bundle,
				})
				assert.NoError(t, err)
			}
			rolesRes, err := roleService.ListRoles(context.Background(), &proto.ListBundlesRequest{})
			assert.NoError(t, err)

			for _, bundle := range rolesRes.Bundles {
				assert.Contains(t, tt.expectedBundles, expectedBundle{
					displayName: bundle.DisplayName,
					name:        bundle.Name,
				})
			}
			assert.Equal(t, len(tt.expectedBundles) , len(rolesRes.Bundles))
		})
	}
}

func setFullReadWriteOnBundle(t *testing.T, accountID, bundleID string) {
	permissionRequest := proto.AddSettingToBundleRequest{
		BundleId: svc.BundleUUIDRoleAdmin,
		Setting: &proto.Setting{
			Name: "test-bundle-permission-readwrite",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_BUNDLE,
				Id:   bundleID,
			},
			Value: &proto.Setting_PermissionValue{
				PermissionValue: &proto.Permission{
					Operation:  proto.Permission_OPERATION_READWRITE,
					Constraint: proto.Permission_CONSTRAINT_ALL,
				},
			},
		},
	}
	addPermissionResponse, err := bundleService.AddSettingToBundle(context.Background(), &permissionRequest)
	assert.NoError(t, err)
	if err == nil {
		assert.NotEmpty(t, addPermissionResponse.Setting)
	}

	_, err = roleService.AssignRoleToUser(
		context.Background(),
		&proto.AssignRoleToUserRequest{AccountUuid: accountID, RoleId: svc.BundleUUIDRoleAdmin},
	)
	assert.NoError(t, err)
}
