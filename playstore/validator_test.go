package playstore

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type testSignature struct {
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
	Type         string `json:"type"`
}

var testJSON = testSignature{
	PrivateKeyID: "dummyKeyID",
	PrivateKey:   "-----BEGIN PRIVATE KEY-----\nMIIBOQIBAAJBANXOa7wgs5KHMEVJmVo2eoRxEgeqiYF2oABPGYrebU+cQiE7Mwdy\nxv153DHME+9L9QzAj+fR4y5Rwva/fAsGAssCAwEAAQJATQwrFMtwCtC+22kvYywY\nsJuSlMKm9MmL1TCsErgfCj2rksRK1U+/ZY709tE3XJVYlZalWCeVhHTjs5p0pnk6\nYQIhAOw0FksytfIfpdfcREbful+LhFp1um5WjcVf7kQ73JDxAiEA57nJkG9pwnUd\nBCyIcElTVIAKU0+iFpd1208OnGxyT3sCIGaEBNkGXWmEytnxQ8DvAVjOmNcaGZwh\n/M4ZYLREtupBAiAsrpFkTWdqPKTcsi2Y4Tq1N39GMzvA+XGbWTIrDWo5UwIgHhp9\nEOnHuUuPCjpLfYM2vSFiYzaj8UJCImjkMtDwzbA=\n-----END PRIVATE KEY-----\n",
	ClientEmail:  "dummyEmail",
	ClientID:     "dummyClientID",
	Type:         "service_account",
}

func TestSetTimeout(t *testing.T) {
	_timeout := time.Second * 3
	SetTimeout(_timeout)

	if timeout != _timeout {
		t.Errorf("got %#v\nwant %#v", timeout, _timeout)
	}
}

func TestVerifySubscriptionAndroidPublisherError(t *testing.T) {
	client := Client{nil}
	expected := errors.New("client is nil")
	_, actual := client.VerifySubscription("package", "subscriptionID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifyProductAndroidPublisherError(t *testing.T) {
	client := Client{nil}
	expected := errors.New("client is nil")
	_, actual := client.VerifyProduct("package", "productID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
