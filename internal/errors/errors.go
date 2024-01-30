package app

import "errors"

var ErrCouldNotOpen = errors.New("could not open")
var ErrAlreadyExists = errors.New("already exists")
var ErrParse = errors.New("parse error")
var ErrorNotFound = errors.New("not found")
var ErrorInvalid = errors.New("invalid")
