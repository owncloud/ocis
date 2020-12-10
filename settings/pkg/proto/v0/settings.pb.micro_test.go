package proto_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	ocislog "github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/settings/pkg/config"
	"github.com/owncloud/ocis/settings/pkg/proto/v0"
	svc "github.com/owncloud/ocis/settings/pkg/service/v0"
	store "github.com/owncloud/ocis/settings/pkg/store/filesystem"
	"github.com/stretchr/testify/assert"
)

var (
	service           grpc.Service
	handler           svc.Service
	bundleService     proto.BundleService
	valueService      proto.ValueService
	roleService       proto.RoleService
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

const dataPath = "/tmp/grpc-tests-ocis-settings"

func init() {
	service = grpc.NewService(
		grpc.Namespace("com.owncloud.api"),
		grpc.Name("settings"),
		grpc.Address("localhost:9992"),
	)

	cfg := config.New()
	cfg.Service.DataPath = dataPath
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
			merrors.New("ocis-settings", "extension: must be in a valid format; name: must be in a valid format.", http.StatusBadRequest),
		},
		{
			"UTF validation on display name",
			"सिम्प्ले-bundle-name",
			"सिम्प्ले-display-name",
			"simple-extension-name",
			merrors.New("ocis-settings", "name: must be in a valid format.", http.StatusBadRequest),
		},
		{
			"extension name with ../ in the name",
			"bundle-name",
			"simple-display-name",
			"../folder-a-level-higher-up",
			merrors.New("ocis-settings", "extension: must be in a valid format.", http.StatusBadRequest),
		},
		{
			"extension name with \\ in the name",
			"bundle-name",
			"simple-display-name",
			"\\",
			merrors.New("ocis-settings", "extension: must be in a valid format.", http.StatusBadRequest),
		},
		{
			"spaces are disallowed in bundle names",
			"bundle name",
			"simple display name",
			"simple extension name",
			merrors.New("ocis-settings", "extension: must be in a valid format; name: must be in a valid format.", http.StatusBadRequest),
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
			merrors.New("ocis-settings", "extension: cannot be blank.", http.StatusBadRequest),
		},
		{
			"display name missing",
			"bundleName",
			"",
			"simple-extension-name",
			merrors.New("ocis-settings", "display_name: cannot be blank.", http.StatusBadRequest),
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
			ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

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

			cresponse, err := bundleService.SaveBundle(ctx, &createRequest)
			if err != nil || scenario.expectedError != nil {
				t.Log(err)
				assert.Equal(t, scenario.expectedError, err)
			} else {
				assert.Equal(t, scenario.extensionName, cresponse.Bundle.Extension)
				assert.Equal(t, scenario.displayName, cresponse.Bundle.DisplayName)

				// we want to test input validation, so just allow the request permission-wise
				setFullReadWriteOnBundleForAdmin(ctx, t, cresponse.Bundle.Id)

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

	ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
	ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

	createRequest := proto.SaveBundleRequest{
		Bundle: &proto.Bundle{
			DisplayName: "Alice's Bundle",
		},
	}
	response, err := bundleService.SaveBundle(ctx, &createRequest)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, merrors.New("ocis-settings", "extension: cannot be blank; name: cannot be blank; settings: cannot be blank.", http.StatusBadRequest), err)
}

func TestGetBundleOfABundleSavedWithoutPermissions(t *testing.T) {
	teardown := setup()
	defer teardown()

	ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
	ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

	saveRequest := proto.SaveBundleRequest{
		Bundle: &bundleStub,
	}
	saveResponse, err := bundleService.SaveBundle(ctx, &saveRequest)
	assert.NoError(t, err)
	assert.Equal(t, bundleStub.Id, saveResponse.Bundle.Id)

	getRequest := proto.GetBundleRequest{BundleId: bundleStub.Id}
	getResponse, err := bundleService.GetBundle(ctx, &getRequest)
	assert.Empty(t, getResponse)

	assert.Equal(t, merrors.New("ocis-settings", "could not read bundle: b1b8c9d0-fb3c-4e12-b868-5a8508218d2e", http.StatusNotFound), err)
}

func TestGetBundleHavingFullPermissionsOnAnotherRole(t *testing.T) {
	teardown := setup()
	defer teardown()

	ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
	ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

	saveRequest := proto.SaveBundleRequest{
		Bundle: &bundleStub,
	}
	saveResponse, err := bundleService.SaveBundle(ctx, &saveRequest)
	assert.NoError(t, err)
	assert.Equal(t, bundleStub.Id, saveResponse.Bundle.Id)

	setFullReadWriteOnBundleForAdmin(ctx, t, bundleStub.Id)

	ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleUser))
	getRequest := proto.GetBundleRequest{BundleId: bundleStub.Id}
	getResponse, err := bundleService.GetBundle(ctx, &getRequest)
	assert.Empty(t, getResponse)

	assert.Equal(t, merrors.New("ocis-settings", "could not read bundle: b1b8c9d0-fb3c-4e12-b868-5a8508218d2e", http.StatusNotFound), err)
}

