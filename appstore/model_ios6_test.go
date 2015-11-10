package appstore

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIAPResponseIOS6(t *testing.T) {
	assert := assert.New(t)

	rc := "receipt"
	r := NewIAPResponseIOS6(rc)
	assert.Equal(rc, r.rawReceipt)
}

func TestIAPResponseIOS6(t *testing.T) {
	assert := assert.New(t)

	rc := "receipt"
	r := NewIAPResponseIOS6(rc)

	err := json.Unmarshal([]byte(_receipt), r)
	assert.Nil(err)
	assert.Equal(0, r.Status)
	assert.Equal("dummy_latest_receipt", r.LatestReceipt)
	assert.Equal("900000001", r.Receipt.AppItemID)
	assert.Equal("com.example.app", r.Receipt.BundleID)
	assert.Equal("0.1", r.Receipt.ApplicationVersion)
	assert.Equal("", r.Receipt.OriginalApplicationVersion)
	assert.Equal("90000000000001", r.Receipt.OriginalTransactionID)
	assert.Equal("com.example.product.item", r.Receipt.ProductID)
	assert.Equal("1", r.Receipt.Quantity)
	assert.Equal("90000000000001", r.Receipt.TransactionID)
	assert.Equal("900000000", r.Receipt.VersionExternalIdentifier)
	assert.Equal("70000000000001", r.Receipt.WebOrderLineItemID)
	assert.Equal("2015-06-09 23:31:25 Etc/GMT", r.Receipt.ExpiresDate)
	assert.Equal("1433892685000", r.Receipt.ExpiresDateMS)
	assert.Equal("2015-06-09 16:31:25 America/Los_Angeles", r.Receipt.ExpiresDatePST)

	assert.Equal("", r.Receipt.RequestDate.RequestDate)
	assert.Equal("2015-05-09 23:31:25 Etc/GMT", r.Receipt.PurchaseDate.PurchaseDate)
	assert.Equal("2015-05-09 23:31:28 Etc/GMT", r.Receipt.OriginalPurchaseDate.OriginalPurchaseDate)
}

func TestToIOS7(t *testing.T) {
	assert := assert.New(t)

	rc := "receipt"
	r6 := NewIAPResponseIOS6(rc)
	err := json.Unmarshal([]byte(_receipt), r6)
	assert.Nil(err)

	r7 := r6.ToIOS7()
	assert.Equal(0, r7.Status)
	assert.Equal(6, r7.responseVersion)
	assert.Equal("receipt", r7.rawReceipt)
	assert.Equal("dummy_latest_receipt", r7.LatestReceipt)
	assert.Equal("", r7.Environment)
	assert.Equal("", r7.Receipt.ReceiptType)
	assert.Equal(int64(0), r7.Receipt.AdamID)
	assert.Equal(int64(900000001), r7.Receipt.AppItemID)
	assert.Equal("com.example.app", r7.Receipt.BundleID)
	assert.Equal("0.1", r7.Receipt.ApplicationVersion)
	assert.Equal(int64(0), r7.Receipt.DownloadID)
	assert.Equal("", r7.Receipt.OriginalApplicationVersion)
	assert.Equal("", r7.Receipt.RequestDate.RequestDate)
	assert.Equal("2015-05-09 23:31:28 Etc/GMT", r7.Receipt.OriginalPurchaseDate.OriginalPurchaseDate)

	inapp := r7.Receipt.InApp[0]
	assert.Equal("1", inapp.Quantity)
	assert.Equal("com.example.product.item", inapp.ProductID)
	assert.Equal("90000000000001", inapp.TransactionID)
	assert.Equal("90000000000001", inapp.OriginalTransactionID)
	assert.Equal("", inapp.IsTrialPeriod)
	assert.Equal("", inapp.AppItemID)
	assert.Equal("900000000", inapp.VersionExternalIdentifier)
	assert.Equal("70000000000001", inapp.WebOrderLineItemID)
	assert.Equal("2015-05-09 23:31:25 Etc/GMT", inapp.PurchaseDate.PurchaseDate)
	assert.Equal("2015-05-09 23:31:28 Etc/GMT", inapp.OriginalPurchaseDate.OriginalPurchaseDate)
	assert.Equal("2015-06-09 23:31:25 Etc/GMT", inapp.ExpiresDate.ExpiresDate)
	assert.Equal("", inapp.CancellationDate.CancellationDate)

	latest := r7.LatestReceiptInfo[0]
	assert.Equal("1", latest.Quantity)
	assert.Equal("com.example.product.item", latest.ProductID)
	assert.Equal("90000000000002", latest.TransactionID)
	assert.Equal("10000010000000", latest.OriginalTransactionID)
	assert.Equal("", latest.IsTrialPeriod)
	assert.Equal("", latest.AppItemID)
	assert.Equal("", latest.VersionExternalIdentifier)
	assert.Equal("70000000000002", latest.WebOrderLineItemID)
	assert.Equal("2015-11-10 00:31:25 Etc/GMT", latest.PurchaseDate.PurchaseDate)
	assert.Equal("2015-05-09 23:31:28 Etc/GMT", latest.OriginalPurchaseDate.OriginalPurchaseDate)
	assert.Equal("2015-12-10 00:31:25 Etc/GMT", latest.ExpiresDate.ExpiresDate)
	assert.Equal("", latest.CancellationDate.CancellationDate)
}

