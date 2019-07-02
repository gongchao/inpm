package utils

import (
	"github.com/kataras/iris/core/errors"
	"io"
	"net/http"
)

func Request(url string) (body io.ReadCloser, err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}

	if response.StatusCode == http.StatusNotFound {
		err = errors.New("not found")

		return
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New("unknown error")

		return
	}

	return response.Body, nil
}