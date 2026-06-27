package handlers

import (
	"midtrans-tester/config"
	"midtrans-tester/models"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go/coreapi"
)

// ─── GoPay ───────────────────────────────────────────────────────────────────

// POST /core/ewallet/gopay
func EWalletGoPay(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errResp(c, "Invalid body", err)
			return
		}
		setDefaults(&req)

		var client coreapi.Client
		client.New(cfg.MidtransServerKey, coreEnv(cfg))

		chargeReq := buildChargeReq(req)
		chargeReq.PaymentType = coreapi.PaymentTypeGopay
		chargeReq.Gopay = &coreapi.GopayDetails{
			EnableCallback: req.EnableCallback,
			CallbackUrl:    req.CallbackURL,
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "GoPay charge failed", err)
			return
		}
		successResp(c, "gopay", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── ShopeePay ───────────────────────────────────────────────────────────────

// POST /core/ewallet/shopeepay
func EWalletShopeePay(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errResp(c, "Invalid body", err)
			return
		}
		setDefaults(&req)

		var client coreapi.Client
		client.New(cfg.MidtransServerKey, coreEnv(cfg))

		chargeReq := buildChargeReq(req)
		chargeReq.PaymentType = coreapi.PaymentTypeShopeepay
		chargeReq.ShopeePay = &coreapi.ShopeePayDetails{
			CallbackUrl: req.CallbackURL,
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "ShopeePay charge failed", err)
			return
		}
		successResp(c, "shopeepay", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── DANA ────────────────────────────────────────────────────────────────────

// POST /core/ewallet/dana
// Note: DANA uses payment_type "qris" with acquirer "dana" in Midtrans Core API.
func EWalletDANA(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errResp(c, "Invalid body", err)
			return
		}
		setDefaults(&req)

		var client coreapi.Client
		client.New(cfg.MidtransServerKey, coreEnv(cfg))

		chargeReq := buildChargeReq(req)
		chargeReq.PaymentType = coreapi.PaymentTypeQris
		chargeReq.Qris = &coreapi.QrisDetails{
			Acquirer: "dana",
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "DANA charge failed", err)
			return
		}
		successResp(c, "qris/dana", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── OVO ─────────────────────────────────────────────────────────────────────

// POST /core/ewallet/ovo
// OVO requires the customer's OVO-registered phone number.
func EWalletOVO(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errResp(c, "Invalid body", err)
			return
		}
		setDefaults(&req)

		phone := req.OVOPhone
		if phone == "" {
			phone = req.Phone
		}

		var client coreapi.Client
		client.New(cfg.MidtransServerKey, coreEnv(cfg))

		chargeReq := buildChargeReq(req)
		// OVO uses payment_type "qris" with acquirer "ovo"
		chargeReq.PaymentType = coreapi.PaymentTypeQris
		chargeReq.Qris = &coreapi.QrisDetails{
			Acquirer: "ovo",
		}
		// Override customer phone with OVO-registered phone
		if chargeReq.CustomerDetails != nil {
			chargeReq.CustomerDetails.Phone = phone
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "OVO charge failed", err)
			return
		}
		successResp(c, "qris/ovo", resp.OrderID, resp.TransactionID, resp)
	}
}
