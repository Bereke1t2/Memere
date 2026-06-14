package dto

// InitiatePaymentRequest is the body for POST /payments/initiate. CourseID is a
// course purchase; omit it (with a Plan) for a subscription purchase. The
// Idempotency-Key travels in the header, not the body, so a retry is detectable
// before the body is even parsed.
type InitiatePaymentRequest struct {
	CourseID   *string `json:"course_id"`
	Plan       string  `json:"plan"`
	Provider   string  `json:"provider" binding:"required"`
	CouponCode *string `json:"coupon_code"`
}

// SubscribeRequest is the body for POST /subscriptions: the chosen plan and the
// provider to pay with. It is sugar over initiate with no course id.
type SubscribeRequest struct {
	Plan       string  `json:"plan" binding:"required"`
	Provider   string  `json:"provider" binding:"required"`
	CouponCode *string `json:"coupon_code"`
}
