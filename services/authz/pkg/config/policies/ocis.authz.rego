package ocis.authz

import future.keywords.if

default allow = true

is_stage_http if "http" == input.stage

is_stage_pp if "pp" == input.stage

is_admin if "admin" == input.user.username

is_http_put if {
    is_stage_http

    input.method == "PUT"
}


# disallow txt files for all users except admin
allow := false if {
    not is_admin

    is_http_put

    endswith(input.name, ".txt")
}

# disallow png files for all users except admin
allow := false if {
    not is_admin

    is_stage_pp

    dataB := loadResource(input.url)

    hasMimetype(dataB, "image/png")
}

# disallow files with content `voldemort` for all users except admin
allow := false if {
    not is_admin

    is_stage_pp

    dataB := loadResource(input.url)
    dataS := convertBtoS(dataB)

    contains(dataS, "voldemort")
}

