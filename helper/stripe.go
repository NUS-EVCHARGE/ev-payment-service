package helper

import (
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/paymentintent"
)

var (
	StripeKey string
)

func CreateStripeSecret(payableAmount float64) (string, error) {

	stripe.Key = StripeKey

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(payableAmount)),
		Currency: stripe.String(string(stripe.CurrencySGD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)

	return pi.ClientSecret, err
}
