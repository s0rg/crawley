package client

import (
	"net/http"
)

// HTTPError wraps non-200 HTTP state.
type HTTPError struct {
	code int
	msg  string
}

// ErrFromResp creates new HTTPError from response.
func ErrFromResp(resp *http.Response) (err error) {
	code := resp.StatusCode

	// reduce possible statuses to only two - 400 or 500
	switch {
	case code >= http.StatusInternalServerError:
		code = http.StatusInternalServerError
	case code >= http.StatusBadRequest:
		code = http.StatusBadRequest
	}

	return HTTPError{code: code, msg: resp.Status}
}

// Error return error textual representation.
func (herr HTTPError) Error() string {
	return herr.msg
}

// Code returns HTTP status code, caused this error.
func (herr HTTPError) Code() int {
	return herr.code
}
