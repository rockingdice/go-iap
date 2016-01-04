package appstore

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testGetReceipt(t *testing.T) *Receipt {
	assert := assert.New(t)

	result := IAPResponseIOS7{}
	err := json.Unmarshal([]byte(testReceiptString), &result)
	assert.Nil(err, "receipt data must be a valid json")
	return result.ToReceipt()
}

func TestGetLastExpiresByProductID(t *testing.T) {
	assert := assert.New(t)

	rc := testGetReceipt(t)
	assert.Equal("com.example.app", rc.BundleID)

	inApp := rc.GetLastExpiresByProductID("com.example.app.subscription_1.v2")
	expectedExpire, _ := time.Parse(time.RFC3339, "2015-12-07T23:49:14Z")
	assert.EqualValues(1000000183885918, inApp.TransactionID)
	assert.Equal("com.example.app.subscription_1.v2", inApp.ProductID)
	assert.Equal(expectedExpire.Unix(), inApp.ExpiresDate.Unix())
}

func TestLastExpiresByProductID(t *testing.T) {
	assert := assert.New(t)

	rc := testGetReceipt(t)
	assert.Equal("com.example.app", rc.BundleID)

	inApp := rc.InApps.LastExpiresByProductID("com.example.app.subscription_1.v2")
	expectedExpire, _ := time.Parse(time.RFC3339, "2015-12-07T23:24:14Z")
	assert.Equal("com.example.app.subscription_1.v2", inApp.ProductID)
	assert.EqualValues(1000000183882899, inApp.TransactionID)
	assert.Equal(expectedExpire.Unix(), inApp.ExpiresDate.Unix())
	
	inApp = rc.LatestReceiptInfo.LastExpiresByProductID("com.example.app.subscription_1.v2")
	expectedExpire, _ = time.Parse(time.RFC3339, "2015-12-07T23:49:14Z")
	assert.Equal("com.example.app.subscription_1.v2", inApp.ProductID)
	assert.EqualValues(1000000183885918, inApp.TransactionID)
	assert.Equal(expectedExpire.Unix(), inApp.ExpiresDate.Unix())
}

func TestLastExpiresByProductIDForLatest(t *testing.T) {
	assert := assert.New(t)

	rc := testGetReceipt(t)
	assert.Equal("com.example.app", rc.BundleID)

	inApp := rc.LatestReceiptInfo.LastExpiresByProductIDForLatest("com.example.app.subscription_1.v2")
	expectedExpire, _ := time.Parse(time.RFC3339, "2015-12-07T23:49:14Z")
	assert.Equal("com.example.app.subscription_1.v2", inApp.ProductID)
	assert.EqualValues(1000000183885918, inApp.TransactionID)
	assert.Equal(expectedExpire.Unix(), inApp.ExpiresDate.Unix())
}
