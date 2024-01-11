package utils

type AppError string

const (
	AlreadyExists     AppError = "already exists"
	ParseError        AppError = "parse error"
	ErrorNotFound     AppError = "not found"
	ErrorInvalid      AppError = "invalid"
	ErrorUnauthorized AppError = "unauthorized"
)

func (e AppError) Error() string {
	return string(e)
}
