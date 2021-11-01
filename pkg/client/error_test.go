package client

import (
	"errors"
	"net/http"
	"testing"
)

func Test_HTTPError(t *testing.T) {
	const (
		code = 666
		msg  = "infernal server error"
	)

	t.Parallel()

	e := HTTPError{code: code, msg: msg}

	if e.Error() != msg {
		t.Error("unexpected message")
	}

	if e.Code() != code {
		t.Error("unexpected code")
	}
}

func Test_HTTPErrorFromResponse(t *testing.T) {
	var (
		resp http.Response
		herr HTTPError
	)

	t.Parallel()

	resp.Status = "test"
	resp.StatusCode = http.StatusBadGateway

	err := ErrFromResp(&resp)

	if !errors.As(err, &herr) {
		t.Error("1: unexpected error")
	}

	if herr.Code() != 500 {
		t.Error("1: unexpected code")
	}

	resp.StatusCode = http.StatusNotFound

	err = ErrFromResp(&resp)

	if !errors.As(err, &herr) {
		t.Error("2: unexpected error")
	}

	if herr.Code() != 400 {
		t.Error("2: unexpected code")
	}
}
