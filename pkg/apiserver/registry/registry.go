// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"fmt"

	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
)

type REST struct {
	*genericregistry.Store
}

func RESTInPeace(storage *REST, err error) *REST {
	if nil != err {
		err = fmt.Errorf("unable to create REST storage for a resource due to %v, will die", err)
		panic(err)
	}

	return storage
}
