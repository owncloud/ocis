package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/owncloud/ocis/pkg/config"
)

type callbackEndpoint struct {
	server         *http.Server
	code           string
	shutdownSignal chan string
}

func (h *callbackEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	if code != "" {
		h.code = code
		fmt.Fprintln(w, "Login is successful, You may close the browser and goto commandline")
	} else {
		fmt.Fprintln(w, "Login is not successful, You may close the browser and try again")
	}
	h.shutdownSignal <- "shutdown"
}

// HandleOpenIDFlow opens a local callback endpoint and opens a browser to initiate the oidc authorization code flow
func HandleOpenIDFlow(cfg config.OIDC) {
	callbackEndpoint := &callbackEndpoint{}
	callbackEndpoint.shutdownSignal = make(chan string)
	server := &http.Server{
		Addr:           cfg.CallbackAddr,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	callbackEndpoint.server = server
	http.Handle("/callback", callbackEndpoint)
	authzURL, authzURLParseError := url.Parse(cfg.AuthEndpoint)

	if authzURLParseError != nil {
		log.Fatal(authzURLParseError)
	}
	query := authzURL.Query()
	query.Set("response_type", "code")
	query.Set("scope", "openid profile email offline")
	query.Set("client_id", cfg.ClientID)
	query.Set("redirect_uri", cfg.CallbackURL)
	uuid := uuid.Must(uuid.NewV4())
	query.Set("state", uuid.String())
	authzURL.RawQuery = query.Encode()

	cmd := open(authzURL.String())
	cmdErorr := cmd.Start()
	if cmdErorr != nil {
		log.Fatal(authzURLParseError)
	}

	go func() {
		server.ListenAndServe()
	}()

	<-callbackEndpoint.shutdownSignal
	callbackEndpoint.server.Shutdown(context.Background())
	log.Println("Authorization code is ", callbackEndpoint.code)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	vals := url.Values{}
	vals.Set("grant_type", "authorization_code")
	vals.Set("code", callbackEndpoint.code)
	vals.Set("redirect_uri", cfg.CallbackURL)
	req, requestError := http.NewRequest("POST", cfg.TokenEndpoint, strings.NewReader(vals.Encode()))
	if requestError != nil {
		log.Fatal(requestError)
	}
	req.SetBasicAuth(cfg.ClientID, cfg.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, clientError := client.Do(req)
	if clientError != nil {
		log.Fatal(clientError)
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result != nil {
		jsonStr, marshalError := json.Marshal(result)
		if marshalError != nil {
			log.Fatal(marshalError)
		}
		log.Println(string(jsonStr))
	} else {
		log.Println("Error while getting ID token")
	}

}
