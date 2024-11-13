package util

import (
	"context"
	"encoding/json"
	"errors"
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
	message    string
	statusCode int

	err error
}

func NewHTTPError(message string) *HTTPError {
	return &HTTPError{message: message}
}

func (e *HTTPError) WithStatus(statusCode int) *HTTPError {
	e.statusCode = statusCode
	return e
}

func (e *HTTPError) WithError(err error) *HTTPError {
	e.err = err
	return e
}

func (e *HTTPError) WithErrorMessage(errMessage string) *HTTPError {
	e.err = errors.New(errMessage)
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

func (h *HTTPHandler) Wrapper(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(context.Background(), w, r); err != nil {
			httpError, ok := err.(*HTTPError)
			if !ok {
				httpError = NewHTTPError("internal server error").
					WithStatus(http.StatusInternalServerError).
					WithError(err)
			}

			h.logger.WithField("status", httpError.statusCode).
				WithField("message", httpError.message).
				Errorf("%s", httpError.err)
			if err := h.ErrorJSON(w, httpError.message, httpError.statusCode); err != nil {
				h.logger.Errorf("failed to write error json: %s", err)
			}
		}
	}
}

func (h *HTTPHandler) ErrorJSON(w http.ResponseWriter, message string, statusCode int) error {
	var payload JsonResponse

	payload.Error = true
	payload.Message = message

	if err := h.WriteJSON(w, statusCode, payload); err != nil {
		return err
	}

	return nil
}
