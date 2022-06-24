// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package module

import (
	"fmt"
	"strings"
)

// ResourceName generates the name for a given resource
func ResourceName(name, kind string) string {
	return strings.ToLower(fmt.Sprintf("%s+%s", name, kind))
}
