package stripe_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/ymgyt/stripe"
)

var (
	apiKey string
)

const (
	tokenAmex        = "tok_amex"
	tokenInvalid     = "tok_alsdkfja"
	tokenCardExpired = "tok_chargeDeclinedExpiredCard"
)

func init() {
	flag.StringVar(&apiKey, "key", os.Getenv("STRIPE_API_KEY"), "Your TEST secret key for the Stripe API.")
}

func TestClient_Customer(t *testing.T) {
	if apiKey == "" {
		t.Skip("No API key provided")
	}

	type checkFn func(*testing.T, *stripe.Customer, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoErr := func() checkFn {
		return func(t *testing.T, customer *stripe.Customer, err error) {
			if err != nil {
				t.Fatalf("err = %v; want nil", err)
			}
		}
	}
	hasErrType := func(typ string) checkFn {
		return func(t *testing.T, customer *stripe.Customer, err error) {
			se, ok := err.(stripe.Error)
			if !ok {
				t.Fatalf("err isn't a stripe.Error")
			}
			if se.Type != typ {
				t.Errorf("err.Type = %s; want %s", se.Type, typ)
			}
		}
	}
	hasIDPrefix := func() checkFn {
		return func(t *testing.T, customer *stripe.Customer, err error) {
			if !strings.HasPrefix(customer.ID, "cus_") {
				t.Errorf("ID = %s; want prefix %q", customer.ID, "cus_")
			}
		}
	}
	hasCardDefaultSource := func() checkFn {
		return func(t *testing.T, customer *stripe.Customer, err error) {
			if !strings.HasPrefix(customer.DefaultSource, "card_") {
				t.Errorf("DefaultSource = %s; want prefix %q", customer.DefaultSource, "card_")
			}
		}
	}
	hasEmail := func(email string) checkFn {
		return func(t *testing.T, customer *stripe.Customer, err error) {
			if customer.Email != email {
				t.Errorf("Email = %s; want %s", customer.Email, email)
			}
		}

	}

	c := stripe.Client{
		Key: apiKey,
	}

	tests := map[string]struct {
		token  string
		email  string
		checks []checkFn
	}{
		"valid customer with amex": {
			token:  tokenAmex,
			email:  "test@testwithgo.com",
			checks: check(hasNoErr(), hasIDPrefix(), hasCardDefaultSource(), hasEmail("test@testwithgo.com")),
		},
		"invalid token": {
			token:  tokenInvalid,
			email:  "test@testwithgo.com",
			checks: check(hasErrType(stripe.ErrTypeInvalidRequest)),
		},
		"expired card": {
			token:  tokenCardExpired,
			email:  "test@testwithgo.com",
			checks: check(hasErrType(stripe.ErrTypeCardError)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cus, err := c.Customer(tc.token, tc.email)
			for _, check := range tc.checks {
				check(t, cus, err)
			}
		})
	}
}

func TestClient_Charge(t *testing.T) {
	if apiKey == "" {
		t.Skip("No API key provided")
	}

	type checkFn func(*testing.T, *stripe.Charge, error)
	check := func(fns ...checkFn) []checkFn { return fns }

	hasNoErr := func() checkFn {
		return func(t *testing.T, charge *stripe.Charge, err error) {
			if err != nil {
				t.Fatalf("err = %v; want nil", err)
			}
		}
	}
	hasAmount := func(amount int) checkFn {
		return func(t *testing.T, charge *stripe.Charge, err error) {
			if charge.Amount != amount {
				t.Errorf("Amount = %d; want %d", charge.Amount, amount)
			}
		}
	}
	hasErrType := func(typ string) checkFn {
		return func(t *testing.T, charge *stripe.Charge, err error) {
			se, ok := err.(stripe.Error)
			if !ok {
				t.Fatalf("err isn't a stripe.Error")
			}
			if se.Type != typ {
				t.Errorf("err.Type = %s; want %s", se.Type, typ)
			}
		}
	}

	c := stripe.Client{
		Key: apiKey,
	}
	// create a customer for test
	tok := "tok_amex"
	email := "test@testwithgo.com"
	cus, err := c.Customer(tok, email)
	if err != nil {
		t.Fatalf("Customer() err = %v; want nil", err)
	}

	tests := map[string]struct {
		customerID string
		amount     int
		checks     []checkFn
	}{
		"valid charge": {
			customerID: cus.ID,
			amount:     1234,
			checks:     check(hasNoErr(), hasAmount(1234)),
		},
		"invalid customer id": {
			customerID: "cus_missing",
			amount:     1234,
			checks:     check(hasErrType(stripe.ErrTypeInvalidRequest)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			charge, err := c.Charge(tc.customerID, tc.amount)
			for _, check := range tc.checks {
				check(t, charge, err)
			}
		})
	}
}
