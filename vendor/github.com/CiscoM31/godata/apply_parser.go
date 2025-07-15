package godata

import "context"

func ParseApplyString(ctx context.Context, apply string) (*GoDataApplyQuery, error) {
	result := GoDataApplyQuery(apply)
	return &result, nil
}
