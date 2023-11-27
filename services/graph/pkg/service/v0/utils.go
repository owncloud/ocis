package svc

import (
	"encoding/json"
	"io"
	"net/http"

	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
)

// StrictJSONUnmarshal is a wrapper around json.Unmarshal that returns an error if the json contains unknown fields.
func StrictJSONUnmarshal(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// GetDriveAndItemIDParam parses the driveID and itemID from the request,
// validates the common fields and returns the parsed IDs if ok.
func (g Graph) GetDriveAndItemIDParam(w http.ResponseWriter, r *http.Request) (storageprovider.ResourceId, storageprovider.ResourceId, bool) {
	empty := storageprovider.ResourceId{}

	driveID, err := parseIDParam(r, "driveID")
	if err != nil {
		g.logger.Debug().Err(err).Msg("could not parse driveID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid driveID")
		return empty, empty, false
	}

	itemID, err := parseIDParam(r, "itemID")
	if err != nil {
		g.logger.Debug().Err(err).Msg("could not parse itemID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid itemID")
		return empty, empty, false
	}

	if itemID.GetOpaqueId() == "" {
		g.logger.Debug().Interface("driveID", driveID).Interface("itemID", itemID).Msg("empty item opaqueID")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "invalid itemID")
		return empty, empty, false
	}

	if driveID.GetStorageId() != itemID.GetStorageId() || driveID.GetSpaceId() != itemID.GetSpaceId() {
		g.logger.Debug().Interface("driveID", driveID).Interface("itemID", itemID).Msg("driveID and itemID do not match")
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, "driveID and itemID do not match")
		return empty, empty, false
	}

	return driveID, itemID, true
}
