package postprocessing

import future.keywords.if
import data.utils

default granted := true

granted = false if {
    not utils.is_extension_allowed(input.resource.name)
}

granted = false if {
    bytes := ocis.resource.download(input.resource.url)
    mimetype := ocis.mimetype.detect(bytes)

    not utils.is_mimetype_allowed(mimetype)
}
