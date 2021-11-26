package svc

import (
	"net/http"

	"github.com/owncloud/ocis/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/response"
)

// GetConfig renders the ocs config endpoint
func (o Ocs) GetConfig(w http.ResponseWriter, r *http.Request) {
	o.mustRender(w, r, response.DataRender(&data.ConfigData{
		Version: "1.7",  // TODO get from env
		Website: "ocis", // TODO get from env
		Host:    "",     // TODO get from FRONTEND config
		Contact: "",     // TODO get from env
		SSL:     "true", // TODO get from env
	}))
}
