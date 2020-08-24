package proto_test

import (
	"context"
	"encoding/json"
	fmt "fmt"
	"log"
	"os"
	"testing"

	merrors "github.com/micro/go-micro/v2/errors"

	mgrpc "github.com/micro/go-micro/v2/client/grpc"
	ocislog "github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-pkg/v2/service/grpc"
	"github.com/owncloud/ocis-settings/pkg/config"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	svc "github.com/owncloud/ocis-settings/pkg/service/v0"
	store "github.com/owncloud/ocis-settings/pkg/store/filesystem"
	"github.com/stretchr/testify/assert"
)

var service = grpc.Service{}

var (
	settingsStub = []*proto.Setting{
		{
			Id:          "336c4db1-5062-4931-990f-d88e6b02cb02",
			DisplayName: "dummy setting",
			Name:        "dummy-setting",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_BUNDLE,
			},
			Value: &proto.Setting_IntValue{
				IntValue: &proto.Int{
					Default: 42,
				},
			},
			Description: "dummy setting",
		},
	}
)

const dataStore = "/var/tmp/ocis-settings"

func init() {
	service = grpc.NewService(
		grpc.Namespace("com.owncloud.api"),
		grpc.Name("settings"),
		grpc.Address("localhost:9992"),
	)

	cfg := config.New()
	cfg.Storage.DataPath = dataStore
	// Service initialization is not reliable. It most lilely happens
	// asynchronous causing a data race in some tests where it needs
	// as service but this is not available.
	err := proto.RegisterBundleServiceHandler(service.Server(), svc.NewService(cfg, ocislog.NewLogger(ocislog.Color(true), ocislog.Pretty(true))))
	if err != nil {
		log.Fatalf("could not register BundleServiceHandler: %v", err)
	}
	err = proto.RegisterValueServiceHandler(service.Server(), svc.NewService(cfg, ocislog.NewLogger(ocislog.Color(true), ocislog.Pretty(true))))
	if err != nil {
		log.Fatalf("could not register ValueServiceHandler: %v", err)
	}
	err = proto.RegisterRoleServiceHandler(service.Server(), svc.NewService(cfg, ocislog.NewLogger(ocislog.Color(true), ocislog.Pretty(true))))
	if err != nil {
		log.Fatalf("could not register RegisterRoleServiceHandler: %v", err)
	}

	if err = service.Server().Start(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}

/**
Make sure the default roles are created
The tests delete the Datastore after the run, so we need to recreate it
*/
func rebuidDataStore() {
	cfg := config.New()
	cfg.Storage.DataPath = dataStore
	svc.NewService(cfg, ocislog.NewLogger(ocislog.Color(true), ocislog.Pretty(true)))
}

/**
testing that saving a settings bundle and retrieving it again works correctly
using various setting bundle properties
*/
func TestSettingsBundleProperties(t *testing.T) {
	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)

	var scenarios = []struct {
		name          string
		bundleName    string
		displayName   string
		extensionName string
		UUID          string
		expectedError error
	}{
		{
			"ASCII",
			"bundle-name",
			"simple-bundle-key",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			nil,
		},
		{
			"UTF validation on bundle name",
			"सिम्प्ले-bundle-name",
			"सिम्प्ले-display-name",
			"सिम्प्ले-extension-name",
			"सिम्प्ले",
			merrors.New("ocis-settings", "extension: must be in a valid format; name: must be in a valid format.", 400),
		},
		{
			"UTF validation on display name",
			"सिम्प्ले-bundle-name",
			"सिम्प्ले-display-name",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			merrors.New("ocis-settings", "name: must be in a valid format.", 400),
		},
		{
			"extension name with ../ in the name",
			"bundle-name",
			"simple-display-name",
			"../folder-a-level-higher-up",
			"123e4567-e89b-12d3-a456-426652340000",
			merrors.New("ocis-settings", "extension: must be in a valid format.", 400),
		},
		{
			"extension name with \\ in the name",
			"bundle-name",
			"simple-display-name",
			"\\",
			"123e4567-e89b-12d3-a456-426652340000",
			merrors.New("ocis-settings", "extension: must be in a valid format.", 400),
		},
		{
			"spaces are disallowed in bundle names",
			"bundle name",
			"simple display name",
			"simple extension name",
			"123e4567-e89b-12d3-a456-426652340000",
			merrors.New("ocis-settings", "extension: must be in a valid format; name: must be in a valid format.", 400),
		},
		{
			"spaces are allowed in display names",
			"bundle-name",
			"simple display name",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			nil,
		},
		{
			"extension missing",
			"bundle-name",
			"simple-display-name",
			"",
			"123e4567-e89b-12d3-a456-426652340000",
			merrors.New("ocis-settings", "extension: cannot be blank.", 400),
		},
		{
			"display name missing",
			"bundleName",
			"",
			"simple-extension-name",
			"123e4567-e89b-12d3-a456-426652340000",
			merrors.New("ocis-settings", "display_name: cannot be blank.", 400),
		},
		{
			"UUID missing (omitted on bundles)",
			"bundle-name",
			"simple-display-name",
			"simple-extension-name",
			"",
			nil,
		},
	}
	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			bundle := proto.Bundle{
				Name:        scenario.bundleName,
				Extension:   scenario.extensionName,
				DisplayName: scenario.displayName,
				Type:        proto.Bundle_TYPE_DEFAULT,
				Resource: &proto.Resource{
					Type: proto.Resource_TYPE_BUNDLE,
				},
				Settings: settingsStub,
			}
			createRequest := proto.SaveBundleRequest{
				Bundle: &bundle,
			}

			cresponse, err := cl.SaveBundle(context.Background(), &createRequest)
			if err != nil || scenario.expectedError != nil {
				t.Log(err)
				assert.Equal(t, scenario.expectedError, err)
			} else {
				assert.Equal(t, scenario.extensionName, cresponse.Bundle.Extension)
				assert.Equal(t, scenario.displayName, cresponse.Bundle.DisplayName)
				getRequest := proto.GetBundleRequest{BundleId: cresponse.Bundle.Id}
				getResponse, err := cl.GetBundle(context.Background(), &getRequest)
				assert.NoError(t, err)
				assert.Equal(t, scenario.displayName, getResponse.Bundle.DisplayName)
			}
			os.RemoveAll(dataStore)
		})
	}
}

