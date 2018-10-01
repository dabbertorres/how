package how

import (
	"errors"
	"fmt"
)

type errMissingValue string

func (e errMissingValue) Error() string {
	return fmt.Sprintf("missing value for '%s'", string(e))
}

// IsMissingValueError returns true if err signals a missing value
func IsMissingValueError(err error) bool {
	_, ok := err.(errMissingValue)
	return ok
}

type errNotFlag string

func (e errNotFlag) Error() string {
	return fmt.Sprintf("'%s' is not a flag", string(e))
}

// IsNotFlagError returns true if err signals a flag that isn't found in a config
func IsNotFlagError(err error) bool {
	_, ok := err.(errNotFlag)
	return ok
}

type errUnsupportedType string

func (e errUnsupportedType) Error() string {
	return fmt.Sprintf("type '%s' is not supported", string(e))
}

// IsUnsupportedTypeError returns true if err signals that a struct field is an unsupported type
func IsUnsupportedTypeError(err error) bool {
	_, ok := err.(errUnsupportedType)
	return ok
}

var (
	// ErrNotStruct signals that a config passed to a Parse function is not a struct
	ErrNotStruct = errors.New("not a struct")

	// ErrInvalidValue signals that a config passed to a Parse function is nil or not a pointer/interface
	ErrInvalidValue = errors.New("invalid value")
)
