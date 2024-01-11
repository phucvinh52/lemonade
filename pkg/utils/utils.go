package utils

import (
	"bytes"
	"io"
	"net/http"
)

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	//We Read the response body on the line below.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//Convert the body to type string
	return body, nil
}

func Post(posturl string, bodydata []byte) error {
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(bodydata))
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return nil
}