func TestSettingsBundleWithoutSettings(t *testing.T) {
	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)

	createRequest := proto.SaveBundleRequest{
		Bundle: &proto.Bundle{
			DisplayName: "Alice's Bundle",
		},
	}
	response, err := cl.SaveBundle(context.Background(), &createRequest)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, merrors.New("ocis-settings", "extension: cannot be blank; name: cannot be blank; settings: cannot be blank.", 400), err)
	os.RemoveAll(dataStore)
}

// /**
// testing that setting getting and listing a settings bundle works correctly with a set of setting definitions
// */
func TestSaveGetListSettingsBundle(t *testing.T) {
	//options for choice list settings
	options := []*proto.ListOption{
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
	multipleChoiceSetting := proto.MultiChoiceList{
		Options: options,
	}

	//SingleChoiceList
	singleChoiceSetting := proto.SingleChoiceList{
		Options: options,
	}

	settings := []*proto.Setting{
		{
			Name:        "int",
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
				Type: proto.Resource_TYPE_BUNDLE,
			},
		},
		{
			Name:        "string",
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
				Type: proto.Resource_TYPE_BUNDLE,
			},
		},
		{
			Name:        "bool",
			DisplayName: "a bool value",
			Description: "with some description",
			Value: &proto.Setting_BoolValue{
				BoolValue: &proto.Bool{
					Default: false,
					Label:   "bool setting",
				},
			},
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_BUNDLE,
			},
		},
		{
			Name:        "multipleChoice",
			DisplayName: "a multiple choice setting",
			Description: "with some description",
			Value: &proto.Setting_MultiChoiceValue{
				MultiChoiceValue: &multipleChoiceSetting,
			},
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_BUNDLE,
			},
		},
		{
			Name:        "singleChoice",
			DisplayName: "a single choice setting",
			Description: "with some description",
			Value: &proto.Setting_SingleChoiceValue{
				SingleChoiceValue: &singleChoiceSetting,
			},
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_BUNDLE,
			},
		},
	}

	bundle := proto.Bundle{
		Name:        "test",
		DisplayName: "bundleDisplayName",
		Extension:   "testExtension",
		Type:        proto.Bundle_TYPE_DEFAULT,
		Settings:    settings,
		Resource: &proto.Resource{
			Type: proto.Resource_TYPE_BUNDLE,
		},
	}
	saveRequest := proto.SaveBundleRequest{
		Bundle: &bundle,
	}

	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)

	// assert that SaveBundle returns the same bundle as we have sent
	saveResponse, err := cl.SaveBundle(context.Background(), &saveRequest)
	assert.NoError(t, err)
	receivedBundle, _ := json.Marshal(saveResponse.Bundle.Settings)
	expectedBundle, _ := json.Marshal(&bundle.Settings)
	assert.Equal(t, receivedBundle, expectedBundle)

	//assert that GetBundle returns the same bundle as saved
	getRequest := proto.GetBundleRequest{BundleId: saveResponse.Bundle.Id}
	getResponse, err := cl.GetBundle(context.Background(), &getRequest)
	assert.NoError(t, err)
	receivedBundle, _ = json.Marshal(getResponse.Bundle.Settings)
	assert.Equal(t, expectedBundle, receivedBundle)

	os.RemoveAll(dataStore)
}

