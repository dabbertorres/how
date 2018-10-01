package how

import (
	"fmt"
)

type ErrMissingValue string

func (e ErrMissingValue) Error() string {
	return fmt.Sprintf("missing value for '%s'", string(e))
}

// Returns true if err is a ErrMissingValue
func IsMissingValueError(err error) bool {
	_, ok := err.(ErrMissingValue)
	return ok
}

type ErrNotFlag string

func (e ErrNotFlag) Error() string {
	return fmt.Sprintf("'%s' is not a flag", string(e))
}

// Returns true if err is a ErrNotFlag
func IsNotFlagError(err error) bool {
	_, ok := err.(ErrNotFlag)
	return ok
}

type ErrUnsupportedType string

func (e ErrUnsupportedType) Error() string {
	return fmt.Sprintf("type '%s' is not supported", string(e))
}

// Returns true if err is a ErrUnsupportedType
func IsUnsupportedTypeError(err error) bool {
	_, ok := err.(ErrUnsupportedType)
	return ok
}
