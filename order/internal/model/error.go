package model

import "errors"

var (
	ErrOrderNotFound = errors.New("order was not found")
	ErrAlreadyPaid   = errors.New("order has been paid")
	ErrInternal      = errors.New("internal error")
	ErrPartsNotFound = errors.New("not enough parts for this order")
)
