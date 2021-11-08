package shared

// EnvBinding represents a direct binding from an env variable to a go kind. Along with gookit/config, its primal goal
// is to unpack environment variables into a Go value. We do so with reflection, and this data structure is just a step
// in between.
type EnvBinding struct {
	EnvVars     []string    // name of the environment var.
	Destination interface{} // pointer to the original config value to modify.
}

// Log defines the available logging configuration.
type Log struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
	Color  bool   `mapstructure:"color"`
	File   string `mapstructure:"file"`
}
