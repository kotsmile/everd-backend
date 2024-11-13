package util

import (
	"context"
	"encoding/json"
	"net/http"
)

type Handler func(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error

func (e *HTTPError) Error() string {
	return "http error"
}

type HTTPError struct {
	message string
	status  int

	err        error
	errMessage string
}

func NewHTTPError(message string) *HTTPError {
	return &HTTPError{message: message}
}

func (e *HTTPError) WithStatus(status int) *HTTPError {
	e.status = status
	return e
}

func (e *HTTPError) WithError(err error) *HTTPError {
	e.err = err
	return e
}

func (e *HTTPError) WithErrorMessage(errMessage string) *HTTPError {
	e.errMessage = errMessage
	return e
}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type HTTPHandler struct {
	logger Logger
}

func (h *HTTPHandler) ReadJSON(
	w http.ResponseWriter,
	r *http.Request,
	data any,
) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func (h *HTTPHandler) WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers ...http.Header,
) error {
	dataJSON, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(dataJSON)

	return err
}

func (h *HTTPHandler) OkJSON(w http.ResponseWriter, data any) error {
	if err := h.WriteJSON(w, http.StatusOK, JsonResponse{
		Data: data,
	}); err != nil {
		return err
	}

	return nil
}
