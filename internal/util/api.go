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

type Middleware func(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) (context.Context, error)

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

type ApiHelper struct {
	logger Logger
}

func NewApiHelper(logger Logger) *ApiHelper {
	return &ApiHelper{logger: logger}
}

func (h *ApiHelper) ReadJSON(
	w http.ResponseWriter,
	r *http.Request,
	data any,
) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func (h *ApiHelper) WriteJSON(
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

func (h *ApiHelper) OkJSON(w http.ResponseWriter, data any) error {
	if err := h.WriteJSON(w, http.StatusOK, JsonResponse{
		Data: data,
	}); err != nil {
		return err
	}

	return nil
}

func (h *ApiHelper) SendError(w http.ResponseWriter, r *http.Request, err error) {
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

func (h *ApiHelper) Wrapper(handler Handler, middlewares ...Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		for _, middleware := range middlewares {
			var err error
			ctx, err = middleware(ctx, w, r)
			if err != nil {
				h.SendError(w, r, err)
				return
			}
		}

		if err := handler(ctx, w, r); err != nil {
			h.SendError(w, r, err)
			return
		}
	}
}

func (h *ApiHelper) ErrorJSON(w http.ResponseWriter, message string, statusCode int) error {
	var payload JsonResponse

	payload.Error = true
	payload.Message = message

	if err := h.WriteJSON(w, statusCode, payload); err != nil {
		return err
	}

	return nil
}
