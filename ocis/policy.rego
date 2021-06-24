package ocis

default deny = false

deny {
    input.service == "com.owncloud.api.thumbnails"
}

# deny {
#    input.service == "com.owncloud.api.settings"
#    input.endpoint == "AccountsService.ListAccounts"
#    input.method == "RoleService.ListRoleAssignments"
#    input.standard_claims.email == "admin@example.org"
#    input.standard_claims.groups == ""
#    input.standard_claims.iss == "https://localhost:9200"
#    input.standard_claims.name == "admin"
# }
