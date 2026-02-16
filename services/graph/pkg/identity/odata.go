package identity

import (
	"strings"

	"github.com/CiscoM31/godata"
)

// GetExpandValues extracts the values of the $expand query parameter and
// returns them in a []string, rejects any $expand value that consists of more
// than just a single path segment
func GetExpandValues(req *godata.GoDataQuery) ([]string, error) {
	if req == nil || req.Expand == nil {
		return []string{}, nil
	}
	expand := make([]string, 0, len(req.Expand.ExpandItems))
	for _, item := range req.Expand.ExpandItems {
		if item.Filter != nil || item.At != nil || item.Search != nil ||
			item.OrderBy != nil || item.Skip != nil || item.Top != nil ||
			item.Select != nil || item.Compute != nil || item.Expand != nil ||
			item.Levels != 0 {
			return []string{}, godata.NotImplementedError("options for $expand not supported")
		}
		if len(item.Path) > 1 {
			return []string{}, godata.NotImplementedError("multiple segments in $expand not supported")
		}
		expand = append(expand, item.Path[0].Value)
	}
	return expand, nil
}

// GetSelectValues extracts the values of the $select query parameter and
// returns them in a []string, rejects any $select value that consists of more
// than just a single path segment
func GetSelectValues(req *godata.GoDataQuery) ([]string, error) {
	if req == nil || req.Select == nil {
		return []string{}, nil
	}
	sel := make([]string, 0, len(req.Select.SelectItems))
	for _, item := range req.Select.SelectItems {
		if len(item.Segments) > 1 {
			return []string{}, godata.NotImplementedError("multiple segments in $select not supported")
		}
		sel = append(sel, item.Segments[0].Value)
	}
	return sel, nil
}

// GetSearchValues extracts the value of the $search query parameter and returns
// it as a string. Rejects any search query that is more than just a simple string
func GetSearchValues(req *godata.GoDataQuery) (string, error) {
	if req == nil || req.Search == nil {
		return "", nil
	}

	// Only allow simple search queries for now
	if len(req.Search.Tree.Children) != 0 {
		return "", godata.NotImplementedError("complex search queries are not supported")
	}

	searchValue := strings.Trim(req.Search.Tree.Token.Value, "\"")
	return searchValue, nil
}
