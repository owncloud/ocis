package godata

import "fmt"

type GoDataError struct {
	ResponseCode int
	Message      string
	Cause        error
}

func (err *GoDataError) Error() string {
	if err.Cause != nil {
		return fmt.Sprintf("%s. Cause: %s", err.Message, err.Cause.Error())
	}
	return err.Message
}

func (err *GoDataError) Unwrap() error {
	return err.Cause
}

func (err *GoDataError) SetCause(e error) *GoDataError {
	err.Cause = e
	return err
}

func BadRequestError(message string) *GoDataError {
	return &GoDataError{400, message, nil}
}

func NotFoundError(message string) *GoDataError {
	return &GoDataError{404, message, nil}
}

func MethodNotAllowedError(message string) *GoDataError {
	return &GoDataError{405, message, nil}
}

func GoneError(message string) *GoDataError {
	return &GoDataError{410, message, nil}
}

func PreconditionFailedError(message string) *GoDataError {
	return &GoDataError{412, message, nil}
}

func InternalServerError(message string) *GoDataError {
	return &GoDataError{500, message, nil}
}

func NotImplementedError(message string) *GoDataError {
	return &GoDataError{501, message, nil}
}

type UnsupportedQueryParameterError struct {
	Parameter string
}

func (err *UnsupportedQueryParameterError) Error() string {
	return fmt.Sprintf("Query parameter '%s' is not supported", err.Parameter)
}

type DuplicateQueryParameterError struct {
	Parameter string
}

func (err *DuplicateQueryParameterError) Error() string {
	return fmt.Sprintf("Query parameter '%s' cannot be specified more than once", err.Parameter)
}
