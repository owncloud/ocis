package postprocessing

import future.keywords.if
import data.utils

default granted := true

granted = false if {
    not utils.collection_contains(utils.ALLOWED_FILE_EXTENSIONS, input.resource.name)
}

granted = false if {
    bytes := ocis.resource.download(input.resource.url)
    mimetype := ocis.mimetype.detect(bytes)
    extensions := ocis.mimetype.extension_for_mimetype(mimetype)

    extension := extensions[_]
    not utils.collection_contains(utils.ALLOWED_FILE_EXTENSIONS, extension)
}
