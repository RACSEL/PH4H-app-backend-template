package errors

import "fmt"

type HttpError struct {
	StatusCode int
	Body       []map[string]interface{}
	Err        error
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Http error: %d - %s", e.StatusCode, e.Err)
}

func (e *HttpError) Unwrap() error {
	return e.Err
}
