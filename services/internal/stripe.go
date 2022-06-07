package internal

import (
	"io/ioutil"
	"net/http"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
)

type StripeService struct {
}

func NewStripeService(stripeApiKey string) *StripeService {
	stripe.Key = stripeApiKey

	return &StripeService{}
}

func (s *StripeService) ConstructWebhookEvent(w http.ResponseWriter, req *http.Request, endpointSecret string) (stripe.Event, error) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return stripe.Event{}, err
	}

	signatureHeader := req.Header.Get("Stripe-Signature")
	return webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
}
