package app

import (
	"errors"
	"fmt"
	"time"
)

var ErrCouldNotOpen = errors.New("could not open")
var ErrAlreadyExists = errors.New("already exists")
var ErrParse = errors.New("parse error")
var ErrorNotFound = errors.New("not found")
var ErrorInvalid = errors.New("invalid")

type RetryableError struct {
	// Err          error
	Count        int
	Limit        int
	SleepingTime []int
}

func NewRetryableError() *RetryableError {
	return &RetryableError{
		// Err:          err,
		Count:        0,
		Limit:        3,
		SleepingTime: []int{1, 2, 3},
	}
}

// func (e *RetryableError) Unwrap() error {
// 	return e.Err
// }

func (e *RetryableError) Wrap(err error) error {
	return fmt.Errorf("error after %d attempts: %w", e.Count, err)
}

func (e *RetryableError) Error() string {
	return "Retryable Error: "
}

func (e *RetryableError) Sleep() {
	fmt.Println("sleeping for: ", e.SleepingTime[e.Count-1])
	time.Sleep(time.Duration(e.SleepingTime[e.Count-1]) * time.Second)

}

func (e *RetryableError) Increment() {
	e.Count++
}

func (e *RetryableError) SleepAndIncrement() {
	e.Increment()
	e.Sleep()
}

func (e *RetryableError) CanRetry() bool {
	return e.Count < e.Limit
}

// RetryWrapper excecutes the function f and retries it if it returns the error defined in retrErr
func RetryWrapper(f func() error, errIsRetriable func(error) bool, retrErr RetryableError) error {
	var err error
	for retrErr.CanRetry() {
		err = f()
		if err == nil {
			return nil
		}
		if errIsRetriable(err) {
			retrErr.SleepAndIncrement()
		} else {
			return err
		}
	}
	return retrErr.Wrap(err)
}

// RetryWrapper excecutes the function f and retries it if it returns the error defined in retrErr
func RetryWrapperWithResult(f func() (interface{}, error), errIsRetriable func(error) bool, retrErr RetryableError) (interface{}, error) {
	var err error
	for retrErr.CanRetry() {
		result, err := f()
		if err == nil {
			return result, nil
		}
		if errIsRetriable(err) {
			retrErr.SleepAndIncrement()
		} else {
			return nil, err
		}
	}
	return nil, retrErr.Wrap(err)
}
