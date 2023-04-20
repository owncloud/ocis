package keyfunc

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	// defaultRefreshTimeout is the default duration for the context used to create the HTTP request for a refresh of
	// the JWKS.
	defaultRefreshTimeout = time.Minute
)

// Get loads the JWKS at the given URL.
func Get(jwksURL string, options Options) (jwks *JWKS, err error) {
	jwks = &JWKS{
		jwksURL: jwksURL,
	}

	applyOptions(jwks, options)

	if jwks.client == nil {
		jwks.client = http.DefaultClient
	}
	if jwks.requestFactory == nil {
		jwks.requestFactory = defaultRequestFactory
	}
	if jwks.responseExtractor == nil {
		jwks.responseExtractor = ResponseExtractorStatusOK
	}
	if jwks.refreshTimeout == 0 {
		jwks.refreshTimeout = defaultRefreshTimeout
	}
	if !options.JWKUseNoWhitelist && len(jwks.jwkUseWhitelist) == 0 {
		jwks.jwkUseWhitelist = map[JWKUse]struct{}{
			UseOmitted:   {},
			UseSignature: {},
		}
	}

	err = jwks.refresh()
	if err != nil {
		return nil, err
	}

	if jwks.refreshInterval != 0 || jwks.refreshUnknownKID {
		jwks.ctx, jwks.cancel = context.WithCancel(context.Background())
		jwks.refreshRequests = make(chan context.CancelFunc, 1)
		go jwks.backgroundRefresh()
	}

	return jwks, nil
}

// backgroundRefresh is meant to be a separate goroutine that will update the keys in a JWKS over a given interval of
// time.
func (j *JWKS) backgroundRefresh() {
	var lastRefresh time.Time
	var queueOnce sync.Once
	var refreshMux sync.Mutex
	if j.refreshRateLimit != 0 {
		lastRefresh = time.Now().Add(-j.refreshRateLimit)
	}

	// Create a channel that will never send anything unless there is a refresh interval.
	refreshInterval := make(<-chan time.Time)

	// Enter an infinite loop that ends when the background ends.
	for {
		if j.refreshInterval != 0 {
			refreshInterval = time.After(j.refreshInterval)
		}

		select {
		case <-refreshInterval:
			select {
			case <-j.ctx.Done():
				return
			case j.refreshRequests <- func() {}:
			default: // If the j.refreshRequests channel is full, don't send another request.
			}

		case cancel := <-j.refreshRequests:
			refreshMux.Lock()
			if j.refreshRateLimit != 0 && lastRefresh.Add(j.refreshRateLimit).After(time.Now()) {
				// Don't make the JWT parsing goroutine wait for the JWKS to refresh.
				cancel()

				// Launch a goroutine that will get a reservation for a JWKS refresh or fail to and immediately return.
				queueOnce.Do(func() {
					go func() {
						refreshMux.Lock()
						wait := time.Until(lastRefresh.Add(j.refreshRateLimit))
						refreshMux.Unlock()
						select {
						case <-j.ctx.Done():
							return
						case <-time.After(wait):
						}

						refreshMux.Lock()
						defer refreshMux.Unlock()
						err := j.refresh()
						if err != nil && j.refreshErrorHandler != nil {
							j.refreshErrorHandler(err)
						}

						lastRefresh = time.Now()
						queueOnce = sync.Once{}
					}()
				})
			} else {
				err := j.refresh()
				if err != nil && j.refreshErrorHandler != nil {
					j.refreshErrorHandler(err)
				}

				lastRefresh = time.Now()

				// Allow the JWT parsing goroutine to continue with the refreshed JWKS.
				cancel()
			}
			refreshMux.Unlock()

		// Clean up this goroutine when its context expires.
		case <-j.ctx.Done():
			return
		}
	}
}

func defaultRequestFactory(ctx context.Context, url string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewReader(nil))
}

// refresh does an HTTP GET on the JWKS URL to rebuild the JWKS.
func (j *JWKS) refresh() (err error) {
	var ctx context.Context
	var cancel context.CancelFunc
	if j.ctx != nil {
		ctx, cancel = context.WithTimeout(j.ctx, j.refreshTimeout)
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), j.refreshTimeout)
	}
	defer cancel()

	req, err := j.requestFactory(ctx, j.jwksURL)
	if err != nil {
		return fmt.Errorf("failed to create request via factory function: %w", err)
	}

	resp, err := j.client.Do(req)
	if err != nil {
		return err
	}

	jwksBytes, err := j.responseExtractor(ctx, resp)
	if err != nil {
		return fmt.Errorf("failed to extract response via extractor function: %w", err)
	}

	// Only reprocess if the JWKS has changed.
	if len(jwksBytes) != 0 && bytes.Equal(jwksBytes, j.raw) {
		return nil
	}
	j.raw = jwksBytes

	updated, err := NewJSON(jwksBytes)
	if err != nil {
		return err
	}

	j.mux.Lock()
	defer j.mux.Unlock()
	j.keys = updated.keys

	if j.givenKeys != nil {
		for kid, key := range j.givenKeys {
			// Only overwrite the key if configured to do so.
			if !j.givenKIDOverride {
				if _, ok := j.keys[kid]; ok {
					continue
				}
			}

			j.keys[kid] = parsedJWK{public: key.inter}
		}
	}

	return nil
}
