package TestHelpers

var adminUsername string = "admin"
var adminPassword string = "admin"
var baseUrl string = "https://localhost:9200"

func GetAdminUsername() string {
	return adminUsername
}

func GetAdminPassword() string {
	return adminPassword
}

func GetBaseUrl() string {
	return baseUrl
}
