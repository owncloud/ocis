package config

type Engine struct {
	Type  string      `yaml:"type" env:"SEARCH_ENGINE_TYPE" desc:"Defines which search engine to use."`
	Bleve EngineBleve `yaml:"bleve"`
}

type EngineBleve struct {
	Datapath string `yaml:"data_path" env:"SEARCH_ENGINE_BLEVE_DATA_PATH" desc:"Path for the search persistence directory."`
}
