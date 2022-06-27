package godata

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Parse a request from the HTTP server and format it into a GoDaataRequest type
// to be passed to a provider to produce a result.
func ParseRequest(ctx context.Context, path string, query url.Values) (*GoDataRequest, error) {
	r := &GoDataRequest{
		RequestKind: RequestKindUnknown,
	}

	if err := r.ParseUrlPath(path); err != nil {
		return nil, err
	}
	if err := r.ParseUrlQuery(ctx, query); err != nil {
		return nil, err
	}
	return r, nil
}

// Compare a request to a given service, and validate the semantics and update
// the request with semantics included
func (req *GoDataRequest) SemanticizeRequest(service *GoDataService) error {

	// if request kind is a resource
	for segment := req.FirstSegment; segment != nil; segment = segment.Next {
		err := SemanticizePathSegment(segment, service)
		if err != nil {
			return err
		}
	}

	switch req.LastSegment.SemanticReference.(type) {
	case *GoDataEntitySet:
		entitySet := req.LastSegment.SemanticReference.(*GoDataEntitySet)
		entityType, err := service.LookupEntityType(entitySet.EntityType)
		if err != nil {
			return err
		}
		err = SemanticizeFilterQuery(req.Query.Filter, service, entityType)
		if err != nil {
			return err
		}
		err = SemanticizeExpandQuery(req.Query.Expand, service, entityType)
		if err != nil {
			return err
		}
		err = SemanticizeSelectQuery(req.Query.Select, service, entityType)
		if err != nil {
			return err
		}
		err = SemanticizeOrderByQuery(req.Query.OrderBy, service, entityType)
		if err != nil {
			return err
		}
		// TODO: disallow invalid query params
	case *GoDataEntityType:
		entityType := req.LastSegment.SemanticReference.(*GoDataEntityType)
		if err := SemanticizeExpandQuery(req.Query.Expand, service, entityType); err != nil {
			return err
		}
		if err := SemanticizeSelectQuery(req.Query.Select, service, entityType); err != nil {
			return err
		}
	}

	if req.LastSegment.SemanticType == SemanticTypeMetadata {
		req.RequestKind = RequestKindMetadata
	} else if req.LastSegment.SemanticType == SemanticTypeRef {
		req.RequestKind = RequestKindRef
	} else if req.LastSegment.SemanticType == SemanticTypeEntitySet {
		if req.LastSegment.Identifier == nil {
			req.RequestKind = RequestKindCollection
		} else {
			req.RequestKind = RequestKindEntity
		}
	} else if req.LastSegment.SemanticType == SemanticTypeCount {
		req.RequestKind = RequestKindCount
	} else if req.FirstSegment == nil && req.LastSegment == nil {
		req.RequestKind = RequestKindService
	}

	return nil
}

func (req *GoDataRequest) ParseUrlPath(path string) error {
	parts := strings.Split(path, "/")
	req.FirstSegment = &GoDataSegment{
		RawValue:   parts[0],
		Name:       ParseName(parts[0]),
		Identifier: ParseIdentifiers(parts[0]),
	}
	currSegment := req.FirstSegment
	for _, v := range parts[1:] {
		temp := &GoDataSegment{
			RawValue:   v,
			Name:       ParseName(v),
			Identifier: ParseIdentifiers(v),
			Prev:       currSegment,
		}
		currSegment.Next = temp
		currSegment = temp
	}
	req.LastSegment = currSegment

	return nil
}

