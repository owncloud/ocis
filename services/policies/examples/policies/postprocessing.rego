package postprocessing

import future.keywords.if
import data.utils

default granted = true

granted := false if {
    not utils.collection_contains(utils.ALLOWED_FILE_EXTENSIONS, input.resource.name)
}
