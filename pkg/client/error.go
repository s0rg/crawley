package client

import (
	"net/http"
)

type HTTPError struct {
	code int
	msg  string
}

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

func (herr HTTPError) Error() string {
	return herr.msg
}

func (herr HTTPError) Code() int {
	return herr.code
}
