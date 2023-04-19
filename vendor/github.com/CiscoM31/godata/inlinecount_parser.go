package godata

import "context"

const (
	ALLPAGES = "allpages"
	NONE     = "none"
)

func ParseInlineCountString(ctx context.Context, inlinecount string) (*GoDataInlineCountQuery, error) {
	result := GoDataInlineCountQuery(inlinecount)
	if inlinecount == ALLPAGES {
		return &result, nil
	} else if inlinecount == NONE {
		return &result, nil
	} else {
		return nil, BadRequestError("Invalid inlinecount query.")
	}
}
