package utils

import future.keywords.if

is_stage_http {
    input.stage == "http"
}

is_stage_pp {
    input.stage == "pp"
}

##

is_user_admin {
    input.user.username == "admin"
}

##

is_request_type_put {
    is_stage_http
    input.request.method == "PUT"
}

is_request_type_mkcol {
    is_stage_http
    input.request.method == "MKCOL"
}

##

collection_contains(collection, source) {
     current := collection[_]
     endswith(source, current)
}

