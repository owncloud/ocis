package scanners

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

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

	return ICAP{
		client: &ic.Client{
			Timeout: timeout,
		},
		endpoint: endpoint.String(),
	}, nil
}

// ICAP is a Scanner talking to an ICAP server
type ICAP struct {
	client   *ic.Client
	endpoint string
}

// Scan to fulfill Scanner interface
func (s ICAP) Scan(file io.Reader) (ScanResult, error) {
	sr := ScanResult{}

	httpReq, err := http.NewRequest(http.MethodGet, "http://localhost", file)
	if err != nil {
		return sr, err
	}

	req, err := ic.NewRequest(ic.MethodREQMOD, s.endpoint, httpReq, nil)
	if err != nil {
		return sr, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return sr, err
	}

	// TODO: make header configurable. See oc10 documentation: https://doc.owncloud.com/server/10.12/admin_manual/configuration/server/virus-scanner-support.html
	if data, infected := resp.Header["X-Infection-Found"]; infected {
		sr.Infected = infected
		re := regexp.MustCompile(`Threat=(.*);`)
		match := re.FindStringSubmatch(fmt.Sprint(data))

		if len(match) > 1 {
			sr.Description = match[1]
		}
	}

	return sr, nil
}