// https://github.com/owncloud/ocis-settings/issues/18
func TestSaveSettingsBundleWithInvalidSettingValues(t *testing.T) {
	var tests = []proto.Setting{
		{
			Name: "intValue",
			Resource: &proto.Resource{
				Type: proto.Resource_TYPE_BUNDLE,
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
				Type: proto.Resource_TYPE_BUNDLE,
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
				Type: proto.Resource_TYPE_BUNDLE,
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
				Type: proto.Resource_TYPE_BUNDLE,
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
				Type: proto.Resource_TYPE_BUNDLE,
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
				Type: proto.Resource_TYPE_BUNDLE,
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
				Type: proto.Resource_TYPE_BUNDLE,
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
				Type: proto.Resource_TYPE_BUNDLE,
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
				Type: proto.Resource_TYPE_BUNDLE,
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
				Type: proto.Resource_TYPE_BUNDLE,
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
				Type: proto.Resource_TYPE_BUNDLE,
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
		index := index
		t.Run(tests[index].Name, func(t *testing.T) {

			var settings []*proto.Setting

			settings = append(settings, &tests[index])
			bundle := proto.Bundle{
				Name:        "bundle",
				Extension:   "bundleExtension",
				DisplayName: "bundledisplayname",
				Type:        proto.Bundle_TYPE_DEFAULT,
				Resource: &proto.Resource{
					Type: proto.Resource_TYPE_BUNDLE,
				},
				Settings: settings,
			}
			saveRequest := proto.SaveBundleRequest{
				Bundle: &bundle,
			}

			client := service.Client()
			cl := proto.NewBundleService("com.owncloud.api.settings", client)

			//assert that SaveBundle returns the same bundle as we have sent there
			saveResponse, err := cl.SaveBundle(context.Background(), &saveRequest)
			assert.NoError(t, err)
			receivedBundle, _ := json.Marshal(saveResponse.Bundle.Settings)
			expectedBundle, _ := json.Marshal(&bundle.Settings)
			assert.Equal(t, expectedBundle, receivedBundle)
			os.RemoveAll(dataStore)
		})
	}
}

// https://github.com/owncloud/ocis-settings/issues/19
func TestGetSettingsBundleCreatesFolder(t *testing.T) {
	client := service.Client()
	cl := proto.NewBundleService("com.owncloud.api.settings", client)
	getRequest := proto.GetBundleRequest{BundleId: "non-existing-bundle"}

	_, _ = cl.GetBundle(context.Background(), &getRequest)
	assert.NoDirExists(t, store.Name+"/bundles/non-existing-bundle")
	assert.NoFileExists(t, store.Name+"/bundles/non-existing-bundle/not-existing-bundle.json")
	os.RemoveAll(dataStore)
}

// TODO non-deterministic. Fix.
func TestCreateRoleAndAssign(t *testing.T) {
	c := mgrpc.NewClient()
	bundleC := proto.NewBundleService("com.owncloud.api.settings", c)

	res, err := bundleC.SaveBundle(context.Background(), &proto.SaveBundleRequest{
		Bundle: &proto.Bundle{
			// Id:          "f36db5e6-a03c-40df-8413-711c67e40b47", // bug: providing the ID ignores its value for the filename.
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
		rolesC := proto.NewRoleService("com.owncloud.api.settings", c)
		_, err = rolesC.AssignRoleToUser(context.Background(), &proto.AssignRoleToUserRequest{
			AccountUuid: "4c510ada-c86b-4815-8820-42cdf82c3d51",
			RoleId:      res.Bundle.Id,
		})
		if err != nil {
			fmt.Println(err)
		}
		assert.NoError(t, err)
	}
	os.RemoveAll(dataStore)
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
// 	os.RemoveAll(dataStore)
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
			},
		},
	}

	rebuidDataStore()
	client := service.Client()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleService := proto.NewRoleService("com.owncloud.api.settings", client)

			bundleService := proto.NewBundleService("com.owncloud.api.settings", client)
			for _, bundle := range tt.bundles {
				_, err := bundleService.SaveBundle(context.Background(), &proto.SaveBundleRequest{
					Bundle: bundle,
				})
				assert.NoError(t, err)
			}
			rolesRes, err := roleService.ListRoles(context.Background(), &proto.ListBundlesRequest{AccountUuid: "."})
			assert.NoError(t, err)

			for _, bundle := range rolesRes.Bundles {
				assert.Contains(t, tt.expectedBundles, expectedBundle{
					displayName: bundle.DisplayName,
					name:        bundle.Name,
				})
			}
			assert.Equal(t, len(tt.expectedBundles), len(rolesRes.Bundles))

			os.RemoveAll(dataStore)
		})
	}
}

