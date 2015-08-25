package appstore

import (
	"time"
)

type Receipt struct {
	responseVersion int
	rawReceipt      string
	Status          int
	Environment     string

	ReceiptType                string
	AdamID                     int64
	AppItemID                  int64
	BundleID                   string
	ApplicationVersion         string
	DownloadID                 int64
	OriginalApplicationVersion string
	RequestDate                time.Time
	OriginalPurchaseDate       time.Time
	InApps                     ReceiptInApps

	LatestReceiptInfo ReceiptInApps
	LatestReceipt     string
}

func (r *Receipt) String() string {
	return r.rawReceipt
}

func (r *Receipt) ResponseVersion() int {
	if r.responseVersion == 0 {
		return verIOS7
	}
	return r.responseVersion
}

func (r *Receipt) GetStatus() int {
	return r.Status
}

func (r *Receipt) GetEnvironment() string {
	return r.Environment
}

func (r *Receipt) LatestReceiptString() string {
	return r.LatestReceipt
}

func (r *Receipt) IsValidReceipt() bool {
	return r.Status == 0
}

func (r *Receipt) IsAutoRenewable() bool {
	return r.InApps.IsAutoRenewable()
}

func (r *Receipt) HasError() error {
	return HandleError(r.Status)
}

func (r *Receipt) HasExpired() bool {
	return r.Status == 21006
}

func (r *Receipt) GetTransactionIDsByProduct(product string) []int64 {
	return r.InApps.TransactionIDsByProduct(product)
}

func (r *Receipt) GetByTransactionID(id int64) *ReceiptInApp {
	return r.InApps.ByTransactionID(id)
}

type ReceiptInApp struct {
	Quantity                  int64
	ProductID                 string
	TransactionID             int64
	OriginalTransactionID     int64
	IsTrialPeriod             bool
	AppItemID                 int64
	VersionExternalIdentifier int64
	WebOrderLineItemID        int64
	PurchaseDate              time.Time
	OriginalPurchaseDate      time.Time
	ExpiresDate               time.Time
	CancellationDate          time.Time
}

type ReceiptInApps []*ReceiptInApp

func (r ReceiptInApps) IsAutoRenewable() bool {
	for _, v := range r {
		if v.ExpiresDate.IsZero() {
			return false
		}
	}
	return true
}

func (r ReceiptInApps) ByTransactionID(id int64) *ReceiptInApp {
	for _, v := range r {
		if v.TransactionID == id {
			return v
		}
	}
	return nil
}

func (r ReceiptInApps) ByProduct(productID string) ReceiptInApps {
	var matched ReceiptInApps
	for _, v := range r {
		if v.ProductID != productID {
			continue
		}
		matched = append(matched, v)
	}
	return matched
}

func (r ReceiptInApps) TransactionIDsByProduct(productID string) []int64 {
	var matched []int64
	for _, v := range r {
		if v.ProductID != productID {
			continue
		}
		matched = append(matched, v.TransactionID)
	}
	return matched
}

// for auto-renewable
func (r ReceiptInApps) LastExpiresByProductID(productID string) *ReceiptInApp {
	var latest *ReceiptInApp
	for _, v := range r {
		switch {
		case v.ProductID != productID:
			continue
		case latest != nil && latest.ExpiresDate.After(v.ExpiresDate):
			continue
		}
		latest = v
	}
	return latest
}
