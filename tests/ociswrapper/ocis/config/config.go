package config

var config = map[string]string{
	"bin":           "/usr/bin/ocis",
	"url":           "https://localhost:9200",
	"retry":         "5",
	"adminUsername": "",
	"adminPassword": "",
}

var services = map[string]int{
	"ocis":               9250,
	"activitylog":        9197,
	"app-provider":       9165,
	"app-registry":       9243,
	"audit":              9229,
	"auth-app":           9245,
	"auth-bearer":        9149,
	"auth-basic":         9147,
	"auth-machine":       9167,
	"auth-service":       9198,
	"clientlog":          9260,
	"eventhistory":       9270,
	"frontend":           9141,
	"gateway":            9143,
	"graph":              9124,
	"groups":             9161,
	"idm":                9239,
	"idp":                9134,
	"invitations":        9269,
	"nats":               9234,
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

func Set(key string, value string) {
	config[key] = value
}

func Get(key string) string {
	return config[key]
}

func SetServiceDebugPort(key string, value int) {
	services[key] = value
}

func GetServiceDebugPort(key string) int {
	return services[key]
}
