package domain

import (
	"errors"
)

var (
	ParseMotoError = errors.New("parse moto error")
	InternalError = errors.New("internal error")
	RecordNotFound = errors.New("record not found")
	RecordAlreadyExists = errors.New("record already exists")
)
