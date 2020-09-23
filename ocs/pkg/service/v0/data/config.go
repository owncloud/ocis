package data

// ConfigData holds basic config
type ConfigData struct {
	Version string `json:"version" xml:"version"`
	Website string `json:"website" xml:"website"`
	Host    string `json:"host" xml:"host"`
	Contact string `json:"contact" xml:"contact"`
	SSL     string `json:"ssl" xml:"ssl"`
}
