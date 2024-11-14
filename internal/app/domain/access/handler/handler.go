package access_handler

import (
	"context"
	"net/http"
	"strconv"

	access_domain "github.com/kotsmile/everd-backend/internal/app/domain/access"
	"github.com/kotsmile/everd-backend/internal/util"
)

type AccessHandler struct {
	*util.HTTPHandler
}

func NewAccessHandler(httpHandler *util.HTTPHandler) *AccessHandler {
	return &AccessHandler{
		HTTPHandler: httpHandler,
	}
}

func (h *AccessHandler) AuthMiddlerware(ctx context.Context, w http.ResponseWriter, r *http.Request) (context.Context, error) {
	userID_, ok := r.Header["user-id"]
	if !ok {
		return nil, util.
			NewHTTPError("unauthorized").
			WithStatus(http.StatusUnauthorized).
			WithErrorMessage("user id is not provided in header")
	}

	userIDInt, err := strconv.Atoi(userID_[0])
	if err != nil {
		return nil, util.
			NewHTTPError("unauthorized").
			WithStatus(http.StatusUnauthorized).
			WithErrorMessage("user id is not a integer")
	}

	userID, err := access_domain.NewUserID(userIDInt)
	if err != nil {
		return nil, util.
			NewHTTPError("unauthorized").
			WithStatus(http.StatusUnauthorized).
			WithErrorMessage("user id is invalid")
	}

	return context.WithValue(ctx, "userID", userID), nil
}
