package tags

import (
	"encoding/json"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/tags"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"go-micro.dev/v4/metadata"
	"net/http"
)

type tagsService struct {
	r            chi.Router
	logger       log.Logger
	searchClient searchsvc.SearchProviderService
}

// NewTagsService bootstraps the tags service
func NewTagsService(r chi.Router, searchClient searchsvc.SearchProviderService, logger log.Logger) {
	svc := tagsService{
		logger:       logger,
		searchClient: searchClient,
	}

	r.Get("/tags", svc.GetTags)
}

// GetTags lists all available tags as json response.
func (s *tagsService) GetTags(w http.ResponseWriter, r *http.Request) {
	th := r.Header.Get(revactx.TokenHeader)
	ctx := revactx.ContextSetToken(r.Context(), th)
	ctx = metadata.Set(ctx, revactx.TokenHeader, th)
	sr, err := s.searchClient.Search(ctx, &searchsvc.SearchRequest{
		Query:    "Tags:*",
		PageSize: -1,
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("Could not search for tags")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t := tags.FromList("")
	for _, match := range sr.Matches {
		for _, tag := range match.Entity.Tags {
			t.AddList(tag)
		}
	}

	jm, err := json.Marshal(struct {
		Tags []string `json:"tags"`
	}{
		Tags: t.AsSlice(),
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("Could not read tags")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(jm); err != nil {
		s.logger.Error().Err(err).Msg("Could not write tags")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
