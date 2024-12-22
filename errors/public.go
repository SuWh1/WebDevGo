package errors // wrapper around standard Go errors, designed to add a "public-facing" message to an error.

// Public wraps original error with a new error that has a 'Public() string'
// method that will return a message that is acceptable to display to the public.
// This error can also be unwrapped using the traditional 'errors' package approach.
func Public(err error, msg string) error { // err - original error; msg - user-friendly public error
	return publicError{err, msg}
}

type publicError struct {
	err error
	msg string
}

// methods
func (pe publicError) Error() string { // acceptable for public original errors
	return pe.err.Error()
}

func (pe publicError) Public() string { // show use-friendly message, not internal one
	return pe.msg
}

func (pe publicError) Unwrap() error { // to retrieve original error
	return pe.err
}
