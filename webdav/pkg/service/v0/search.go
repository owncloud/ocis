package svc

import (
	"net/http"

	searchsvc "github.com/owncloud/ocis/protogen/gen/ocis/services/search/v0"
	merrors "go-micro.dev/v4/errors"
)

// Search is the endpoint for retrieving search results for REPORT requests
func (g Webdav) Search(w http.ResponseWriter, r *http.Request) {

	rsp, err := g.searchClient.Search(r.Context(), &searchsvc.SearchRequest{})
	if err != nil {
		e := merrors.Parse(err.Error())
		switch e.Code {
		case http.StatusBadRequest:
			renderError(w, r, errBadRequest(err.Error()))
		default:
			renderError(w, r, errInternalError(err.Error()))
		}
		g.log.Error().Err(err).Msg("could not get search results")
		return
	}

	g.sendSearchResponse(rsp, w, r)
}

func (g Webdav) sendSearchResponse(rsp *searchsvc.SearchResponse, w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)

}
