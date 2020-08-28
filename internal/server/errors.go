package server

import "errors"

var (
	errNoAuthenticated          = errors.New("no authenticated")
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errGameAccessDenied         = errors.New("can not access game")
)
