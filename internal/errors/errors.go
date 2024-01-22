package app

import "errors"

var ErrCouldNotOpen = errors.New("could not open")
var AlreadyExists = errors.New("already exists")
var ParseError = errors.New("parse error")
var ErrorNotFound = errors.New("not found")
var ErrorInvalid = errors.New("invalid")
