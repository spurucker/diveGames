package domainErrors

import "fmt"

type DependencyError struct {
	Message string
}

func NewDependencyError(message string) *DependencyError {
	return &DependencyError{message}
}

func (e *DependencyError) Error() string {
	return fmt.Sprintf(e.Message)
}
