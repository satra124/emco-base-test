// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation

/*
	Package emcoerror standardizes the error handling
*/
package emcoerror

// ErrorReason is an enumeration of potential failure reasons
// Each emco error type must be associated with an ErrorReason

type ErrorReason int

const (
	BadRequest ErrorReason = iota
	Conflict
	NotFound
	PreconditionFailed
	RequestTimeout
	Unknown
	UnprocessableEntity
	// Add additional reason(s)
)

// Type Error implements the emcoerror
type Error struct {
	error
	Message string
	Reason  ErrorReason
	Cause   *Error
	// Add additional properties(s)

}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Message + e.Cause.Error()
	}

	return e.Message
}
