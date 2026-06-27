package handlers

import (
	"midtrans-tester/config"
	"midtrans-tester/models"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go/coreapi"
)

// ─── QRIS (Generic) ───────────────────────────────────────────────────────────

// POST /core/qris
// Uses default acquirer (airpay/shopeepay). Can be changed to "gopay", "dana", etc.
func QRIS(cfg *config.Config) gin.HandlerFunc {
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
			Acquirer: "gopay", // or "airpay shopee" | "dana" | "ovo" | "ntt_data"
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "QRIS charge failed", err)
			return
		}
		successResp(c, "qris", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── Credit Card ──────────────────────────────────────────────────────────────

// POST /core/credit-card
// Requires token_id from Midtrans.js (frontend tokenization).
func CreditCard(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errResp(c, "Invalid body", err)
			return
		}
		setDefaults(&req)

		if req.TokenID == "" {
			errResp(c, "token_id is required for credit card payment", nil)
			return
		}

		var client coreapi.Client
		client.New(cfg.MidtransServerKey, coreEnv(cfg))

		chargeReq := buildChargeReq(req)
		chargeReq.PaymentType = coreapi.PaymentTypeCreditCard
		chargeReq.CreditCard = &coreapi.CreditCardDetails{
			TokenID:        req.TokenID,
			Authentication: true, // enable 3DS
		}

		// Installment (optional)
		if req.Installment && req.InstallmentTerm > 0 {
			chargeReq.CreditCard.InstallmentTerm = int8(req.InstallmentTerm)
			chargeReq.CreditCard.Bank = req.Bank
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "Credit card charge failed", err)
			return
		}
		successResp(c, "credit_card", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── Indomaret (Convenience Store) ────────────────────────────────────────────

// POST /core/cstore/indomaret
func CStoreIndomaret(cfg *config.Config) gin.HandlerFunc {
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
		chargeReq.PaymentType = coreapi.PaymentTypeConvenienceStore
		chargeReq.ConvStore = &coreapi.ConvStoreDetails{
			Store:   "indomaret",
			Message: "Pembayaran " + req.ProductName,
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "Indomaret charge failed", err)
			return
		}
		successResp(c, "cstore/indomaret", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── Alfamart (Convenience Store) ─────────────────────────────────────────────

// POST /core/cstore/alfamart
func CStoreAlfamart(cfg *config.Config) gin.HandlerFunc {
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
		chargeReq.PaymentType = coreapi.PaymentTypeConvenienceStore
		chargeReq.ConvStore = &coreapi.ConvStoreDetails{
			Store:             "alfamart",
			Message:           "Pembayaran " + req.ProductName,
			AlfamartFreeText1: req.ProductName,
			AlfamartFreeText2: "Terima kasih",
			AlfamartFreeText3: "Simpan struk ini",
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "Alfamart charge failed", err)
			return
		}
		successResp(c, "cstore/alfamart", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── Akulaku (Pay Later) ──────────────────────────────────────────────────────

// POST /core/akulaku
func Akulaku(cfg *config.Config) gin.HandlerFunc {
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
		chargeReq.PaymentType = coreapi.PaymentTypeAkulaku

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "Akulaku charge failed", err)
			return
		}
		successResp(c, "akulaku", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── Kredivo (Pay Later) ──────────────────────────────────────────────────────

// POST /core/kredivo
func Kredivo(cfg *config.Config) gin.HandlerFunc {
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
		chargeReq.PaymentType = coreapi.CoreapiPaymentType("kredivo")

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "Kredivo charge failed", err)
			return
		}
		successResp(c, "kredivo", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── UOB EzPay ────────────────────────────────────────────────────────────────

// POST /core/uob-ezpay
func UOBEzPay(cfg *config.Config) gin.HandlerFunc {
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
		chargeReq.PaymentType = coreapi.CoreapiPaymentType("uob_ezpay")

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "UOB EzPay charge failed", err)
			return
		}
		successResp(c, "uob_ezpay", resp.OrderID, resp.TransactionID, resp)
	}
}
