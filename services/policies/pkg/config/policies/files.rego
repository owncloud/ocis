package files

import future.keywords.if
import data.utils

default granted = true

ALLOWED_FILE_EXTENSIONS := [
    ".apk", ".avi", ".bat", ".bmp", ".css", ".csv", ".doc", ".docm", ".docx",
    ".docxf", ".dotx", ".eml", ".epub", ".htm", ".html", ".ipa", ".jar", ".java",
    ".jpg", ".js", ".json", ".mp3", ".mp4", ".msg", ".odp", ".ods", ".odt", ".oform",
    ".ots", ".ott", ".pdf", ".php", ".png", ".potm", ".potx", ".ppsm", ".ppsx", ".ppt",
    ".pptm", ".pptx", ".py", ".rtf", ".sb3", ".sprite3", ".sql", ".svg", ".tif", ".tiff",
    ".txt", ".xls", ".xlsm", ".xlsx", ".xltm", ".xltx", ".xml", ".zip", ".md"
]

granted := false if {
    utils.is_request_type_put
    not utils.collection_contains(ALLOWED_FILE_EXTENSIONS, input.request.path)
}

granted := false if {
    utils.is_stage_pp
    not utils.collection_contains(ALLOWED_FILE_EXTENSIONS, input.resource.name)

    #bytes := ocis_get_resource(input.resource.url)
    #mimetype := ocis_get_mimetype(bytes)
    #not utils.collection_contains(["image/png"], mimetype)
}
