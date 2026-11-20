package config

import "time"

// Extractor defines which extractor to use
type Extractor struct {
	Type             string        `yaml:"type" env:"SEARCH_EXTRACTOR_TYPE" desc:"Defines the content extraction engine. Defaults to 'basic'. Supported values are: 'basic' and 'tika'." introductionVersion:"pre5.0"`
	CS3AllowInsecure bool          `yaml:"cs3_allow_insecure" env:"OCIS_INSECURE;SEARCH_EXTRACTOR_CS3SOURCE_INSECURE" desc:"Ignore untrusted SSL certificates when connecting to the CS3 source." introductionVersion:"pre5.0"`
	Tika             ExtractorTika `yaml:"tika"`

	MaxConsecutiveFailures int `yaml:"max_consecutive_failures" env:"SEARCH_EXTRACTOR_MAX_CONSECUTIVE_FAILURES" desc:"Abort a space (re)index walk after this many files in a row fail content extraction, which usually means the extraction service is down. Files that cannot be resolved at all are logged and skipped without counting towards this limit. Set to 0 to never abort." introductionVersion:"9.0.0"`

	FailedRetryInterval time.Duration `yaml:"failed_retry_interval" env:"SEARCH_EXTRACTOR_FAILED_RETRY_INTERVAL" desc:"How long a space (re)index walk leaves a file alone after its content extraction failed, as long as the file has not changed. Until then the file is indexed by its name and metadata only. A changed file is always extracted again. For example '24h'. Set to 0 to retry such files on every walk. See the Environment Variable Types description for more details." introductionVersion:"9.0.0"`
}

// ExtractorTika configures the Tika extractor
type ExtractorTika struct {
	TikaURL        string `yaml:"tika_url" env:"SEARCH_EXTRACTOR_TIKA_TIKA_URL" desc:"URL of the tika server." introductionVersion:"pre5.0"`
	CleanStopWords bool   `yaml:"clean_stop_words" env:"SEARCH_EXTRACTOR_TIKA_CLEAN_STOP_WORDS" desc:"Defines if stop words should be cleaned or not. See the documentation for more details." introductionVersion:"5.0"`
}