func SemanticizePathSegment(segment *GoDataSegment, service *GoDataService) error {
	var err error

	if segment.RawValue == "$metadata" {
		if segment.Next != nil || segment.Prev != nil {
			return BadRequestError("A metadata segment must be alone.")
		}

		segment.SemanticType = SemanticTypeMetadata
		segment.SemanticReference = service.Metadata
		return nil
	}

	if segment.RawValue == "$ref" {
		// this is a ref segment
		if segment.Next != nil {
			return BadRequestError("A $ref segment must be last.")
		}
		if segment.Prev == nil {
			return BadRequestError("A $ref segment must be preceded by something.")
		}

		segment.SemanticType = SemanticTypeRef
		segment.SemanticReference = segment.Prev
		return nil
	}

	if segment.RawValue == "$count" {
		// this is a ref segment
		if segment.Next != nil {
			return BadRequestError("A $count segment must be last.")
		}
		if segment.Prev == nil {
			return BadRequestError("A $count segment must be preceded by something.")
		}

		segment.SemanticType = SemanticTypeCount
		segment.SemanticReference = segment.Prev
		return nil
	}

	if _, ok := service.EntitySetLookup[segment.Name]; ok {
		// this is an entity set
		segment.SemanticType = SemanticTypeEntitySet
		segment.SemanticReference, err = service.LookupEntitySet(segment.Name)
		if err != nil {
			return err
		}

		if segment.Prev == nil {
			// this is the first segment
			if segment.Next == nil {
				// this is the only segment
				return nil
			} else {
				// there is at least one more segment
				if segment.Identifier != nil {
					return BadRequestError("An entity set must be the last segment.")
				}
				// if it has an identifier, it is allowed
				return nil
			}
		} else if segment.Next == nil {
			// this is the last segment in a sequence of more than one
			return nil
		} else {
			// this is a middle segment
			if segment.Identifier != nil {
				return BadRequestError("An entity set must be the last segment.")
			}
			// if it has an identifier, it is allowed
			return nil
		}
	}

	if segment.Prev != nil && segment.Prev.SemanticType == SemanticTypeEntitySet {
		// previous segment was an entity set
		semanticRef := segment.Prev.SemanticReference.(*GoDataEntitySet)

		entity, err := service.LookupEntityType(semanticRef.EntityType)

		if err != nil {
			return err
		}

		for _, p := range entity.Properties {
			if p.Name == segment.Name {
				segment.SemanticType = SemanticTypeProperty
				segment.SemanticReference = p
				return nil
			}
		}

		return BadRequestError("A valid entity property must follow entity set.")
	}

	return BadRequestError("Invalid segment " + segment.RawValue)
}

var supportedOdataKeywords = map[string]bool{
	"$filter":      true,
	"$apply":       true,
	"$expand":      true,
	"$select":      true,
	"$orderby":     true,
	"$top":         true,
	"$skip":        true,
	"$count":       true,
	"$inlinecount": true,
	"$search":      true,
	"$format":      true,
	"at":           true,
	"tags":         true,
}

type OdataComplianceConfig int

const (
	ComplianceStrict OdataComplianceConfig = 0
	// Ignore duplicate ODATA keywords in the URL query.
	ComplianceIgnoreDuplicateKeywords OdataComplianceConfig = 1 << iota
	// Ignore unknown ODATA keywords in the URL query.
	ComplianceIgnoreUnknownKeywords
	// Ignore extraneous comma as the last character in a list of function arguments.
	ComplianceIgnoreInvalidComma
	ComplianceIgnoreAll OdataComplianceConfig = ComplianceIgnoreDuplicateKeywords |
		ComplianceIgnoreUnknownKeywords |
		ComplianceIgnoreInvalidComma
)

type parserConfigKey int

const (
	odataCompliance parserConfigKey = iota
)

// If the lenient mode is set, the 'failOnConfig' bits are used to determine the ODATA compliance.
// This is mostly for historical reasons because the original parser had compliance issues.
// If the lenient mode is not set, the parser returns an error.
func WithOdataComplianceConfig(ctx context.Context, cfg OdataComplianceConfig) context.Context {
	return context.WithValue(ctx, odataCompliance, cfg)
}