/**
testing that setting getting and listing a settings bundle works correctly with a set of setting definitions
*/
func TestSaveAndGetBundle(t *testing.T) {
	teardown := setup()
	defer teardown()

	ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
	ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

	saveRequest := proto.SaveBundleRequest{
		Bundle: &bundleStub,
	}

	// assert that SaveBundle returns the same bundle as we have sent
	saveResponse, err := bundleService.SaveBundle(ctx, &saveRequest)
	assert.NoError(t, err)
	receivedBundle, _ := json.Marshal(saveResponse.Bundle.Settings)
	expectedBundle, _ := json.Marshal(&bundleStub.Settings)
	assert.Equal(t, receivedBundle, expectedBundle)

	// set full permissions for getting the created bundle
	setFullReadWriteOnBundleForAdmin(ctx, t, saveResponse.Bundle.Id)

	//assert that GetBundle returns the same bundle as saved
	getRequest := proto.GetBundleRequest{BundleId: saveResponse.Bundle.Id}
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
			// https://github.com/owncloud/ocis/settings/issues/57
		},
		{
			name:  "less than Min",
			value: proto.Value_IntValue{IntValue: 0},
			// https://github.com/owncloud/ocis/settings/issues/57
		},
		{
			name:  "more than Max",
			value: proto.Value_IntValue{IntValue: 128},
			// https://github.com/owncloud/ocis/settings/issues/57
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
			ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

			saveResponse, err := bundleService.SaveBundle(ctx, &proto.SaveBundleRequest{
				Bundle: &bundleStub,
			})
			assert.NoError(t, err)

			saveValueResponse, err := valueService.SaveValue(ctx, &proto.SaveValueRequest{
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
				ctx, &proto.GetValueRequest{Id: saveValueResponse.Value.Value.Id},
			)
			assert.NoError(t, err)
			assert.Equal(t, tt.value.IntValue, getValueResponse.Value.Value.GetIntValue())
		})
	}
}

/**
try to save a wrong type of the value
https://github.com/owncloud/ocis/settings/issues/57
*/
func TestSaveGetIntValueIntoString(t *testing.T) {
	teardown := setup()
	defer teardown()

	ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
	ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

	saveResponse, err := bundleService.SaveBundle(ctx, &proto.SaveBundleRequest{
		Bundle: &bundleStub,
	})
	assert.NoError(t, err)

	saveValueResponse, err := valueService.SaveValue(ctx, &proto.SaveValueRequest{
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
		ctx, &proto.GetValueRequest{Id: saveValueResponse.Value.Value.Id},
	)
	assert.NoError(t, err)
	assert.Equal(t, "forty two", getValueResponse.Value.Value.GetStringValue())
}

// https://github.com/owncloud/ocis/settings/issues/18
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

			ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
			ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

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
			saveResponse, err := bundleService.SaveBundle(ctx, &saveRequest)
			assert.NoError(t, err)
			receivedBundle, _ := json.Marshal(saveResponse.Bundle.Settings)
			expectedBundle, _ := json.Marshal(&bundle.Settings)
			assert.Equal(t, expectedBundle, receivedBundle)
		})
	}
}

// https://github.com/owncloud/ocis/settings/issues/19
func TestGetBundleNoSideEffectsOnDisk(t *testing.T) {
	teardown := setup()
	defer teardown()

	ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
	ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

	getRequest := proto.GetBundleRequest{BundleId: "non-existing-bundle"}

	_, _ = bundleService.GetBundle(ctx, &getRequest)
	assert.NoDirExists(t, store.Name+"/bundles/non-existing-bundle")
	assert.NoFileExists(t, store.Name+"/bundles/non-existing-bundle/not-existing-bundle.json")
}

