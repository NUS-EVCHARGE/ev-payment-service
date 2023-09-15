package userpayment

import (
	"ev-payment-service/dao"
	"ev-payment-service/dto"
)

type UserPaymentController interface {
	GetUserPaymentInfo(bookingId uint) (dto.UserPayment, error)
	CreateUserPayment(userPayment dto.UserPayment) error
	DeleteUserPayment(id uint) error
	UpdateUserPayment(userPayment dto.UserPayment) error
}

type UserControllerImpl struct {
}

func (u *UserControllerImpl) GetUserPaymentInfo(bookingId uint) (dto.UserPayment, error) {
	return dao.Db.GetUserPaymentEntry(bookingId)
}

func (u *UserControllerImpl) CreateUserPayment(userPayment dto.UserPayment) error {
	return dao.Db.CreateUserPaymentEntry(userPayment)
}

func (u *UserControllerImpl) DeleteUserPayment(id uint) error {
	return dao.Db.DeleteUserPaymentEntry(id)
}

func (u *UserControllerImpl) UpdateUserPayment(userPayment dto.UserPayment) error {
	return dao.Db.UpdateUserPaymentEntry(userPayment)
}

var (
	UserPaymentControllerObj UserPaymentController
)
