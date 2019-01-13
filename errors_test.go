package stripe_test

import (
	"encoding/json"
	"testing"

	"github.com/ymgyt/stripe"
)

var errorJSON = []byte(`{
	"error": {
	  "code": "resource_missing",
	  "doc_url": "https://stripe.com/docs/error-codes/resource-missing",
	  "message": "No such customer: cus_123",
	  "param": "customer",
	  "type": "invalid_request_error"
	}
  }`)

func TestError_Unmarshal(t *testing.T) {
	var se stripe.Error
	err := json.Unmarshal(errorJSON, &se)
	if err != nil {
		t.Fatalf("Unmarshal() err = %v; want nil", err)
	}
	wantDocURL := "https://stripe.com/docs/error-codes/resource-missing"
	if se.DocURL != wantDocURL {
		t.Errorf("DocURL = %s; want %s", se.DocURL, wantDocURL)
	}
	wantType := "invalid_request_error"
	if se.Type != wantType {
		t.Errorf("Type = %s; want %s", se.Type, wantType)
	}
}

func TestError_Marshal(t *testing.T) {
	se := stripe.Error{
		Code:    "test-code",
		DocURL:  "test-DocURL",
		Message: "test-Message",
		Param:   "test-param",
		Type:    "test-type",
	}
	data, err := json.Marshal(se)
	if err != nil {
		t.Fatalf("Marshal() err = %v; want nil", err)
	}
	var got stripe.Error
	err = json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("Unmarshal() err = %v; want nil", err)
	}
	if got != se {
		t.Errorf("got = %v; want %v", got, se)
	}
}