// TODO non-deterministic. Fix.
func TestCreateRoleAndAssign(t *testing.T) {
	teardown := setup()
	defer teardown()

	ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
	ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

	res, err := bundleService.SaveBundle(ctx, &proto.SaveBundleRequest{
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
		_, err = roleService.AssignRoleToUser(ctx, &proto.AssignRoleToUserRequest{
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

			ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
			ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

			for _, bundle := range tt.bundles {
				_, err := bundleService.SaveBundle(ctx, &proto.SaveBundleRequest{
					Bundle: bundle,
				})
				assert.NoError(t, err)
			}
			rolesRes, err := roleService.ListRoles(ctx, &proto.ListBundlesRequest{})
			assert.NoError(t, err)

			for _, bundle := range rolesRes.Bundles {
				assert.Contains(t, tt.expectedBundles, expectedBundle{
					displayName: bundle.DisplayName,
					name:        bundle.Name,
				})
			}
			assert.Equal(t, len(tt.expectedBundles), len(rolesRes.Bundles))
		})
	}
}

func TestListFilteredBundle(t *testing.T) {
	type expectedBundle struct {
		displayName string
		name        string
	}

	type permission struct {
		permission proto.Permission_Operation
		roleUUID   string
	}

	type bundleForTest struct {
		bundle     *proto.Bundle
		permission permission
	}

	tests := []struct {
		name            string
		bundles         []bundleForTest
		expectedBundles []expectedBundle
	}{
		{
			name: "multiple bundles, all have READ(WRITE) permission",
			bundles: []bundleForTest{
				{
					bundle: &proto.Bundle{
						Name:        "test",
						Id:          "b1b8c9d0-fb3c-4e12-b868-5a8508218d2e",
						DisplayName: "bundleDisplayName",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READWRITE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					bundle: &proto.Bundle{
						Name:        "one-more",
						Id:          "3b9f230a-fc9e-4605-89ee-a21e24728c64",
						DisplayName: "an other bundle",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READ,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
			},
			expectedBundles: []expectedBundle{
				{displayName: "bundleDisplayName", name: "test"},
				{displayName: "an other bundle", name: "one-more"},
			},
		},
		{
			name: "multiple bundles, only one with READ permission",
			bundles: []bundleForTest{
				{
					bundle: &proto.Bundle{
						Name:        "Permission_OPERATION_WRITE",
						Id:          "12fe2b67-4a08-4f17-9cb6-924da943da0e",
						DisplayName: "Permission_OPERATION_WRITE",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_WRITE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					bundle: &proto.Bundle{
						Name:        "Permission_OPERATION_DELETE",
						Id:          "1a0b65b0-fdbf-4738-b41e-41d36a01376e",
						DisplayName: "Permission_OPERATION_DELETE",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_DELETE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					bundle: &proto.Bundle{
						Name:        "Permission_OPERATION_UPDATE",
						Id:          "511fe78e-89c9-4237-a01e-6af5457a135e",
						DisplayName: "Permission_OPERATION_UPDATE",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_UPDATE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					bundle: &proto.Bundle{
						Name:        "Permission_OPERATION_CREATE",
						Id:          "aa42fb12-57aa-40c0-b458-3a91f398deba",
						DisplayName: "Permission_OPERATION_CREATE",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_CREATE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					bundle: &proto.Bundle{
						Name:        "Permission_OPERATION_UNKNOWN",
						Id:          "eabb2a18-09e2-4b06-aa62-987e8dc5e908",
						DisplayName: "Permission_OPERATION_UNKNOWN",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_UNKNOWN,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					bundle: &proto.Bundle{
						Name:        "Permission_OPERATION_READ",
						Id:          "3b9f230a-fc9e-4605-89ee-a21e24728c64",
						DisplayName: "Permission_OPERATION_READ",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READ,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
			},
			expectedBundles: []expectedBundle{
				{displayName: "Permission_OPERATION_READ", name: "Permission_OPERATION_READ"},
			},
		},
		{
			name: "multiple bundles, all have READ permission, but only one matching role",
			bundles: []bundleForTest{
				{
					bundle: &proto.Bundle{
						Name:        "matching-role",
						Id:          "b1b8c9d0-fb3c-4e12-b868-5a8508218d2e",
						DisplayName: "bundleDisplayName",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READWRITE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					bundle: &proto.Bundle{
						Name:        "NOT-matching-role",
						Id:          "3b9f230a-fc9e-4605-89ee-a21e24728c64",
						DisplayName: "an other bundle",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READ,
						roleUUID:   svc.BundleUUIDRoleGuest,
					},
				},
				{
					bundle: &proto.Bundle{
						Name:        "NOT-matching-role2",
						Id:          "714a5917-627c-40ac-8dc7-5fdac013e4b7",
						DisplayName: "an other bundle",
						Extension:   "testExtension",
						Type:        proto.Bundle_TYPE_DEFAULT,
						Settings:    complexSettingsStub,
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_SYSTEM,
						},
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READ,
						roleUUID:   svc.BundleUUIDRoleUser,
					},
				},
			},
			expectedBundles: []expectedBundle{
				{displayName: "bundleDisplayName", name: "matching-role"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
			ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

			for _, testBundle := range tt.bundles {
				_, err := bundleService.SaveBundle(ctx, &proto.SaveBundleRequest{
					Bundle: testBundle.bundle,
				})
				assert.NoError(t, err)

				setPermissionOnBundleOrSetting(
					ctx, t, testBundle.bundle.Id, proto.Resource_TYPE_BUNDLE,
					testBundle.permission.permission, testBundle.permission.roleUUID,
				)
			}

			listRes, err := bundleService.ListBundles(ctx, &proto.ListBundlesRequest{})
			assert.NoError(t, err)

			for _, bundle := range listRes.Bundles {
				assert.Contains(t, tt.expectedBundles, expectedBundle{
					displayName: bundle.DisplayName,
					name:        bundle.Name,
				})
			}
			assert.Equal(t, len(tt.expectedBundles), len(listRes.Bundles))
		})
	}
}

func TestListGetBundleSettingMixedPermission(t *testing.T) {
	type expectedSetting struct {
		displayName string
		name        string
	}

	type permission struct {
		permission proto.Permission_Operation
		roleUUID   string
	}

	type settingsForTest struct {
		setting    *proto.Setting
		permission permission
	}

	tests := []struct {
		name             string
		settings         []settingsForTest
		expectedSettings []expectedSetting
	}{
		{
			name: "all settings have R/RW permissions",
			settings: []settingsForTest{
				{
					setting: &proto.Setting{
						Id:          "b86fdb0a-801f-4749-ab84-5c99e90dbd6d",
						DisplayName: "RW setting",
						Name:        "RW-setting",
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_USER,
						},
						Value: &proto.Setting_IntValue{
							IntValue: &proto.Int{
								Default: 42,
							},
						},
						Description: "RW setting",
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READWRITE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					setting: &proto.Setting{
						Id:          "cb1bbe58-27e7-461b-91b1-a9c85c488789",
						DisplayName: "RO setting",
						Name:        "RO-setting",
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_USER,
						},
						Value: &proto.Setting_IntValue{
							IntValue: &proto.Int{
								Default: 42,
							},
						},
						Description: "RO setting",
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READWRITE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
			},
			expectedSettings: []expectedSetting{
				{displayName: "RW setting", name: "RW-setting"},
				{displayName: "RO setting", name: "RO-setting"},
			},
		},
		{
			name: "all settings have R/RW permissions but only one the matching user",
			settings: []settingsForTest{
				{
					setting: &proto.Setting{
						Id:          "b86fdb0a-801f-4749-ab84-5c99e90dbd6d",
						DisplayName: "matching user",
						Name:        "matching-user",
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_USER,
						},
						Value: &proto.Setting_IntValue{
							IntValue: &proto.Int{
								Default: 42,
							},
						},
						Description: "matching user",
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READWRITE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					setting: &proto.Setting{
						Id:          "cb1bbe58-27e7-461b-91b1-a9c85c488789",
						DisplayName: "NOT matching user",
						Name:        "NOT-matching-user",
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_USER,
						},
						Value: &proto.Setting_IntValue{
							IntValue: &proto.Int{
								Default: 42,
							},
						},
						Description: "NOT matching user",
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READWRITE,
						roleUUID:   svc.BundleUUIDRoleGuest,
					},
				},
			},
			expectedSettings: []expectedSetting{
				{displayName: "matching user", name: "matching-user"},
			},
		},
		{
			name: "only one settings has READ permissions",
			settings: []settingsForTest{
				{
					setting: &proto.Setting{
						Id:          "b86fdb0a-801f-4749-ab84-5c99e90dbd6d",
						DisplayName: "WRITE setting",
						Name:        "WRITE-setting",
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_USER,
						},
						Value: &proto.Setting_IntValue{
							IntValue: &proto.Int{
								Default: 42,
							},
						},
						Description: "WRITE setting",
					},
					permission: permission{
						permission: proto.Permission_OPERATION_WRITE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					setting: &proto.Setting{
						Id:          "6163c6bf-79f2-43f7-b0ba-1493534bfc10",
						DisplayName: "UNKNOWN setting",
						Name:        "UNKNOWN-setting",
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_USER,
						},
						Value: &proto.Setting_IntValue{
							IntValue: &proto.Int{
								Default: 42,
							},
						},
						Description: "UNKNOWN setting",
					},
					permission: permission{
						permission: proto.Permission_OPERATION_UNKNOWN,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					setting: &proto.Setting{
						Id:          "79eda727-9fa1-459f-aaff-f73ed5693419",
						DisplayName: "CREATE setting",
						Name:        "CREATE-setting",
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_USER,
						},
						Value: &proto.Setting_IntValue{
							IntValue: &proto.Int{
								Default: 42,
							},
						},
						Description: "CREATE setting",
					},
					permission: permission{
						permission: proto.Permission_OPERATION_CREATE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					setting: &proto.Setting{
						Id:          "2be7ca51-89fb-4968-b9d2-0ac43197adff",
						DisplayName: "UPDATE setting",
						Name:        "UPDATE-setting",
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_USER,
						},
						Value: &proto.Setting_IntValue{
							IntValue: &proto.Int{
								Default: 42,
							},
						},
						Description: "UPDATE setting",
					},
					permission: permission{
						permission: proto.Permission_OPERATION_UPDATE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					setting: &proto.Setting{
						Id:          "f1a0005e-e570-4bd8-a18c-b4afaaa8d7d9",
						DisplayName: "DELETE setting",
						Name:        "DELETE-setting",
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_USER,
						},
						Value: &proto.Setting_IntValue{
							IntValue: &proto.Int{
								Default: 42,
							},
						},
						Description: "DELETE setting",
					},
					permission: permission{
						permission: proto.Permission_OPERATION_DELETE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
				{
					setting: &proto.Setting{
						Id:          "cb1bbe58-27e7-461b-91b1-a9c85c488789",
						DisplayName: "RO setting",
						Name:        "RO-setting",
						Resource: &proto.Resource{
							Type: proto.Resource_TYPE_USER,
						},
						Value: &proto.Setting_IntValue{
							IntValue: &proto.Int{
								Default: 42,
							},
						},
						Description: "RO setting",
					},
					permission: permission{
						permission: proto.Permission_OPERATION_READWRITE,
						roleUUID:   svc.BundleUUIDRoleAdmin,
					},
				},
			},
			expectedSettings: []expectedSetting{
				{displayName: "RO setting", name: "RO-setting"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
			ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

			// create bundle with the defined settings
			var settings []*proto.Setting

			for _, testSetting := range tt.settings {
				settings = append(settings, testSetting.setting)
			}

			bundle := proto.Bundle{
				Name:        "test",
				Id:          "b1b8c9d0-fb3c-4e12-b868-5a8508218d2e",
				DisplayName: "bundleDisplayName",
				Extension:   "testExtension",
				Type:        proto.Bundle_TYPE_DEFAULT,
				Settings:    settings,
				Resource: &proto.Resource{
					Type: proto.Resource_TYPE_SYSTEM,
				},
			}

			_, err := bundleService.SaveBundle(ctx, &proto.SaveBundleRequest{
				Bundle: &bundle,
			})
			assert.NoError(t, err)

			// set permissions for each setting
			for _, testSetting := range tt.settings {
				setPermissionOnBundleOrSetting(
					ctx, t, testSetting.setting.Id, proto.Resource_TYPE_SETTING,
					testSetting.permission.permission, testSetting.permission.roleUUID,
				)
			}

			listRes, err := bundleService.ListBundles(ctx, &proto.ListBundlesRequest{})
			assert.NoError(t, err)

			for _, setting := range listRes.Bundles[0].Settings {
				assert.Contains(t, tt.expectedSettings, expectedSetting{
					displayName: setting.DisplayName,
					name:        setting.Name,
				})
			}
			assert.Equal(t, len(tt.expectedSettings), len(listRes.Bundles[0].Settings))

			getRes, err := bundleService.GetBundle(ctx, &proto.GetBundleRequest{BundleId: bundle.Id})
			assert.NoError(t, err)

			for _, setting := range getRes.Bundle.Settings {
				assert.Contains(t, tt.expectedSettings, expectedSetting{
					displayName: setting.DisplayName,
					name:        setting.Name,
				})
			}
			assert.Equal(t, len(tt.expectedSettings), len(getRes.Bundle.Settings))
		})
	}
}

func TestListFilteredBundle_SetPermissionsOnSettingAndBundle(t *testing.T) {
	tests := []struct {
		name                     string
		settingPermission        proto.Permission_Operation
		bundlePermission         proto.Permission_Operation
		expectedAmountOfSettings int
	}{
		{
			"setting has read permission bundle not",
			proto.Permission_OPERATION_READ,
			proto.Permission_OPERATION_UNKNOWN,
			1,
		},
		{
			"bundle has read permission setting not",
			proto.Permission_OPERATION_UNKNOWN,
			proto.Permission_OPERATION_READ,
			5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teardown := setup()
			defer teardown()

			ctx := metadata.Set(context.Background(), middleware.AccountID, testAccountID)
			ctx = metadata.Set(ctx, middleware.RoleIDs, getRoleIDAsJSON(svc.BundleUUIDRoleAdmin))

			_, err := bundleService.SaveBundle(ctx, &proto.SaveBundleRequest{
				Bundle: &bundleStub,
			})
			assert.NoError(t, err)

			setPermissionOnBundleOrSetting(
				ctx, t, bundleStub.Id, proto.Resource_TYPE_BUNDLE, tt.bundlePermission, svc.BundleUUIDRoleAdmin,
			)

			setPermissionOnBundleOrSetting(
				ctx, t, bundleStub.Settings[0].Id, proto.Resource_TYPE_SETTING,
				tt.settingPermission, svc.BundleUUIDRoleAdmin,
			)

			listRes, err := bundleService.ListBundles(ctx, &proto.ListBundlesRequest{})
			assert.NoError(t, err)
			assert.Equal(t, 1, len(listRes.Bundles))
			assert.Equal(t, tt.expectedAmountOfSettings, len(listRes.Bundles[0].Settings))
			assert.Equal(t, bundleStub.Id, listRes.Bundles[0].Id)
			assert.Equal(t, bundleStub.Settings[0].Id, listRes.Bundles[0].Settings[0].Id)
		})
	}
}

func setFullReadWriteOnBundleForAdmin(ctx context.Context, t *testing.T, bundleID string) {
	setPermissionOnBundleOrSetting(
		ctx, t, bundleID, proto.Resource_TYPE_BUNDLE, proto.Permission_OPERATION_READWRITE, svc.BundleUUIDRoleAdmin,
	)
}

func setPermissionOnBundleOrSetting(
	ctx context.Context,
	t *testing.T,
	bundleID string,
	resourceType proto.Resource_Type,
	permission proto.Permission_Operation,
	roleUUID string,
) {
	permissionRequest := proto.AddSettingToBundleRequest{
		BundleId: roleUUID,
		Setting: &proto.Setting{
			Name: "test-bundle-permission-readwrite",
			Resource: &proto.Resource{
				Type: resourceType,
				Id:   bundleID,
			},
			Value: &proto.Setting_PermissionValue{
				PermissionValue: &proto.Permission{
					Operation:  permission,
					Constraint: proto.Permission_CONSTRAINT_OWN,
				},
			},
		},
	}
	addPermissionResponse, err := bundleService.AddSettingToBundle(ctx, &permissionRequest)
	assert.NoError(t, err)
	if err == nil {
		assert.NotEmpty(t, addPermissionResponse.Setting)
	}
}

func getRoleIDAsJSON(roleID string) string {
	roleIDsJSON, _ := json.Marshal([]string{roleID})
	return string(roleIDsJSON)
}
