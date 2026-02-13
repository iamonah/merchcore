package payment

import "fmt"

type PaymentStatus string

var (
	PaymentPending  = newPaymentStatus("pending")
	PaymentSuccess  = newPaymentStatus("success")
	PaymentFailed   = newPaymentStatus("failed")
	PaymentCanceled = newPaymentStatus("canceled")
	PaymentRefunded = newPaymentStatus("refunded")
)

var paymentStatuses = make(map[string]PaymentStatus)

func newPaymentStatus(v string) PaymentStatus {
	ps := PaymentStatus(v)
	paymentStatuses[v] = ps
	return ps
}

func ParsePaymentStatus(v string) (PaymentStatus, error) {
	status, ok := paymentStatuses[v]
	if !ok {
		return "", fmt.Errorf("invalid payment status: %v", v)
	}
	return status, nil
}
