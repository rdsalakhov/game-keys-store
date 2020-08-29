package server

import "errors"

var (
	errNoAuthenticated          = errors.New("no authenticated")
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errGameAccessDenied         = errors.New("can not access game")
	errItemNotFound             = errors.New("item not found")
	errNoKeys                   = errors.New("no keys available")
)
