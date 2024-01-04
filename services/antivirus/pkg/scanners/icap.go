package scanners

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"

	"github.com/cs3org/reva/v2/pkg/mime"
	ic "github.com/egirna/icap-client"
)

// NewICAP returns a Scanner talking to an ICAP server
func NewICAP(icapURL string, icapService string, timeout time.Duration) (ICAP, error) {
	endpoint, err := url.Parse(icapURL)
	if err != nil {
		return ICAP{}, err
	}

	endpoint.Scheme = "icap"
	endpoint.Path = icapService

	client, err := ic.NewClient(ic.Options{
		Timeout: timeout,
	})

	return ICAP{client: client, url: *endpoint}, nil
}

// ICAP is responsible for scanning files using an ICAP server
type ICAP struct {
	client *ic.Client
	url    url.URL
}

// Scan scans a file using the ICAP server
func (s ICAP) Scan(in Input) (Result, error) {
	ctx := context.TODO()
	result := Result{}

	httpReq, err := http.NewRequest(http.MethodPost, in.Url, in.Body)
	if err != nil {
		return result, err
	}

	httpReq.ContentLength = in.Size
	if mt := mime.Detect(path.Ext(in.Name) == "", in.Name); mt != "" {
		httpReq.Header.Set("Content-Type", mt)
	}

	optReq, err := ic.NewRequest(ctx, ic.MethodOPTIONS, s.url.String(), nil, nil)
	if err != nil {
		return result, err
	}

	optRes, err := s.client.Do(optReq)
	if err != nil {
		return result, err
	}

	req, err := ic.NewRequest(ctx, ic.MethodREQMOD, s.url.String(), httpReq, nil)
	if err != nil {
		return result, err
	}

	err = req.SetPreview(optRes.PreviewBytes)
	if err != nil {
		return result, err
	}

	res, err := s.client.Do(req)
	if err != nil {
		return result, err
	}

	if data, infected := res.Header["X-Infection-Found"]; infected {
		result.Infected = infected

		match := regexp.MustCompile(`Threat=(.*);`).FindStringSubmatch(fmt.Sprint(data))

		if len(match) > 1 {
			result.Description = match[1]
		}
	}

	return result, nil
}
