// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

package common

// EmcoEvent is the supported life-cycle events
type EmcoEvent string

const (
	AddChildContext EmcoEvent = "AddChildContext"
	Instantiate     EmcoEvent = "Instantiate"
	Read            EmcoEvent = "Read"
	Terminate       EmcoEvent = "Terminate"
	Update          EmcoEvent = "Update"
	UpdateDelete    EmcoEvent = "UpdateDelete"
)
