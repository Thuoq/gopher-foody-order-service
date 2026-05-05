package domain

import "errors"

var (
	ErrOrderNotFound           = errors.New("order not found")
	ErrInvalidStatusTransition = errors.New("invalid order status transition")
	ErrUnauthorized            = errors.New("unauthorized access to this order")
	ErrEmptyOrder              = errors.New("cannot create an empty order")
)
