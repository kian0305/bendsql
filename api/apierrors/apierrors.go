package apierrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var ProvisionWarehouseTimeout = "ProvisionWarehouseTimeout"

type APIErrorResponseBody struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type APIError struct {
	RespBody   APIErrorResponseBody
	RespText   string
	StatusCode int
	Hint       string
}

func (e APIError) Error() string {
	message := e.RespBody.Message
	if message == "" {
		message = e.RespText
	}
	message = fmt.Sprintf("%d %s", e.StatusCode, message)
	if e.Hint != "" {
		message = strings.Trim(message, ".")
		message += ". " + e.Hint
	}
	return message
}

func New(hint string, status int, respBuf []byte) error {
	respBody := APIErrorResponseBody{}
	_ = json.Unmarshal(respBuf, &respBody)
	return APIError{
		RespBody:   respBody,
		RespText:   string(respBuf),
		StatusCode: status,
		Hint:       hint,
	}
}

func IsNotFound(err error) bool {
	var apiErr APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == 404
}

func IsProxyErr(err error) bool {
	var apiErr APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == 520
}

func IsAuthFailed(err error) bool {
	var apiErr APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == 401
}

func RespBody(err error) APIErrorResponseBody {
	var apiErr APIError
	if !errors.As(err, &apiErr) {
		return APIErrorResponseBody{}
	}
	return apiErr.RespBody
}
