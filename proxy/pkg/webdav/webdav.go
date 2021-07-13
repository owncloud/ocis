package webdav

import "net/http"

var methods = []string{"PROPFIND", "DELETE", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK"}

// This is a non exhaustive way to detect if a request is directed to a webdav server. This naïve implementation
// only deals with the set of methods exclusive to WebDAV. Since WebDAV is a superset of HTTP, GET, POST and so on
// are valid methods, but this implementation would require a larger effort than we can build upon in this file.
// This is needed because the proxy might need to create a response with a webdav body; such as unauthorized.
func IsWebdavRequest(r *http.Request) bool {
	for i := range methods {
		if methods[i] == r.Method {
			return true
		}
	}
	return false
}
