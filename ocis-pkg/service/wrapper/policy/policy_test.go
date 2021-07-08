package policy

import (
	"context"
	"io/ioutil"
	"testing"

	mgrpcc "github.com/asim/go-micro/plugins/client/grpc/v3"
	mbreaker "github.com/asim/go-micro/plugins/wrapper/breaker/gobreaker/v3"

	"github.com/asim/go-micro/v3/client"
)

type scenario struct {
	name     string
	service  string
	endpoint string
	cw       *clientWrapper
	error    bool
}

// because we have an init state we need a higher closure.
var (
	fname             string
	scenarios         []scenario
	defaultRegoPolicy = `package ocis

default deny = false

deny {
    input.service == "com.owncloud.api.thumbnails"
}

# deny {
#    input.service == "com.owncloud.api.settings"
#    input.endpoint == "AccountsService.ListAccounts"
#    input.method == "RoleService.ListRoleAssignments"
#    input.standard_claims.email == "admin@example.org"
#    input.standard_claims.groups == ""
#    input.standard_claims.iss == "https://localhost:9200"
#    input.standard_claims.name == "admin"
# }
`
)

func TestMain(m *testing.M) {
	f, err := ioutil.TempFile("", "policy.rego")
	if err != nil {
		panic(err)
	}

	// set package scoped global fname so fixtures have access to the generated file name
	fname = f.Name()

	// populate scenarios
	loadScenarios()

	// write default policy contents
	if _, err := f.Write([]byte(defaultRegoPolicy)); err != nil {
		panic(err)
	}

	m.Run()
}

func loadScenarios() {
	scenarios = []scenario{
		{
			name:     "grpc service non-matching request configured [should pass] [real world client]",
			service:  "should.pass",
			endpoint: "irrelevant",
			error:    false,
			cw: &clientWrapper{
				Client:     realWorldClient(),
				policyPath: fname,
			},
		},
		{
			name:     "grpc service non-matching request configured [should pass] [isolated client]",
			service:  "should.pass",
			endpoint: "irrelevant",
			error:    false,
			cw: &clientWrapper{
				Client:     isolatedClient(),
				policyPath: fname,
			},
		},
		{
			name:     "grpc service matching request configured [should fail] [real world client]",
			service:  "com.owncloud.api.thumbnails",
			endpoint: "irrelevant",
			error:    true,
			cw: &clientWrapper{
				Client:     realWorldClient(),
				policyPath: fname,
			},
		},
		{
			name:     "grpc service matching request configured [should fail] [isolated client]",
			service:  "com.owncloud.api.thumbnails",
			endpoint: "irrelevant",
			error:    true,
			cw: &clientWrapper{
				Client:     isolatedClient(),
				policyPath: fname,
			},
		},
	}
}

func BenchmarkCheckPolicySuccess(b *testing.B) {
	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if err := scenario.cw.checkPolicy(context.Background(), client.NewRequest(scenario.service, scenario.endpoint, nil)); err != nil {
					if !scenario.error {
						b.Error(err)
					}
				}
			}
		})
	}
}

// FIXTURES

// get as close as a real world scenario by using both wrappers
func realWorldClient() client.Client {
	return mgrpcc.NewClient(
		client.Wrap(mbreaker.NewClientWrapper()),
		client.Wrap(NewClientWrapper()),
	)
}

// create a client with only the policy middleware
func isolatedClient() client.Client {
	return mgrpcc.NewClient(
		client.Wrap(NewClientWrapper()),
	)
}
