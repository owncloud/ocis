package godata

import (
	"context"
	"strconv"
)

func ParseCountString(ctx context.Context, count string) (*GoDataCountQuery, error) {
	i, err := strconv.ParseBool(count)
	if err != nil {
		return nil, err
	}

	result := GoDataCountQuery(i)

	return &result, nil
}
