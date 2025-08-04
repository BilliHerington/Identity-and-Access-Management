package models

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error, please try again later")
	ErrIncorrectDataFormat = errors.New("incorrect data format, please check your input data")

	ErrUserDoesNotExist        = errors.New("user does not exist")
	ErrPasswordMismatch        = errors.New("password mismatch")
	ErrEmailAlreadyRegistered  = errors.New("email already registered")
	ErrInvalidVerificationCode = errors.New("verification code is invalid")
	ErrInvalidToken            = errors.New("invalid token")

	ErrUserAlreadyExists = errors.New("user already exists")

	ErrRequestLimitExceeded = errors.New("request limit exceeded")

	ErrRoleDoesNotExist = errors.New("role does not exist")
	ErrRoleAlreadyExist = errors.New("role already exist")
	ErrRolesListEmpty   = errors.New("roles list is empty")
	//ErrNoRoleDataFound  = errors.New("no role data found")
	//ErrEmailNotFound           = errors.New("email not found")
)
