package config

// Engine defines which search engine to use
type Engine struct {
	Type  string      `yaml:"type" env:"SEARCH_ENGINE_TYPE" desc:"Defines which search engine to use. Defaults to 'bleve'. Supported values are: 'bleve'." introductionVersion:"pre5.0"`
	Bleve EngineBleve `yaml:"bleve"`
}

// EngineBleve configures the bleve engine
type EngineBleve struct {
	Datapath string `yaml:"data_path" env:"SEARCH_ENGINE_BLEVE_DATA_PATH" desc:"The directory where the filesystem will store search data. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/search." introductionVersion:"pre5.0"`
	Scale    bool   `yaml:"scale" env:"SEARCH_ENGINE_BLEVE_SCALE" desc:"Enable scaling of the search index (bleve). If set to 'true', the instance of the search service will no longer have exclusive write access to the index. Note when scaling search, all instances of the search service must be set to true! For 'false', which is the default, the running search service has exclusive access to the index as long it is running. This locks out other search processes tying to access the index." introductionVersion:"7.2.0"`
}
