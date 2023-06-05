// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package request

import (
	"context"
	"net/url"
)

type queryKeyType int

const queryKey queryKeyType = iota

func WithRequestQuery(parent context.Context, query url.Values) context.Context {
	if query == nil {
		return parent
	}
	return context.WithValue(parent, queryKey, query)
}

func RequestQueryFrom(ctx context.Context) url.Values {
	query, _ := ctx.Value(queryKey).(url.Values)
	return query
}

func HasRequestQuery(ctx context.Context) bool {
	return ctx.Value(queryKey) != nil
}
