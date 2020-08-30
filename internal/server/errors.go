package server

import "errors"

var (
	errNoAuthenticated          = errors.New("no authenticated")
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errGameAccessDenied         = errors.New("can not access game")
	errItemNotFound             = errors.New("item not found")
	errNoKeys                   = errors.New("no keys available")
	errInvalidCardNumber        = errors.New("invalid card number")
	errPerformedSession         = errors.New("session is not exist or already performed")
)
