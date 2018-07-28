package services

import (
	"errors"
	"net/http"
)

func raiseForStatus(res *Response) (err error) {
	if res.StatusCode < 300 {
		return nil
	}

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	} else {
		return errors.New(body)
	}
}
