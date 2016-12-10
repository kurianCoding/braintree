package braintree

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

func (bt *Braintree) execute(method, path string, v interface{}, payload interface{}) error {

	url := "https://" + bt.environment + ".braintreegateway.com/merchants/" + bt.merchantID + "/" + path
	buf := new(bytes.Buffer)
	if err := xml.NewEncoder(buf).Encode(payload); err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("X-ApiVersion", "4")
	req.Header.Set("Content-Type", "application/xml")
	req.SetBasicAuth(bt.publicKey, bt.privateKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {

	case http.StatusOK, http.StatusCreated:
		return xml.NewDecoder(resp.Body).Decode(v)

	case http.StatusNotFound:
		return errors.New("404: not found")

	default:
		return fmt.Errorf("unexpected error: %d: %s", resp.StatusCode, resp.Status)
	}
}