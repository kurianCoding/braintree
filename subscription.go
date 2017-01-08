package braintree

import (
	"encoding/xml"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

const (
	// SubscriptionStatusActive means the subscription is active.
	SubscriptionStatusActive = "Active"

	// SubscriptionStatusCanceled means the subscription was canceled.
	SubscriptionStatusCanceled = "Canceled"

	// SubscriptionStatusExpired means the subscription has expired.
	SubscriptionStatusExpired = "Expired"

	// SubscriptionStatusPastDue means the subscription is past due.
	SubscriptionStatusPastDue = "Past Due"

	// SubscriptionStatusPending means the the subscription will begin in the future.
	SubscriptionStatusPending = "Pending"
)

// A Subscription on braintree
type Subscription struct {
	// AddOns
	// Balance
	// BillingPeriodEndDate
	BillingDayOfMonth int `xml:"billing-day-of-month"`
	// BillingPeriodStartDate
	CreatedAt           time.Time `xml:"created-at"`
	CurrentBillingCycle int       `xml:"current-billing-cycle"`
	// DaysPastDue         int       `xml:"days-past-due"`
	// Descriptor
	// Discounts
	// FailureCount int    `xml:"failure-count"`
	ID string `xml:"id"`
	// MerchantAccountID string `xml:"merchant-account-id"`
	// NeverExpires      bool   `xml:"never-expires"`
	// NextBillAmount
	// NextBillingDate
	// NextBillingPeriodAmount
	// NumberOfBillingCycles int `xml:"number-of-billing-cycles"`
	// PaidThroughDate
	PaymentMethodToken string          `xml:"payment-method-token"`
	PlanID             string          `xml:"plan-id"`
	Price              decimal.Decimal `xml:"price"`
	Status             string          `xml:"status"`
	// StatusHistory
	// Transactions
	// TrialDuration     int       `xml:"trial-duration"`
	// TrialDurationUnit string    `xml:"trial-duration-unit"`
	// TrialPeriod bool      `xml:"trial-period"`
	UpdatedAt time.Time `xml:"updated-at"`
}

// SubscriptionInput is used to create or update a subscription.
type SubscriptionInput struct {
	// AddOns
	BillingDayOfMonth int
	// Descriptor
	// Discounts
	// FirstBillingDate
	ID                string
	MerchantAccountID string
	NeverExpires      bool
	// Options
	PaymentMethodNonce string
	PaymentMethodToken string
	PlanID             string
	Price              decimal.Decimal
	TrialDuration      int
	TrialDurationUnit  string
	TrialPeriod        bool
}

// SubscriptionGW is a Subscription Gateway.
type SubscriptionGW struct {
	bt *Braintree
}

// Cancel a subscription on braintree.
func (sgw SubscriptionGW) Cancel(id string) (*Subscription, error) {
	subscription := &Subscription{}
	if err := sgw.bt.execute(http.MethodPut, "subscriptions/"+id+"/cancel", subscription, nil); err != nil {
		return nil, err
	}
	return subscription, nil
}

// Create a subscription on braintree.
//
// One of PaymentMethodNonce or PaymentMethodToken is required.
// PlanID is required.
func (sgw SubscriptionGW) Create(subscriptionInput SubscriptionInput) (*Subscription, error) {
	subscription := &Subscription{}
	if err := sgw.bt.execute(http.MethodPost, "subscriptions", subscription, subscriptionInput.sanitize()); err != nil {
		return nil, err
	}
	return subscription, nil
}

// Find a subscription with a given subscription id on braintree.
func (sgw SubscriptionGW) Find(id string) (*Subscription, error) {
	subscription := &Subscription{}
	if err := sgw.bt.execute(http.MethodGet, "subscriptions/"+id, subscription, nil); err != nil {
		return nil, err
	}
	return subscription, nil
}

// Update a subscription on braintree.
//
// ID is required.
func (sgw SubscriptionGW) Update(subscriptionInput SubscriptionInput) (*Subscription, error) {
	subscription := &Subscription{}
	if err := sgw.bt.execute(http.MethodPut, "subscriptions/"+subscriptionInput.ID, subscription, subscriptionInput.sanitize()); err != nil {
		return nil, err
	}
	return subscription, nil
}

type subscriptionInputSanitized struct {
	XMLName xml.Name `xml:"subscription"`
	// AddOns
	BillingDayOfMonth int `xml:"billing-day-of-month,omitempty"`
	// Descriptor
	// Discounts
	// FirstBillingDate
	ID                string `xml:"id,omitempty"`
	MerchantAccountID string `xml:"merchant-account-id,omitempty"`
	NeverExpires      bool   `xml:"never-expires,omitempty"`
	// Options
	PaymentMethodNonce string           `xml:"payment-method-nonce,omitempty"`
	PaymentMethodToken string           `xml:"payment-method-token,omitempty"`
	PlanID             string           `xml:"plan-id"`
	Price              *decimal.Decimal `xml:"price,omitempty"`
	TrialDuration      int              `xml:"trial-duration,omitempty"`
	TrialDurationUnit  string           `xml:"trial-duration-unit,omitempty"`
	TrialPeriod        bool             `xml:"trial-period,omitempty"`
}

func (si SubscriptionInput) sanitize() subscriptionInputSanitized {
	sis := subscriptionInputSanitized{
		BillingDayOfMonth:  si.BillingDayOfMonth,
		ID:                 si.ID,
		MerchantAccountID:  si.MerchantAccountID,
		PaymentMethodNonce: si.PaymentMethodNonce,
		PaymentMethodToken: si.PaymentMethodToken,
		PlanID:             si.PlanID,
		TrialDuration:      si.TrialDuration,
		TrialDurationUnit:  si.TrialDurationUnit,
		TrialPeriod:        si.TrialPeriod,
	}
	if si.Price != (decimal.Decimal{}) {
		sis.Price = &si.Price
	}
	return sis
}