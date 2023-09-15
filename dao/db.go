package dao

import (
	"ev-payment-service/dto"
)

type Database interface {
	CreateUserPaymentEntry(userPayment dto.UserPayment) error
	GetUserPaymentEntry(id uint) (dto.UserPayment, error)
	UpdateUserPaymentEntry(userPayment dto.UserPayment) error
	DeleteUserPaymentEntry(id uint) error
}

var (
	Db Database
)

type dbImpl struct {
}
