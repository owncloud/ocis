package proxy

import future.keywords.if
import data.utils

default granted := true

granted = false if {
    print("PRINT MESSAGE EXAMPLE")
    input.request.method == "PUT"
    not startswith(input.request.path, "/ocs")
    not utils.is_extension_allowed(input.request.path)
}

granted = false if {
    input.request.method == "POST"
    startswith(input.request.path, "/remote.php")
    not utils.is_extension_allowed(input.resource.name)
}
