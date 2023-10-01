package dao

import (
	"context"
	"ev-payment-service/dto"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

var (
	providerCollectionName = "provider_payment"
	databaseName           = "ev"
)

func (db *DbImpl) GetProviderPaymentEntry(providerId uint, billingPeriod dto.ProviderBillingPeriod) ([]dto.ProviderPayment, error) {
	providerCollection := db.MongoClient.Database(databaseName).Collection(providerCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filterCursor, err := providerCollection.Find(ctx, bson.M{"ProviderId": providerId, "providerbillingperiod": billingPeriod})
	if err != nil {
		return nil, err
	}

	var providerPaymentFiltered []dto.ProviderPayment
	if err = filterCursor.All(ctx, &providerPaymentFiltered); err != nil {
		return nil, err
	} else {
		return providerPaymentFiltered, nil
	}
}

func (db *DbImpl) CreateProviderPaymentEntry(providerPayment *dto.ProviderPayment) error {

	providerCollection := db.MongoClient.Database(databaseName).Collection(providerCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	items, err := db.GetProviderPaymentEntry(providerPayment.ProviderId, providerPayment.ProviderBillingPeriod)
	if items != nil {
		return fmt.Errorf("provider payment billing period already exists")
	}

	_, err = providerCollection.InsertOne(ctx, &providerPayment)
	return err
}

func (db *DbImpl) UpdateProviderPaymentEntry(providerPayment *dto.ProviderPayment) error {
	collection := db.MongoClient.Database(databaseName).Collection(providerCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"ProviderId": providerPayment.ProviderId, "providerbillingperiod": providerPayment.ProviderBillingPeriod}
	_, err := collection.UpdateOne(ctx, filter, bson.M{"$set": providerPayment})
	return err
}

func (db *DbImpl) DeleteProviderPaymentEntry(id uint, billingPeriod dto.ProviderBillingPeriod) error {
	collection := db.MongoClient.Database(databaseName).Collection(providerCollectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"ProviderId": id, "providerbillingperiod": billingPeriod}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}