// ParseUrlQuery parses the URL query, applying optional logic specified in the context.
func (req *GoDataRequest) ParseUrlQuery(ctx context.Context, query url.Values) error {
	cfg, hasComplianceConfig := ctx.Value(odataCompliance).(OdataComplianceConfig)
	if !hasComplianceConfig {
		// Strict ODATA compliance by default.
		cfg = ComplianceStrict
	}
	// Validate each query parameter is a valid ODATA keyword.
	for key, val := range query {
		if _, ok := supportedOdataKeywords[key]; !ok && (cfg&ComplianceIgnoreUnknownKeywords == 0) {
			return BadRequestError(fmt.Sprintf("Query parameter '%s' is not supported", key)).
				SetCause(&UnsupportedQueryParameterError{key})
		}
		if (cfg&ComplianceIgnoreDuplicateKeywords == 0) && (len(val) > 1) {
			return BadRequestError(fmt.Sprintf("Query parameter '%s' cannot be specified more than once", key)).
				SetCause(&DuplicateQueryParameterError{key})
		}
	}
	filter := query.Get("$filter")
	at := query.Get("at")
	apply := query.Get("$apply")
	expand := query.Get("$expand")
	sel := query.Get("$select")
	orderby := query.Get("$orderby")
	top := query.Get("$top")
	skip := query.Get("$skip")
	count := query.Get("$count")
	inlinecount := query.Get("$inlinecount")
	search := query.Get("$search")
	format := query.Get("$format")

	result := &GoDataQuery{}

	var err error = nil
	if filter != "" {
		result.Filter, err = ParseFilterString(ctx, filter)
	}
	if err != nil {
		return err
	}
	if at != "" {
		result.At, err = ParseFilterString(ctx, at)
	}
	if err != nil {
		return err
	}
	if at != "" {
		result.At, err = ParseFilterString(ctx, at)
	}
	if err != nil {
		return err
	}
	if apply != "" {
		result.Apply, err = ParseApplyString(ctx, apply)
	}
	if err != nil {
		return err
	}
	if expand != "" {
		result.Expand, err = ParseExpandString(ctx, expand)
	}
	if err != nil {
		return err
	}
	if sel != "" {
		result.Select, err = ParseSelectString(ctx, sel)
	}
	if err != nil {
		return err
	}
	if orderby != "" {
		result.OrderBy, err = ParseOrderByString(ctx, orderby)
	}
	if err != nil {
		return err
	}
	if top != "" {
		result.Top, err = ParseTopString(ctx, top)
	}
	if err != nil {
		return err
	}
	if skip != "" {
		result.Skip, err = ParseSkipString(ctx, skip)
	}
	if err != nil {
		return err
	}
	if count != "" {
		result.Count, err = ParseCountString(ctx, count)
	}
	if err != nil {
		return err
	}
	if inlinecount != "" {
		result.InlineCount, err = ParseInlineCountString(ctx, inlinecount)
	}
	if err != nil {
		return err
	}
	if search != "" {
		result.Search, err = ParseSearchString(ctx, search)
	}
	if err != nil {
		return err
	}
	if format != "" {
		err = NotImplementedError("Format is not supported")
	}
	if err != nil {
		return err
	}
	req.Query = result
	return err
}

func ParseIdentifiers(segment string) *GoDataIdentifier {
	if !(strings.Contains(segment, "(") && strings.Contains(segment, ")")) {
		return nil
	}

	rawIds := segment[strings.LastIndex(segment, "(")+1 : strings.LastIndex(segment, ")")]
	parts := strings.Split(rawIds, ",")

	result := make(GoDataIdentifier)

	for _, v := range parts {
		if strings.Contains(v, "=") {
			split := strings.SplitN(v, "=", 2)
			result[split[0]] = split[1]
		} else {
			result[v] = ""
		}
	}

	return &result
}

func ParseName(segment string) string {
	if strings.Contains(segment, "(") {
		return segment[:strings.LastIndex(segment, "(")]
	} else {
		return segment
	}
}
