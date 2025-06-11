package dbutils

import (
	"errors"
	"fmt"
	"strings"
)

const (
	CHECK_ERROR       = "check"
	NOT_NULL_ERROR    = "not null"
	UNIQUE_ERROR      = "unique"
	FOREIGN_KEY_ERROR = "foreign key"
	OTHER_ERROR       = "other"
)

type ConstraintError struct {
	Constraint string
	msg        string
}

func (e ConstraintError) Error() string {
	if e.Constraint != "" {
		return fmt.Sprintf("constraint failed (%s) %s", e.Constraint, e.msg)
	}

	return fmt.Sprintf("constraint failed: %s", e.msg)
}

func ConstraintFailed(err error) bool {
	return strings.Contains(err.Error(), "constraint failed")
}

func IsConstraintError(err error) bool {
	var constraintError ConstraintError
	return errors.As(err, &constraintError)
}

func WrapConstraintError(err error) ConstraintError {
	if strings.Contains(err.Error(), "CHECK constraint failed") {
		return ConstraintError{CHECK_ERROR, err.Error()}
	}

	if strings.Contains(err.Error(), "NOT NULL constraint failed") {
		return ConstraintError{NOT_NULL_ERROR, err.Error()}
	}

	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return ConstraintError{UNIQUE_ERROR, err.Error()}
	}

	if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
		return ConstraintError{FOREIGN_KEY_ERROR, err.Error()}
	}

	return ConstraintError{OTHER_ERROR, err.Error()}
}
