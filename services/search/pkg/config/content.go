package config

// Extractor defines which extractor to use
type Extractor struct {
	Type             string        `yaml:"type" env:"SEARCH_EXTRACTOR_TYPE" desc:"Defines the content extraction engine. Defaults to 'basic'. Supported values are: 'basic' and 'tika'." introductionVersion:"pre5.0"`
	CS3AllowInsecure bool          `yaml:"cs3_allow_insecure" env:"OCIS_INSECURE;SEARCH_EXTRACTOR_CS3SOURCE_INSECURE" desc:"Ignore untrusted SSL certificates when connecting to the CS3 source." introductionVersion:"pre5.0"`
	Tika             ExtractorTika `yaml:"tika"`
}

// ExtractorTika configures the Tika extractor
type ExtractorTika struct {
	TikaURL        string `yaml:"tika_url" env:"SEARCH_EXTRACTOR_TIKA_TIKA_URL" desc:"URL of the tika server." introductionVersion:"pre5.0"`
	CleanStopWords bool   `yaml:"clean_stop_words" env:"SEARCH_EXTRACTOR_TIKA_CLEAN_STOP_WORDS" desc:"Defines if stop words should be cleaned or not. See the documentation for more details." introductionVersion:"5.0"`
}
