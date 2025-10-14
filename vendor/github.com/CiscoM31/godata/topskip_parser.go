package godata

import (
	"context"
	"strconv"
)

func ParseTopString(ctx context.Context, top string) (*GoDataTopQuery, error) {
	i, err := strconv.Atoi(top)
	result := GoDataTopQuery(i)
	return &result, err
}

func ParseSkipString(ctx context.Context, skip string) (*GoDataSkipQuery, error) {
	i, err := strconv.Atoi(skip)
	result := GoDataSkipQuery(i)
	return &result, err
}
