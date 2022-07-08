package svc

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"path"
	"strconv"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/net"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/prop"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/propfind"
	merrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"
)

const (
	elementNameSearchFiles = "search-files"
	// TODO elementNameFilterFiles = "filter-files"
)

// Search is the endpoint for retrieving search results for REPORT requests
func (g Webdav) Search(w http.ResponseWriter, r *http.Request) {
	rep, err := readReport(r.Body)
	if err != nil {
		renderError(w, r, errBadRequest(err.Error()))
		g.log.Error().Err(err).Msg("error reading report")
		return
	}

	if rep.SearchFiles == nil {
		renderError(w, r, errBadRequest("missing search-files tag"))
		g.log.Error().Err(err).Msg("error reading report")
		return
	}

	t := r.Header.Get(TokenHeader)
	ctx := revactx.ContextSetToken(r.Context(), t)
	ctx = metadata.Set(ctx, revactx.TokenHeader, t)
	rsp, err := g.searchClient.Search(ctx, &searchsvc.SearchRequest{
		Query:    rep.SearchFiles.Search.Pattern,
		PageSize: int32(rep.SearchFiles.Search.Limit),
	})
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
	responsesXML, err := multistatusResponse(r.Context(), rsp.Matches)
	if err != nil {
		g.log.Error().Err(err).Msg("error formatting propfind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(net.HeaderDav, "1, 3, extended-mkcol")
	w.Header().Set(net.HeaderContentType, "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)
	if _, err := w.Write(responsesXML); err != nil {
		g.log.Err(err).Msg("error writing response")
	}
}

// multistatusResponse converts a list of matches into a multistatus response string
func multistatusResponse(ctx context.Context, matches []*searchmsg.Match) ([]byte, error) {
	responses := make([]*propfind.ResponseXML, 0, len(matches))
	for i := range matches {
		res, err := matchToPropResponse(ctx, matches[i])
		if err != nil {
			return nil, err
		}
		responses = append(responses, res)
	}

	msr := propfind.NewMultiStatusResponseXML()
	msr.Responses = responses
	msg, err := xml.Marshal(msr)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func matchToPropResponse(ctx context.Context, match *searchmsg.Match) (*propfind.ResponseXML, error) {
	response := propfind.ResponseXML{
		Href:     net.EncodePath(path.Join("/dav/spaces/", match.Entity.Ref.ResourceId.StorageId+"!"+match.Entity.Ref.ResourceId.OpaqueId, match.Entity.Ref.Path)),
		Propstat: []propfind.PropstatXML{},
	}

	propstatOK := propfind.PropstatXML{
		Status: "HTTP/1.1 200 OK",
		Prop:   []prop.PropertyXML{},
	}

	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:fileid", match.Entity.Id.StorageId+"!"+match.Entity.Id.OpaqueId))
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("d:getetag", match.Entity.Etag))
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("d:getlastmodified", match.Entity.LastModifiedTime.AsTime().Format(time.RFC3339)))
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("d:getcontenttype", match.Entity.MimeType))

	size := strconv.FormatUint(match.Entity.Size, 10)
	if match.Entity.Type == uint64(provider.ResourceType_RESOURCE_TYPE_CONTAINER) {
		propstatOK.Prop = append(propstatOK.Prop, prop.Raw("d:resourcetype", "<d:collection/>"))
		propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:size", size))
	} else {
		propstatOK.Prop = append(propstatOK.Prop,
			prop.Escaped("d:resourcetype", ""),
			prop.Escaped("d:getcontentlength", size),
		)
	}
	score := strconv.FormatFloat(float64(match.Score), 'f', -1, 64)
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:score", score))

	if len(propstatOK.Prop) > 0 {
		response.Propstat = append(response.Propstat, propstatOK)
	}

	return &response, nil
}

type report struct {
	SearchFiles *reportSearchFiles
	// FilterFiles TODO add this for tag based search
	FilterFiles *reportFilterFiles `xml:"filter-files"`
}
type reportSearchFiles struct {
	XMLName xml.Name                `xml:"search-files"`
	Lang    string                  `xml:"xml:lang,attr,omitempty"`
	Prop    Props                   `xml:"DAV: prop"`
	Search  reportSearchFilesSearch `xml:"search"`
}
type reportSearchFilesSearch struct {
	Pattern string `xml:"pattern"`
	Limit   int    `xml:"limit"`
	Offset  int    `xml:"offset"`
}

type reportFilterFiles struct {
	XMLName xml.Name               `xml:"filter-files"`
	Lang    string                 `xml:"xml:lang,attr,omitempty"`
	Prop    Props                  `xml:"DAV: prop"`
	Rules   reportFilterFilesRules `xml:"filter-rules"`
}

type reportFilterFilesRules struct {
	Favorite  bool `xml:"favorite"`
	SystemTag int  `xml:"systemtag"`
}

// Props represents properties related to a resource
// http://www.webdav.org/specs/rfc4918.html#ELEMENT_prop (for propfind)
type Props []xml.Name

// XML holds the xml representation of a propfind
// http://www.webdav.org/specs/rfc4918.html#ELEMENT_propfind
type XML struct {
	XMLName  xml.Name  `xml:"DAV: propfind"`
	Allprop  *struct{} `xml:"DAV: allprop"`
	Propname *struct{} `xml:"DAV: propname"`
	Prop     Props     `xml:"DAV: prop"`
	Include  Props     `xml:"DAV: include"`
}

func readReport(r io.Reader) (rep *report, err error) {
	decoder := xml.NewDecoder(r)
	rep = &report{}
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			// io.EOF is a successful end
			return rep, nil
		}
		if err != nil {
			return nil, err
		}

		if v, ok := t.(xml.StartElement); ok {
			if v.Name.Local == elementNameSearchFiles {
				var repSF reportSearchFiles
				err = decoder.DecodeElement(&repSF, &v)
				if err != nil {
					return nil, err
				}
				rep.SearchFiles = &repSF
				/*
					} else if v.Name.Local == elementNameFilterFiles {
						var repFF reportFilterFiles
						err = decoder.DecodeElement(&repFF, &v)
						if err != nil {
							return nil, http.StatusBadRequest, err
						}
						rep.FilterFiles = &repFF
				*/
			}
		}
	}
}
