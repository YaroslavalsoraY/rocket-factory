package model

import "errors"

var (
	ErrOrderNotFound = errors.New("Order was not found")
	ErrAlreadyPaid   = errors.New("Order has been paid")
	ErrInternal      = errors.New("Internal error")
	ErrPartsNotFound = errors.New("Not enough parts for this order")
)
