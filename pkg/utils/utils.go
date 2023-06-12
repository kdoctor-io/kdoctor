// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
package utils

import (
	"fmt"
	"strings"
)

func GetObjNameNamespace(ns string) (name, namespace string, err error) {

	s := strings.Split(ns, "/")

	if len(s) != 2 {
		err = fmt.Errorf("failed get name and namespace from %s ", ns)
		return
	}
	name = s[1]
	namespace = s[0]
	return
}
