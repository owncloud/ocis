package svc

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	merrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/metadata"

	searchmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/search/v0"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/constants"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/net"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/prop"
	"github.com/owncloud/ocis/v2/services/webdav/pkg/propfind"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/tags"
	"github.com/owncloud/reva/v2/pkg/utils"
)

const (
	elementNameSearchFiles = "search-files"
	// TODO elementNameFilterFiles = "filter-files"
)

// Search is the endpoint for retrieving search results for REPORT requests
func (g Webdav) Search(w http.ResponseWriter, r *http.Request) {
	logger := g.log.SubloggerWithRequestID(r.Context())
	rep, err := readReport(r.Body)
	if err != nil {
		renderError(w, r, errBadRequest(err.Error()))
		logger.Debug().Err(err).Msg("error reading report")
		return
	}

	if rep.SearchFiles == nil {
		renderError(w, r, errBadRequest("missing search-files tag"))
		logger.Debug().Err(err).Msg("error reading report")
		return
	}

	t := r.Header.Get(revactx.TokenHeader)
	ctx := revactx.ContextSetToken(r.Context(), t)
	ctx = metadata.Set(ctx, revactx.TokenHeader, t)

	req := &searchsvc.SearchRequest{
		Query:    rep.SearchFiles.Search.Pattern,
		PageSize: int32(rep.SearchFiles.Search.Limit),
	}

	rsp, err := g.searchClient.Search(ctx, req)
	if err != nil {
		e := merrors.Parse(err.Error())
		switch e.Code {
		case http.StatusBadRequest:
			renderError(w, r, errBadRequest(e.Detail))
		default:
			renderError(w, r, errInternalError(err.Error()))
		}
		logger.Error().Err(err).Msg("could not get search results")
		return
	}

	g.sendSearchResponse(rsp, w, r)
}

