package dao

import (
	"ev-payment-service/dto"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mySqlDbImpl struct {
	MysqlClient *gorm.DB
}

func (m mySqlDbImpl) CreateUserPaymentEntry(userPayment *dto.UserPayment) error {
	result := m.MysqlClient.Create(userPayment)
	return result.Error
}

func (m mySqlDbImpl) GetUserPaymentEntry(bookingId uint) ([]dto.UserPayment, error) {
	var userPaymentList []dto.UserPayment

	result := m.MysqlClient.Where("booking_id = ?", bookingId).Find(&userPaymentList)
	return userPaymentList, result.Error
}

func (m mySqlDbImpl) GetUserPaymentByProviderId(providerId uint) ([]dto.UserPayment, error) {
	var userPaymentList []dto.UserPayment

	result := m.MysqlClient.Where("provider_id = ?", providerId).Find(&userPaymentList)
	return userPaymentList, result.Error
}

func (m mySqlDbImpl) UpdateUserPaymentEntry(userPayment *dto.UserPayment) error {
	result := m.MysqlClient.Updates(userPayment)
	return result.Error
}

func (m mySqlDbImpl) DeleteUserPaymentEntry(id uint) error {
	result := m.MysqlClient.Where(&dto.UserPayment{}, id)
	return result.Error
}

func (m mySqlDbImpl) CreateProviderPaymentEntry(providerPayment *dto.ProviderPayment) error {
	result := m.MysqlClient.Create(providerPayment)
	return result.Error
}

func (m mySqlDbImpl) GetProviderPaymentEntry(providerId uint, billingPeriod dto.ProviderBillingPeriod) ([]dto.ProviderPayment, error) {
	var paymentList []dto.ProviderPayment
	result := m.MysqlClient.Where("provider_id = ?", providerId).Where("billing_period = ? ", billingPeriod).Find(&paymentList)
	return paymentList, result.Error
}

func (m mySqlDbImpl) UpdateProviderPaymentEntry(providerPayment *dto.ProviderPayment) error {
	result := m.MysqlClient.Updates(providerPayment)
	return result.Error
}

func (m mySqlDbImpl) DeleteProviderPaymentEntry(id uint, billingPeriod dto.ProviderBillingPeriod) error {
	result := m.MysqlClient.Where("billing_period = ?", billingPeriod).Where("id = ?", id).Delete(&dto.ProviderPayment{})
	return result.Error
}

func (m mySqlDbImpl) Close() {
	sql, _ := m.MysqlClient.DB()
	sql.Close()
}

func InitSqlDb(dsn string) {
	client, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Db = NewSqlDatabase(client)
}

func NewSqlDatabase(client *gorm.DB) Database {
	return &mySqlDbImpl{
		MysqlClient: client,
	}
}
