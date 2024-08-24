package database

type DBError struct {
	s string
}

func NewDBError(text string) *DBError {
	return &DBError{text}
}

func (e *DBError) Error() string {
	return e.s
}

func (e *DBError) Unwrap() error {
	return e
}
