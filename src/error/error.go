package error

import "fmt"

type GinHTTPError struct {
	error
	StatusCode int
	Reason     error
	ReasonMsg  string
}

// NewGinHTTPError creates a new Gin HTTP error with the given message and status code.
//
// msg: the error message
// statusCode: the HTTP status code
// error: the created Gin HTTP error
func NewGinHTTPError(msg string, statusCode int) error {
	return GinHTTPError{
		StatusCode: statusCode,
		ReasonMsg:  msg,
	}
}

func WrapGinHTTPError(err error, statusCode int) error {
	return GinHTTPError{
		StatusCode: statusCode,
		Reason:     err,
	}
}

// Error returns the error message string if it exists, otherwise returns the error message of the GinHTTPError.
//
// None.
// string
func (g GinHTTPError) Error() string {
	if len(g.ReasonMsg) > 0 {
		return g.ReasonMsg
	}
	return g.Reason.Error()
}

// JSON returns a map with the error message.
//
// No parameters.
// Returns a map[string]any.
func (g GinHTTPError) JSON() map[string]any {
	return map[string]any{
		"error": g.Error(),
	}
}

// Unwrap returns the wrapped error message if available, otherwise returns the wrapped error itself.
//
// No parameters.
// Returns an error.
func (g GinHTTPError) Unwrap() error {
	if len(g.ReasonMsg) > 0 {
		return fmt.Errorf(g.ReasonMsg)
	}
	return g.Reason
}
