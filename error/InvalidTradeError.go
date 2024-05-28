package domainErrors

import "fmt"

type InvalidTradeError struct {
	Message string
}

func NewInvalidTradeError(message string) *InvalidTradeError {
	return &InvalidTradeError{message}
}

func (e *InvalidTradeError) Error() string {
	return fmt.Sprintf(e.Message)
}
