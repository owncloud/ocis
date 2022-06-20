package svc

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0/response"
)

// ListUserGroups lists a users groups
func (o Ocs) ListUserGroups(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")
	userid, _ = url.PathUnescape(userid)
	switch o.config.AccountBackend {
	case "cs3":
		// TODO
		o.mustRender(w, r, response.DataRender(&data.Groups{}))
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
	return
}

// AddToGroup adds a user to a group
func (o Ocs) AddToGroup(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		// TODO
		o.cs3WriteNotSupported(w, r)
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
}

// RemoveFromGroup removes a user from a group
func (o Ocs) RemoveFromGroup(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		// TODO
		o.cs3WriteNotSupported(w, r)
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
}

// ListGroups lists all groups
func (o Ocs) ListGroups(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		// TODO
		o.mustRender(w, r, response.DataRender(&data.Groups{}))
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
	return
}

// AddGroup adds a group
// oC10 implementation: https://github.com/owncloud/core/blob/762780a23c9eadda4fb5fa8db99eba66a5100b6e/apps/provisioning_api/lib/Groups.php#L126-L154
func (o Ocs) AddGroup(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		o.cs3WriteNotSupported(w, r)
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
}

// DeleteGroup deletes a group
func (o Ocs) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		o.cs3WriteNotSupported(w, r)
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
}

// GetGroupMembers lists all members of a group
func (o Ocs) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		// TODO
		o.mustRender(w, r, response.DataRender(&data.Users{}))
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
	return
}
