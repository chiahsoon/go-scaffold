package models

type UnauthorizedError struct {
	message string
}

func (e UnauthorizedError) Error() string {
	return e.message
}

func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{
		message: message,
	}
}

// ==================================================================================================================

type InternalServerError struct {
	message string
}

func (e InternalServerError) Error() string {
	return e.message
}

func NewInternalServerError(message string) InternalServerError {
	return InternalServerError{
		message: message,
	}
}

// ==================================================================================================================

type BadRequestError struct {
	message string
}

func (e BadRequestError) Error() string {
	return e.message
}

func NewBadRequestError(message string) BadRequestError {
	return BadRequestError {
		message: message,
	}
}
