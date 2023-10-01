package userpayment

import (
	"ev-payment-service/config"
	"ev-payment-service/dao"
	"ev-payment-service/dto"
	helper "ev-payment-service/helper"
	"fmt"
	"github.com/sirupsen/logrus"
)

type UserPaymentController interface {
	GetUserPaymentInfo(bookingId uint) ([]dto.UserPayment, error)
	CreateUserPayment(userPayment *dto.UserPayment, token string) (string, error)
	DeleteUserPayment(id uint) error
	UpdateUserPayment(userPayment dto.UserPayment) error
	CompleteUserPayment(userPayment *dto.UserPayment) error
}

type UserControllerImpl struct {
	stripeKey string
}

func (u *UserControllerImpl) GetUserPaymentInfo(bookingId uint) ([]dto.UserPayment, error) {
	userPaymentEntries, err := dao.Db.GetUserPaymentEntry(bookingId)
	if err != nil {
		return nil, err
	}

	return userPaymentEntries, nil
}

func (u *UserControllerImpl) CreateUserPayment(userPayment *dto.UserPayment, token string) (string, error) {

	if userPayment.Coupon != "" {
		userPayment.FinalBill = userPayment.TotalBill * 0.9
	} else {
		userPayment.FinalBill = userPayment.TotalBill
	}

	stripeClientSecret, err := helper.CreateStripeSecret(userPayment.FinalBill)

	if err != nil {
		logrus.WithField("err", err).Info("error creating payment intent")
		return "", err
	}

	// Get Booking information from booking service
	booking, err := helper.Getbooking(config.GetBookingUrl, token, userPayment.BookingId)
	//booking, err := helper.Getbooking("http://localhost:8081/api/v1/booking", token, userPayment.BookingId)
	if err != nil {
		logrus.WithField("err", err).Info("error getting booking")
		return "", err
	}

	userPayment.Booking = booking

	dbErr := dao.Db.CreateUserPaymentEntry(userPayment)

	if dbErr != nil {
		logrus.WithField("err", err).Info("error creating user payment saving into mongoDB")
		return "", dbErr
	} else {
		return stripeClientSecret, nil
	}
}

func (u *UserControllerImpl) DeleteUserPayment(id uint) error {
	return dao.Db.DeleteUserPaymentEntry(id)
}

func (u *UserControllerImpl) UpdateUserPayment(userPayment dto.UserPayment) error {
	return dao.Db.UpdateUserPaymentEntry(&userPayment)
}

func (u *UserControllerImpl) CompleteUserPayment(userPayment *dto.UserPayment) error {
	items, err := u.GetUserPaymentInfo(userPayment.BookingId)
	if err != nil {
		return err
	}

	if len(items) >= 1 {
		userPayment = &items[0]
		logrus.WithField("userPayment", userPayment).Info("userPayment")
	} else {
		return fmt.Errorf("booking id has no pening payment")
	}

	if userPayment.Status != "waiting" {
		return fmt.Errorf("payment has already been completed")
	}

	userPayment.Status = "completed"
	return dao.Db.UpdateUserPaymentEntry(userPayment)
}

var (
	UserControllerObj UserPaymentController
)

func NewUserController(stripeKey string) {
	UserControllerObj = &UserControllerImpl{
		stripeKey: stripeKey,
	}
}
