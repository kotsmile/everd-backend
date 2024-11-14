package access_domain

import (
	"errors"
	"fmt"
)

type UserID uint

var NilUserID UserID

var (
	Err       = errors.New("access")
	ErrUserID = fmt.Errorf("%w: user id", Err)
)

func NewUserID(id int) (UserID, error) {
	if id < 0 {
		return UserID(0), ErrUserID
	}

	return UserID(uint(id)), nil
}
