package errors

import "errors"

// these vars give acceess to existing funcs in the std errors package.
// We can wrap them in custom functionality as needed if we want, or mock them during testing.
var (
	As = errors.As
	Is = errors.Is
)
