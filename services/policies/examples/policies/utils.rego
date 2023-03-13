package utils

import future.keywords.if

ALLOWED_FILE_EXTENSIONS := [
    ".apk", ".avi", ".bat", ".bmp", ".css", ".csv", ".doc", ".docm", ".docx",
    ".docxf", ".dotx", ".eml", ".epub", ".htm", ".html", ".ipa", ".jar", ".java",
    ".jpg", ".js", ".json", ".mp3", ".mp4", ".msg", ".odp", ".ods", ".odt", ".oform",
    ".ots", ".ott", ".pdf", ".php", ".png", ".potm", ".potx", ".ppsm", ".ppsx", ".ppt",
    ".pptm", ".pptx", ".py", ".rtf", ".sb3", ".sprite3", ".sql", ".svg", ".tif", ".tiff",
    ".txt", ".xls", ".xlsm", ".xlsx", ".xltm", ".xltx", ".xml", ".zip", ".md"
]

##

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

is_request_path_file {
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

