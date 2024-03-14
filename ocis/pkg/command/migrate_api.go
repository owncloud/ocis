package command

import (
	"encoding/json"
	"errors"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
)

type TokenExchangeRequest struct {
	EMail string `json:"email"`
}
type TokenExchangeResponse struct {
	AccessToken string `json:"accessToken"`
}
type httpHandler struct {
	SharedSecret string
	Logger       *zerolog.Logger
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		w.WriteHeader(400)
		return
	}
	if auth != "Bearer "+h.SharedSecret {
		w.WriteHeader(400)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	request := TokenExchangeRequest{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	// TODO: create app token with some lifetime - up for a different PR

	// set access token back
	response := TokenExchangeResponse{AccessToken: "foobar"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func StartMigrationAPI(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "start-migration-api",
		Usage: "starts the migration api service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "shared_secret",
				Value:    "string",
				Usage:    "shared secret to gain api access",
				Required: true,
			},
		},
		Before: func(c *cli.Context) error {
			sharedSecret := c.String("shared_secret")
			if len(sharedSecret) < 32 {
				return errors.New("shared secret is too short. Minimum length is 32 characters")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			log := logger()
			// ctx := log.WithContext(context.Background())

			sharedSecret := c.String("shared_secret")
			log.Info().Str("Secret", sharedSecret).Msg("User input")

			mux := http.NewServeMux()
			// API POST /satellites/tokenExchange with POST body {"email": "user1@example.com"}
			h := &httpHandler{
				SharedSecret: sharedSecret,
				Logger:       log,
			}
			mux.Handle("/satellites/tokenExchange", h)

			return http.ListenAndServe(":9999", mux)
		},
	}
}