func TestToIOS7ToReceipt(t *testing.T) {
	assert := assert.New(t)

	rc := "receipt"
	r6 := NewIAPResponseIOS6(rc)
	err := json.Unmarshal([]byte(_receipt), r6)
	assert.Nil(err)

	r7 := r6.ToIOS7()
	r := r7.ToReceipt()
	assert.Equal(6, r.responseVersion)
}

var _receipt = `{
"receipt":{"expires_date_formatted":"2015-06-09 23:31:25 Etc/GMT", "original_purchase_date_pst":"2015-05-09 16:31:28 America/Los_Angeles", "unique_identifier":"f000000000000000000000000000000000000000", "original_transaction_id":"90000000000001", "expires_date":"1433892685000", "app_item_id":"900000001", "transaction_id":"90000000000001", "quantity":"1", "expires_date_formatted_pst":"2015-06-09 16:31:25 America/Los_Angeles", "product_id":"com.example.product.item", "bvrs":"0.1", "unique_vendor_identifier":"F0000000-F000-F000-F000-F00000000000", "web_order_line_item_id":"70000000000001", "original_purchase_date_ms":"1431214288317", "version_external_identifier":"900000000", "bid":"com.example.app", "purchase_date_ms":"1431214285000", "purchase_date":"2015-05-09 23:31:25 Etc/GMT", "purchase_date_pst":"2015-05-09 16:31:25 America/Los_Angeles", "original_purchase_date":"2015-05-09 23:31:28 Etc/GMT", "item_id":"826168883"},
"latest_receipt_info":{"original_purchase_date_pst":"2015-05-09 16:31:28 America/Los_Angeles", "unique_identifier":"f000000000000000000000000000000000000000", "original_transaction_id":"10000010000000", "expires_date":"1449707485000", "app_item_id":"900000001", "transaction_id":"90000000000002", "quantity":"1", "product_id":"com.example.product.item", "bvrs":"0.1", "bid":"com.example.app", "unique_vendor_identifier":"F0000000-F000-F000-F000-F00000000000", "web_order_line_item_id":"70000000000002", "original_purchase_date_ms":"1431214288000", "expires_date_formatted":"2015-12-10 00:31:25 Etc/GMT", "purchase_date":"2015-11-10 00:31:25 Etc/GMT", "purchase_date_ms":"1447115485000", "expires_date_formatted_pst":"2015-12-09 16:31:25 America/Los_Angeles", "purchase_date_pst":"2015-11-09 16:31:25 America/Los_Angeles", "original_purchase_date":"2015-05-09 23:31:28 Etc/GMT", "item_id":"826168883"}, "status":0,
"latest_receipt":"dummy_latest_receipt"}`
