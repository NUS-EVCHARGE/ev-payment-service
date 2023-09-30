package userpayment

import (
	"ev-payment-service/config"
	"ev-payment-service/dao"
	"ev-payment-service/dto"
	"ev-payment-service/helper"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/paymentintent"
)

type UserPaymentController interface {
	GetUserPaymentInfo(bookingId uint) ([]dto.UserPayment, error)
	CreateUserPayment(userPayment *dto.UserPayment, token string) (string, error)
	DeleteUserPayment(id uint) error
	UpdateUserPayment(userPayment dto.UserPayment) error
	CompleteUserPayment(userPayment *dto.UserPayment) error
}

type UserControllerImpl struct {
	stripe_key string
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

	stripe.Key = u.stripe_key

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(userPayment.FinalBill)),
		Currency: stripe.String(string(stripe.CurrencySGD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}
	pi, err := paymentintent.New(params)
	logrus.WithField("pi New", pi.ClientSecret).Info("payment intent")

	if err != nil {
		logrus.WithField("err", err).Info("error creating payment intent")
		return "", err
	}

	// Get Booking information from booking service
	booking, err := helper.Getbooking(config.GetBookingUrl, token, userPayment.BookingId)
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
		return pi.ClientSecret, nil
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
		stripe_key: stripeKey,
	}
}
