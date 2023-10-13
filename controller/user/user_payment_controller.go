package userpayment

import (
	"ev-payment-service/config"
	"ev-payment-service/dao"
	"ev-payment-service/dto"
	helper "ev-payment-service/helper"
	"fmt"
	dto2 "github.com/NUS-EVCHARGE/ev-user-service/dto"
	"github.com/sirupsen/logrus"
)

type UserPaymentController interface {
	GetUserPaymentInfo(bookingId uint) ([]dto.UserPayment, error)
	CreateUserPayment(userPayment *dto.UserPayment, token string) (string, error)
	DeleteUserPayment(id uint) error
	UpdateUserPayment(userPayment dto.UserPayment) error
	CompleteUserPayment(userPayment *dto.UserPayment) error
	GetAllUserPayments(token string, user dto2.User) (map[string][]dto.UserPayment, error)
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

	// Temporary hardcode if user has coupon will be 90% discount
	if userPayment.Coupon != "" {
		userPayment.FinalBill = userPayment.TotalBill * 0.9
	} else {
		userPayment.FinalBill = userPayment.TotalBill
	}
	userPayment.PaymentStatus = "pending"

	stripeClientSecret, err := helper.CreateStripeSecret(userPayment.FinalBill)

	if err != nil {
		logrus.WithField("err", err).Info("error creating payment intent")
		return "", err
	}

	return stripeClientSecret, nil
}

func (u *UserControllerImpl) DeleteUserPayment(id uint) error {
	return dao.Db.DeleteUserPaymentEntry(id)
}

func (u *UserControllerImpl) UpdateUserPayment(userPayment dto.UserPayment) error {
	return dao.Db.UpdateUserPaymentEntry(&userPayment)
}

func (u *UserControllerImpl) CompleteUserPayment(userPayment *dto.UserPayment) error {
	err := userPayment.SetCompleteStatus()
	if err != nil {
		return err
	}
	return dao.Db.CreateUserPaymentEntry(userPayment)
}

func (u *UserControllerImpl) GetAllUserPayments(token string, user dto2.User) (map[string][]dto.UserPayment, error) {

	var (
		userPaymentResults = make(map[string][]dto.UserPayment)
	)

	bookings, err := helper.GetBooking(config.GetBookingUrl, token, 0)

	if err != nil {
		return nil, err
	}

	if len(bookings) == 0 {
		return nil, fmt.Errorf("no booking found for user")
	}

	for _, booking := range bookings {
		// Get Charger Info specific to the booking
		charger, err := helper.GetChargerById(config.GetProviderUrl, token, booking.ChargerId)
		if err != nil {
			return nil, err
		}

		if booking.Status == "completed" {
			payment, err := u.GetUserPaymentInfo(booking.ID)
			if err != nil {
				return nil, err
			}
			if len(payment) == 1 {
				if payment[0].PaymentStatus == "completed" {
					userPaymentResults["completed"] = append(userPaymentResults["completed"], payment[0])
				} else {
					userPaymentResults["pending"] = append(userPaymentResults["pending"], payment[0])
				}
			} else {
				//Payment is not performed and not saved into documentDB yet
				newUserPayment := dto.UserPayment{}
				newUserPayment.Booking = booking
				newUserPayment.BookingId = booking.ID
				newUserPayment.UserEmail = user.Email
				newUserPayment.PaymentStatus = "pending"

				if len(charger) == 1 {
					newUserPayment.ChargerAddress = charger[0].Address
					// Get Rate for the charger
					rates, err := helper.GetRateById(config.GetProviderUrl, token, charger[0].RatesId)
					if err != nil {
						newUserPayment.NormalRate = 1.0
					} else {
						newUserPayment.NormalRate = rates[0].NormalRate
					}
				} else {
					newUserPayment.NormalRate = 1.0
				}

				// Calculate total bill
				timeDiff := booking.EndTime.Sub(booking.StartTime)
				newUserPayment.TotalBill = timeDiff.Minutes() * newUserPayment.NormalRate
				userPaymentResults["pending"] = append(userPaymentResults["pending"], newUserPayment)
			}
		}
	}

	return userPaymentResults, nil
}

var (
	UserControllerObj UserPaymentController
)

func NewUserController(stripeKey string) {
	UserControllerObj = &UserControllerImpl{
		stripeKey: stripeKey,
	}
}
