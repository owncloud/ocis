package ocis.authz

default allow = true

# disable txt updates
allow := false {
    input.method == "PUT"
    endswith(input.path, ".txt")
}

# disable pdf uploads, expect for admin
allow := false {
    input.method == "POST"
    input.user.username != "admin"
    hasMimetype(input.body, "application/pdf")
}

