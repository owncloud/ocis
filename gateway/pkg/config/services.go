package config

// ServiceMap holds the addresses of other services
type ServiceMap struct {
	// TODO: replace all these static addresses by a registry

	// registries is located on the gateway
	AuthRegistryAddr    string
	StorageRegistryAddr string
	AppRegistryAddr     string

	// user metadata is located on the users services
	PreferenceAddr    string
	UserProviderAddr  string
	GroupProviderAddr string

	// sharing is located on the sharing service
	UserShareProviderAddr   string
	PublicShareProviderAddr string
	OCMShareProviderAddr    string

	StorageHomeAddr    string
	StorageUsersAddr   string
	StoragePublicShare string

	AuthBasicAddr        string
	AuthBearerAddr       string
	AuthMachineAddr      string
	AuthPublicSharesAddr string
}
