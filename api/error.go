package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type ApiError struct{
	Response *http.Response
}

func (e *ApiError) Error() string {
	if e.Response != nil {

		bodyText, _ := ioutil.ReadAll(e.Response.Body)

		return fmt.Sprintf(
			"PIOT Api Call Failed, status: %s (status code: %d), body: %s",
			e.Response.Status,
			e.Response.StatusCode,
			bodyText,
		)
	}

	return "PIOT Api Call Error"
}

func isApiAuthError(err error) bool {

	// try to typecast err to ApiError
	if apiErr, ok := err.(*ApiError); ok {
		if apiErr.Response.StatusCode == 401 {
			return true
		}
	}

	return false
}

func isApiNotFoundError(err error) bool {

	// try to typecast err to ApiError
	if apiErr, ok := err.(*ApiError); ok {
		if apiErr.Response.StatusCode == 404 {
			return true
		}
	}

	return false
}
