package ocis

default deny = false
default maxCount = 10

deny {
    input.service == "com.owncloud.api.thumbnails"
    input.external.users_count > maxCount
}

# deny {
#    input.service == "com.owncloud.api.settings"
#    input.endpoint == "AccountsService.ListAccounts"
#    input.method == "RoleService.ListRoleAssignments"
##    input.standard_claims.email == "admin@example.org"
##    input.standard_claims.groups == ""
##    input.standard_claims.iss == "https://localhost:9200"
##    input.standard_claims.name == "admin"
# }