func TestListRolesVariousAccountUuid(t *testing.T) {
	tests := []struct {
		name        string
		accountUUID string
		error       interface{}
	}{
		{"space",
			" ",
			"{\"id\":\"ocis-settings\",\"code\":400,\"detail\":\"must be in a valid format\",\"status\":\"Bad Request\"}",
		},
		{"empty",
			"",
			"{\"id\":\"ocis-settings\",\"code\":400,\"detail\":\"cannot be blank\",\"status\":\"Bad Request\"}",
		},
		{"numbers",
			"123",
			nil,
		},
		{"dot",
			".",
			nil,
		},
		{"string",
			"abc",
			nil,
		},
		{"UTF",
			"नेपाल",
			"{\"id\":\"ocis-settings\",\"code\":400,\"detail\":\"must be in a valid format\",\"status\":\"Bad Request\"}",
		},
		{"string with spaces",
			"abc de",
			"{\"id\":\"ocis-settings\",\"code\":400,\"detail\":\"must be in a valid format\",\"status\":\"Bad Request\"}",
		},
	}

	rebuidDataStore()
	client := service.Client()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roleService := proto.NewRoleService("com.owncloud.api.settings", client)

			roles, err := roleService.ListRoles(
				context.Background(), &proto.ListBundlesRequest{AccountUuid: tt.accountUUID},
			)
			if tt.error != nil {
				assert.EqualError(t, err, tt.error.(string))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 3, len(roles.Bundles))
			}
		})
	}
}
