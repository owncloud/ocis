package config

var config = map[string]string{
	"bin":           "/usr/bin/ocis",
	"url":           "https://localhost:9200",
	"retry":         "5",
	"adminUsername": "",
	"adminPassword": "",
	"namespace":     "ocis",
}

var debugPorts = map[string]int{
	"ocis":               9250,
	"activitylog":        9197,
	"antivirus":          9277,
	"app-registry":       9243,
	"app-provider":       9165,
	"audit":              9229,
	"auth-app":           9245,
	"auth-basic":         9147,
	"auth-bearer":        9149,
	"auth-machine":       9167,
	"auth-service":       9198,
	"clientlog":          9260,
	"collaboration":      9304,
	"eventhistory":       9270,
	"frontend":           9141,
	"gateway":            9143,
	"graph":              9124,
	"groups":             9161,
	"idm":                9239,
	"idp":                9134,
	"invitations":        9269,
	"nats":               9234,
	"notifications":      9174,
	"ocdav":              9163,
	"ocm":                9281,
	"ocs":                9114,
	"policies":           9129,
	"postprocessing":     9255,
	"proxy":              9205,
	"search":             9224,
	"settings":           9194,
	"sharing":            9151,
	"sse":                9139,
	"storage-publiclink": 9179,
	"storage-shares":     9156,
	"storage-system":     9217,
	"storage-users":      9159,
	"thumbnails":         9189,
	"userlog":            9214,
	"users":              9145,
	"web":                9104,
	"webdav":             9119,
	"webfinger":          9279,
}

// only some services have gRPC ports
var grpcPorts = map[string]int{
	"app-registry":       9242,
	"app-provider":       9164,
	"auth-app":           9246,
	"auth-basic":         9146,
	"auth-bearer":        9148,
	"auth-machine":       9166,
	"auth-service":       9616, // 9199
	"collaboration":      9301,
	"eventhistory":       8080, // 9274
	"gateway":            9142,
	"groups":             9160,
	"ocm":                9282,
	"policies":           9125,
	"search":             9220,
	"settings":           9191,
	"sharing":            9150,
	"storage-publiclink": 9178,
	"storage-shares":     9154,
	"storage-system":     9215,
	"storage-users":      9157,
	"thumbnails":         9185,
	"users":              9144,
}

func Set(key string, value string) {
	config[key] = value
}

func Get(key string) string {
	return config[key]
}

func SetServiceDebugPort(key string, value int) {
	debugPorts[key] = value
}

func GetServiceDebugPort(key string) int {
	return debugPorts[key]
}

func SetServiceGRPCPort(key string, value int) {
	grpcPorts[key] = value
}

func GetServiceGRPCPort(key string) int {
	return grpcPorts[key]
}
