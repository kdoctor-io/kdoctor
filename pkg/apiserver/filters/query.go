// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package filters

import (
	"net/http"

	"github.com/kdoctor-io/kdoctor/pkg/apiserver/request"
)

func WithRequestQuery(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(request.WithRequestQuery(req.Context(), req.URL.Query()))
		handler.ServeHTTP(w, req)
	})
}
