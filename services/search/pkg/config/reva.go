package config

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `ocisConfig:"address" env:"OCIS_REVA_GATEWAY;REVA_GATEWAY" desc:"The CS3 gateway endpoint." deprecationVersion:"3.0" removalVersion:"4.0.0" deprecationInfo:"REVA_GATEWAY changing name for consistency" deprecationReplacement:"OCIS_REVA_GATEWAY"`
}
