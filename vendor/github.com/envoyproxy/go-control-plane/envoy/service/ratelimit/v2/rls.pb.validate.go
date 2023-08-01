// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/service/ratelimit/v2/rls.proto

package ratelimitv2

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on RateLimitRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *RateLimitRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimitRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RateLimitRequestMultiError, or nil if none found.
func (m *RateLimitRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimitRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Domain

	for idx, item := range m.GetDescriptors() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, RateLimitRequestValidationError{
						field:  fmt.Sprintf("Descriptors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, RateLimitRequestValidationError{
						field:  fmt.Sprintf("Descriptors[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RateLimitRequestValidationError{
					field:  fmt.Sprintf("Descriptors[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for HitsAddend

	if len(errors) > 0 {
		return RateLimitRequestMultiError(errors)
	}

	return nil
}

// RateLimitRequestMultiError is an error wrapping multiple validation errors
// returned by RateLimitRequest.ValidateAll() if the designated constraints
// aren't met.
type RateLimitRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimitRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimitRequestMultiError) AllErrors() []error { return m }

// RateLimitRequestValidationError is the validation error returned by
// RateLimitRequest.Validate if the designated constraints aren't met.
type RateLimitRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimitRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimitRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimitRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimitRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimitRequestValidationError) ErrorName() string { return "RateLimitRequestValidationError" }

// Error satisfies the builtin error interface
func (e RateLimitRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimitRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimitRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimitRequestValidationError{}

// Validate checks the field values on RateLimitResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *RateLimitResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimitResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RateLimitResponseMultiError, or nil if none found.
func (m *RateLimitResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimitResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for OverallCode

	for idx, item := range m.GetStatuses() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, RateLimitResponseValidationError{
						field:  fmt.Sprintf("Statuses[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, RateLimitResponseValidationError{
						field:  fmt.Sprintf("Statuses[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RateLimitResponseValidationError{
					field:  fmt.Sprintf("Statuses[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetHeaders() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, RateLimitResponseValidationError{
						field:  fmt.Sprintf("Headers[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, RateLimitResponseValidationError{
						field:  fmt.Sprintf("Headers[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RateLimitResponseValidationError{
					field:  fmt.Sprintf("Headers[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	for idx, item := range m.GetRequestHeadersToAdd() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, RateLimitResponseValidationError{
						field:  fmt.Sprintf("RequestHeadersToAdd[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, RateLimitResponseValidationError{
						field:  fmt.Sprintf("RequestHeadersToAdd[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RateLimitResponseValidationError{
					field:  fmt.Sprintf("RequestHeadersToAdd[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return RateLimitResponseMultiError(errors)
	}

	return nil
}

// RateLimitResponseMultiError is an error wrapping multiple validation errors
// returned by RateLimitResponse.ValidateAll() if the designated constraints
// aren't met.
type RateLimitResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimitResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimitResponseMultiError) AllErrors() []error { return m }

// RateLimitResponseValidationError is the validation error returned by
// RateLimitResponse.Validate if the designated constraints aren't met.
type RateLimitResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimitResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimitResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimitResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimitResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimitResponseValidationError) ErrorName() string {
	return "RateLimitResponseValidationError"
}

// Error satisfies the builtin error interface
func (e RateLimitResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimitResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimitResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimitResponseValidationError{}

// Validate checks the field values on RateLimitResponse_RateLimit with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RateLimitResponse_RateLimit) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimitResponse_RateLimit with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RateLimitResponse_RateLimitMultiError, or nil if none found.
func (m *RateLimitResponse_RateLimit) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimitResponse_RateLimit) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

	// no validation rules for RequestsPerUnit

	// no validation rules for Unit

	if len(errors) > 0 {
		return RateLimitResponse_RateLimitMultiError(errors)
	}

	return nil
}

// RateLimitResponse_RateLimitMultiError is an error wrapping multiple
// validation errors returned by RateLimitResponse_RateLimit.ValidateAll() if
// the designated constraints aren't met.
type RateLimitResponse_RateLimitMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimitResponse_RateLimitMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimitResponse_RateLimitMultiError) AllErrors() []error { return m }

// RateLimitResponse_RateLimitValidationError is the validation error returned
// by RateLimitResponse_RateLimit.Validate if the designated constraints
// aren't met.
type RateLimitResponse_RateLimitValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimitResponse_RateLimitValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimitResponse_RateLimitValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimitResponse_RateLimitValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimitResponse_RateLimitValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimitResponse_RateLimitValidationError) ErrorName() string {
	return "RateLimitResponse_RateLimitValidationError"
}

// Error satisfies the builtin error interface
func (e RateLimitResponse_RateLimitValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimitResponse_RateLimit.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimitResponse_RateLimitValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimitResponse_RateLimitValidationError{}

// Validate checks the field values on RateLimitResponse_DescriptorStatus with
// the rules defined in the proto definition for this message. If any rules
// are violated, the first error encountered is returned, or nil if there are
// no violations.
func (m *RateLimitResponse_DescriptorStatus) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RateLimitResponse_DescriptorStatus
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// RateLimitResponse_DescriptorStatusMultiError, or nil if none found.
func (m *RateLimitResponse_DescriptorStatus) ValidateAll() error {
	return m.validate(true)
}

func (m *RateLimitResponse_DescriptorStatus) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Code

	if all {
		switch v := interface{}(m.GetCurrentLimit()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RateLimitResponse_DescriptorStatusValidationError{
					field:  "CurrentLimit",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RateLimitResponse_DescriptorStatusValidationError{
					field:  "CurrentLimit",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCurrentLimit()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RateLimitResponse_DescriptorStatusValidationError{
				field:  "CurrentLimit",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for LimitRemaining

	if len(errors) > 0 {
		return RateLimitResponse_DescriptorStatusMultiError(errors)
	}

	return nil
}

// RateLimitResponse_DescriptorStatusMultiError is an error wrapping multiple
// validation errors returned by
// RateLimitResponse_DescriptorStatus.ValidateAll() if the designated
// constraints aren't met.
type RateLimitResponse_DescriptorStatusMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RateLimitResponse_DescriptorStatusMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RateLimitResponse_DescriptorStatusMultiError) AllErrors() []error { return m }

// RateLimitResponse_DescriptorStatusValidationError is the validation error
// returned by RateLimitResponse_DescriptorStatus.Validate if the designated
// constraints aren't met.
type RateLimitResponse_DescriptorStatusValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RateLimitResponse_DescriptorStatusValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RateLimitResponse_DescriptorStatusValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RateLimitResponse_DescriptorStatusValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RateLimitResponse_DescriptorStatusValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RateLimitResponse_DescriptorStatusValidationError) ErrorName() string {
	return "RateLimitResponse_DescriptorStatusValidationError"
}

// Error satisfies the builtin error interface
func (e RateLimitResponse_DescriptorStatusValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRateLimitResponse_DescriptorStatus.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RateLimitResponse_DescriptorStatusValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RateLimitResponse_DescriptorStatusValidationError{}
