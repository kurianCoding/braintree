package braintree

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// A GatewayError is returned from braintree if a payment is blocked by
// gateway settings of the braintree account.
type GatewayError struct {
	Message string
}

// A ProcessorError is returned from braintree if a payment failed.
type ProcessorError struct {
	Code    int
	Message string
}

// A ValidationError is returned from braintree in response to an invalid api call.
type ValidationError struct {
	Attribute string `xml:"attribute"`
	Code      int    `xml:"code"`
	Message   string `xml:"message"`
}

func (e *GatewayError) Error() string {
	return e.Message
}

func (e *ProcessorError) Error() string {
	return fmt.Sprintf("Code %d: %s", e.Code, e.Message)
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Code %d: %s", e.Code, e.Message)
}

func parseError(resp *http.Response) error {

	if resp.StatusCode != http.StatusUnprocessableEntity {
		return errors.New(resp.Status)
	}

	errs := struct {
		XMLName             xml.Name          `xml:"api-error-response"`
		Address             []ValidationError `xml:"errors>address>errors>error"`
		ClientToken         []ValidationError `xml:"errors>client-token>errors>error"`
		CreditCard          []ValidationError `xml:"errors>credit-card>errors>error"`
		Customer            []ValidationError `xml:"errors>customer>errors>error"`
		Subscription        []ValidationError `xml:"errors>subscription>errors>error"`
		Transaction         []ValidationError `xml:"errors>transaction>errors>error"`
		TransactionResponse *Transaction      `xml:"transaction"`
	}{}

	if err := xml.NewDecoder(resp.Body).Decode(&errs); err != nil {
		return err
	}

	// look for ValidationErrors
	valErrs := append(errs.Address, errs.ClientToken...)
	valErrs = append(valErrs, errs.CreditCard...)
	valErrs = append(valErrs, errs.Customer...)
	valErrs = append(valErrs, errs.Subscription...)
	valErrs = append(valErrs, errs.Transaction...)
	if len(valErrs) > 0 {
		return &valErrs[0]
	}

	if tx := errs.TransactionResponse; tx != nil {

		// look for GatewayErrors
		if tx.Status == TransactionStatusGatewayRejected {
			return &GatewayError{
				Message: tx.GatewayRejectionReason,
			}
		}

		// look or ProcessorErrors
		if tx.ProcessorResponseCode != "" {
			code, err := strconv.Atoi(tx.ProcessorResponseCode)
			if err != nil {
				return err
			}
			return &ProcessorError{
				Code:    code,
				Message: tx.ProcessorResponseText,
			}
		}
	}

	// return the status for errors we aren't handling yet
	return errors.New(resp.Status)
}
