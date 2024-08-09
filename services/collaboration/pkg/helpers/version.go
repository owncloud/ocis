package helpers

import (
	"net/http"
	"strconv"

	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
)

// SetVersionHeader sets a WOPI version header on the response writer
func SetVersionHeader(w http.ResponseWriter, t *typesv1beta1.Timestamp) {
	// non-canonical headers can only be set directly on the header map
	w.Header().Set("X-WOPI-ItemVersion", GetVersion(t))
}

// GetVersion returns a string representation of the timestamp
func GetVersion(timestamp *typesv1beta1.Timestamp) string {
	return "v" + strconv.FormatUint(timestamp.GetSeconds(), 10) +
		strconv.FormatUint(uint64(timestamp.GetNanos()), 10)
}
