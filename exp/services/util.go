package services

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func raiseForStatus(res *http.Response) error {
	if res.StatusCode < 300 {
		return nil
	}

	message, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 500 {
		return errors.New(string(message))
	} else {
		return ClientError{
			Message:    string(message),
			StatusCode: res.StatusCode,
		}
	}
}
