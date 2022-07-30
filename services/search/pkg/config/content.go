package config

type Extractor struct {
	Type             string        `yaml:"type" env:"SEARCH_EXTRACTOR_TYPE" desc:"Defines the content extraction engine."`
	CS3AllowInsecure bool          `yaml:"cs3_allow_insecure" env:"OCIS_INSECURE;SEARCH_EXTRACTOR_CS3SOURCE_INSECURE" desc:"Ignore untrusted SSL certificates when connecting to the CS3 source."`
	Tika             ExtractorTika `yaml:"tika"`
}

type ExtractorTika struct {
	TikaURL string `yaml:"tika_url" env:"SEARCH_EXTRACTOR_TIKA_TIKA_URL" desc:"URL of the tika server."`
}
