package config

// Activities defines the available activities configuration.
type Activities struct {
	Enabled bool              `yaml:"enabled" env:"EXPERIMENTAL_ACTIVITIES_ENABLED" desc:"enable activities app."`
	Storage ActivitiesStorage `yaml:"storage" env:"EXPERIMENTAL_ACTIVITIES_STORAGE" desc:"activities storage."`
}

// ActivitiesStorage defines the available activities storage configuration.
type ActivitiesStorage struct {
	Type     string               `yaml:"type" env:"EXPERIMENTAL_ACTIVITIES_STORAGE_TYPE" desc:"Defines the activity storage type."`
	MemStore ActivitiesMemStorage `yaml:"mem_storage"`
}

// ActivitiesMemStorage defines the available activities, mem_storage configuration.
type ActivitiesMemStorage struct {
	Capacity uint64 `yaml:"capacity" env:"EXPERIMENTAL_ACTIVITIES_STORE_MEM_STORAGE_KEEP" desc:"specifies how many activities should be saved across all users."`
}