func (g Webdav) sendSearchResponse(rsp *searchsvc.SearchResponse, w http.ResponseWriter, r *http.Request) {
	logger := g.log.SubloggerWithRequestID(r.Context())

	hrefPrefix := "/dav/spaces"
	if strings.HasPrefix(r.URL.Path, "/remote.php/dav/spaces") {
		hrefPrefix = "/remote.php/dav/spaces"
	}

	responsesXML, err := multistatusResponse(r.Context(), rsp.Matches, hrefPrefix)
	if err != nil {
		logger.Error().Err(err).Msg("error formatting propfind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set(net.HeaderDav, "1, 3, extended-mkcol")
	w.Header().Set(net.HeaderContentType, "application/xml; charset=utf-8")
	if len(rsp.Matches) > 0 {
		w.Header().Set(net.HeaderContentRange, fmt.Sprintf("rows 0-%d/%d", len(rsp.Matches)-1, rsp.TotalMatches))
	}
	w.WriteHeader(http.StatusMultiStatus)
	if _, err := w.Write(responsesXML); err != nil {
		logger.Err(err).Msg("error writing response")
	}
}

// multistatusResponse converts a list of matches into a multistatus response string
func multistatusResponse(ctx context.Context, matches []*searchmsg.Match, hrefPrefix string) ([]byte, error) {
	responses := make([]*propfind.ResponseXML, 0, len(matches))
	for i := range matches {
		res, err := matchToPropResponse(ctx, matches[i], hrefPrefix)
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

// appendPhotoProps appends photo metadata properties to the property list
func appendPhotoProps(props []prop.PropertyXML, photo *searchmsg.Photo) []prop.PropertyXML {
	if photo == nil {
		return props
	}

	if photo.TakenDateTime != nil {
		props = append(props, prop.Escaped("oc:photo-taken-date-time", photo.TakenDateTime.AsTime().Format("2006-01-02T15:04:05Z07:00")))
	}
	if photo.CameraMake != nil && *photo.CameraMake != "" {
		props = append(props, prop.Escaped("oc:photo-camera-make", *photo.CameraMake))
	}
	if photo.CameraModel != nil && *photo.CameraModel != "" {
		props = append(props, prop.Escaped("oc:photo-camera-model", *photo.CameraModel))
	}
	if photo.FNumber != nil && *photo.FNumber != 0 {
		props = append(props, prop.Escaped("oc:photo-fnumber", strconv.FormatFloat(float64(*photo.FNumber), 'f', 2, 64)))
	}
	if photo.FocalLength != nil && *photo.FocalLength != 0 {
		props = append(props, prop.Escaped("oc:photo-focal-length", strconv.FormatFloat(float64(*photo.FocalLength), 'f', 2, 64)))
	}
	if photo.Iso != nil && *photo.Iso != 0 {
		props = append(props, prop.Escaped("oc:photo-iso", strconv.FormatInt(int64(*photo.Iso), 10)))
	}
	if photo.Orientation != nil && *photo.Orientation != 0 {
		props = append(props, prop.Escaped("oc:photo-orientation", strconv.FormatInt(int64(*photo.Orientation), 10)))
	}
	if photo.ExposureNumerator != nil && *photo.ExposureNumerator != 0 {
		props = append(props, prop.Escaped("oc:photo-exposure-numerator", strconv.FormatFloat(float64(*photo.ExposureNumerator), 'f', 0, 64)))
	}
	if photo.ExposureDenominator != nil && *photo.ExposureDenominator != 0 {
		props = append(props, prop.Escaped("oc:photo-exposure-denominator", strconv.FormatFloat(float64(*photo.ExposureDenominator), 'f', 0, 64)))
	}

	return props
}

// appendLocationProps appends location metadata properties to the property list
func appendLocationProps(props []prop.PropertyXML, location *searchmsg.GeoCoordinates) []prop.PropertyXML {
	if location == nil {
		return props
	}

	if location.Latitude != nil {
		props = append(props, prop.Escaped("oc:photo-location-latitude", strconv.FormatFloat(*location.Latitude, 'f', 6, 64)))
	}
	if location.Longitude != nil {
		props = append(props, prop.Escaped("oc:photo-location-longitude", strconv.FormatFloat(*location.Longitude, 'f', 6, 64)))
	}
	if location.Altitude != nil {
		props = append(props, prop.Escaped("oc:photo-location-altitude", strconv.FormatFloat(*location.Altitude, 'f', 1, 64)))
	}

	return props
}

func matchToPropResponse(ctx context.Context, match *searchmsg.Match, hrefPrefix string) (*propfind.ResponseXML, error) {
	// unfortunately search uses own versions of ResourceId and Ref. So we need to assert them here
	var (
		ref string
		err error
	)

	// to copy PROPFIND behaviour we need to deliver different ids
	// for shares it needs to be sharestorageproviderid!shareid
	// for other spaces it needs to be storageproviderid$spaceid
	switch match.Entity.Ref.ResourceId.StorageId {
	case utils.ShareStorageProviderID:
		ref, err = storagespace.FormatReference(&provider.Reference{
			ResourceId: &provider.ResourceId{
				SpaceId:  match.Entity.Ref.ResourceId.SpaceId,
				OpaqueId: match.Entity.Ref.ResourceId.OpaqueId,
			},
			Path: match.Entity.Ref.Path,
		})
	default:
		ref, err = storagespace.FormatReference(&provider.Reference{
			ResourceId: &provider.ResourceId{
				StorageId: match.Entity.Ref.ResourceId.StorageId,
				SpaceId:   match.Entity.Ref.ResourceId.SpaceId,
			},
			Path: match.Entity.Ref.Path,
		})
	}
	if err != nil {
		return nil, err
	}

	response := propfind.ResponseXML{
		Href:     net.EncodePath(path.Join(hrefPrefix, ref)),
		Propstat: []propfind.PropstatXML{},
	}

	propstatOK := propfind.PropstatXML{
		Status: "HTTP/1.1 200 OK",
		Prop:   []prop.PropertyXML{},
	}

	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:fileid", storagespace.FormatResourceID(&provider.ResourceId{
		StorageId: match.Entity.Id.StorageId,
		SpaceId:   match.Entity.Id.SpaceId,
		OpaqueId:  match.Entity.Id.OpaqueId,
	})))

	if match.Entity.ParentId != nil {
		propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:file-parent", storagespace.FormatResourceID(&provider.ResourceId{
			StorageId: match.Entity.ParentId.StorageId,
			SpaceId:   match.Entity.ParentId.SpaceId,
			OpaqueId:  match.Entity.ParentId.OpaqueId,
		})))
	}

	if match.Entity.Ref.ResourceId.StorageId == utils.ShareStorageProviderID {
		propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:shareid", match.Entity.Ref.ResourceId.OpaqueId))
		propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:shareroot", match.Entity.ShareRootName))
		propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:remote-item-id", storagespace.FormatResourceID(&provider.ResourceId{
			StorageId: match.Entity.GetRemoteItemId().GetStorageId(),
			SpaceId:   match.Entity.GetRemoteItemId().GetSpaceId(),
			OpaqueId:  match.Entity.GetRemoteItemId().GetOpaqueId(),
		})))
	}

	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:name", match.Entity.Name))
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("d:getlastmodified", match.Entity.LastModifiedTime.AsTime().Format(constants.RFC1123)))
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("d:getcontenttype", match.Entity.MimeType))
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:permissions", match.Entity.Permissions))
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:highlights", match.Entity.Highlights))

	t := tags.New(match.Entity.Tags...)
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("oc:tags", t.AsList()))

	// those seem empty - bug?
	propstatOK.Prop = append(propstatOK.Prop, prop.Escaped("d:getetag", match.Entity.Etag))

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

	// Add photo and location metadata using helper functions
	propstatOK.Prop = appendPhotoProps(propstatOK.Prop, match.Entity.Photo)
	propstatOK.Prop = appendLocationProps(propstatOK.Prop, match.Entity.Location)

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
