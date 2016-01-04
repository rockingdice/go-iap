package appstore

import (
	"time"
)

// Receipt is struct for iap receipt data
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

// ResponseVersion returns receipt style;
// iOS6 style returns `6`
// iOS7 style returns `7`
func (r *Receipt) ResponseVersion() int {
	if r.responseVersion == 0 {
		return verIOS7
	}
	return r.responseVersion
}

// GetStatus returns status code of the receipt
// see: https://developer.apple.com/library/ios/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html
func (r *Receipt) GetStatus() int {
	return r.Status
}

func (r *Receipt) GetEnvironment() string {
	return r.Environment
}

// LatestReceiptString returns raw receipt of `latest_receipt`
func (r *Receipt) LatestReceiptString() string {
	return r.LatestReceipt
}

// IsValidReceipt checks this receipt is valid receipt or not
// if this receipt is auto-renewable and iOS6 style, expired one returns false
func (r *Receipt) IsValidReceipt() bool {
	return r.Status == 0
}

// IsAutoRenewable checks this receipt is auto-renewable subscription or not
func (r *Receipt) IsAutoRenewable() bool {
	return r.InApps.IsAutoRenewable()
}

func (r *Receipt) HasError() error {
	return HandleError(r.Status)
}

// HasExpired checks this receipt is expired or not (only for iOS6 style)
func (r *Receipt) HasExpired() bool {
	return r.Status == 21006
}

// GetTransactionIDs returns all of transaction_id from `in_app`
func (r *Receipt) GetTransactionIDs() []int64 {
	return r.InApps.TransactionIDs()
}

// GetTransactionIDsByProduct returns all of transaction_id from `in_app` filtered by `product_id`
func (r *Receipt) GetTransactionIDsByProduct(product string) []int64 {
	return r.InApps.TransactionIDsByProduct(product)
}

// GetByTransactionID returns receipt data by `transaction_id`
func (r *Receipt) GetByTransactionID(id int64) *ReceiptInApp {
	return r.InApps.ByTransactionID(id)
}

// GetLastExpiresByProductID returns latest expires receipt data by `product_id`
func (r *Receipt) GetLastExpiresByProductID(productID string) *ReceiptInApp {
	inAppLatest := r.LatestReceiptInfo.LastExpiresByProductIDForLatest(productID)
	inApp := r.InApps.LastExpiresByProductID(productID)
	switch {
	case inApp == nil:
		return inAppLatest
	case inAppLatest == nil:
		return inApp
	case inApp.ExpiresDate.After(inAppLatest.ExpiresDate):
		return inApp
	}
	return inAppLatest
}

// GetLastExpiresByTransactionIDs returns latest expires receipt data from `transaction_id` list
func (r *Receipt) GetLastExpiresByTransactionIDs(ids []int64) *ReceiptInApp {
	inAppLatest := r.LatestReceiptInfo.LastExpiresByTransactionIDsForLatest(ids)
	inApp := r.InApps.LastExpiresByTransactionIDs(ids)
	switch {
	case inApp == nil:
		return inAppLatest
	case inAppLatest == nil:
		return inApp
	case inApp.ExpiresDate.After(inAppLatest.ExpiresDate):
		return inApp
	}
	return inAppLatest
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

func (r ReceiptInApps) TransactionIDs() []int64 {
	var ids []int64
	for _, v := range r {
		ids = append(ids, v.TransactionID)
	}
	return ids
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

// for LatestReceiptInfo
func (r ReceiptInApps) LastExpiresByProductIDForLatest(productID string) *ReceiptInApp {
	for i := len(r) - 1; i >= 0; i-- {
		v := r[i]
		if v.ProductID == productID {
			return v
		}
	}
	return nil
}

func (r ReceiptInApps) LastExpiresByTransactionIDs(ids []int64) *ReceiptInApp {
	idMap := make(map[int64]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	var latest *ReceiptInApp
	for _, v := range r {
		_, ok := idMap[v.TransactionID]
		switch {
		case !ok:
			continue
		case latest != nil && latest.ExpiresDate.After(v.ExpiresDate):
			continue
		}
		latest = v
	}
	return latest
}

func (r ReceiptInApps) LastExpiresByTransactionIDsForLatest(ids []int64) *ReceiptInApp {
	idMap := make(map[int64]bool)
	for _, id := range ids {
		idMap[id] = true
	}

	for i := len(r) - 1; i >= 0; i-- {
		v := r[i]
		if _, ok := idMap[v.TransactionID]; ok {
			return v
		}
	}
	return nil
}
