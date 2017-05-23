package fitbit

import (
	"net/http"
)

type SubscriptionService struct {
	authedClient *http.Client
}

func newSubscription(authClient *http.Client) *SubscriptionService {
	var theService = new(SubscriptionService)
	theService.authedClient = authClient
	return theService
}
