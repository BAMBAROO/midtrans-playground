package handlers

import (
	"midtrans-tester/config"
	"midtrans-tester/models"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

// ─── BCA Virtual Account ──────────────────────────────────────────────────────

// POST /core/bank-transfer/bca
func BankTransferBCA(cfg *config.Config) gin.HandlerFunc {
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
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.BankBca,
			FreeText: &coreapi.BCABankTransferDetailFreeText{
				Inquiry: []coreapi.BCABankTransferLangDetail{{LangID: "Pembayaran", LangEN: "Payment"}},
				Payment: []coreapi.BCABankTransferLangDetail{{LangID: "Terima kasih", LangEN: "Thank you"}},
			},
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "BCA VA charge failed", err)
			return
		}
		successResp(c, "bank_transfer/bca", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── BNI Virtual Account ──────────────────────────────────────────────────────

// POST /core/bank-transfer/bni
func BankTransferBNI(cfg *config.Config) gin.HandlerFunc {
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
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.BankBni,
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "BNI VA charge failed", err)
			return
		}
		successResp(c, "bank_transfer/bni", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── BRI Virtual Account ──────────────────────────────────────────────────────

// POST /core/bank-transfer/bri
func BankTransferBRI(cfg *config.Config) gin.HandlerFunc {
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
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.BankBri,
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "BRI VA charge failed", err)
			return
		}
		successResp(c, "bank_transfer/bri", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── Mandiri Bill ─────────────────────────────────────────────────────────────

// POST /core/bank-transfer/mandiri
func BankTransferMandiri(cfg *config.Config) gin.HandlerFunc {
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
		// Mandiri uses echannel (Bill Payment), not bank_transfer
		chargeReq.PaymentType = coreapi.PaymentTypeEChannel
		chargeReq.EChannel = &coreapi.EChannelDetail{
			BillInfo1: "Pembayaran:",
			BillInfo2: req.ProductName,
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "Mandiri Bill charge failed", err)
			return
		}
		successResp(c, "echannel/mandiri", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── Permata Virtual Account ──────────────────────────────────────────────────

// POST /core/bank-transfer/permata
func BankTransferPermata(cfg *config.Config) gin.HandlerFunc {
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
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.BankPermata,
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "Permata VA charge failed", err)
			return
		}
		successResp(c, "bank_transfer/permata", resp.OrderID, resp.TransactionID, resp)
	}
}

// ─── CIMB (OCBC) Virtual Account ─────────────────────────────────────────────

// POST /core/bank-transfer/cimb
func BankTransferCIMB(cfg *config.Config) gin.HandlerFunc {
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
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.BankCimb,
		}

		resp, err := client.ChargeTransaction(chargeReq)
		if err != nil {
			errResp(c, "CIMB VA charge failed", err)
			return
		}
		successResp(c, "bank_transfer/cimb", resp.OrderID, resp.TransactionID, resp)
	}
}

// setDefaults fills zero-value fields with test defaults.
func setDefaults(req *models.PaymentRequest) {
	if req.Amount <= 0 {
		req.Amount = 10000
	}
	if req.FirstName == "" {
		req.FirstName = "Test"
	}
	if req.LastName == "" {
		req.LastName = "User"
	}
	if req.Email == "" {
		req.Email = "test@example.com"
	}
	if req.Phone == "" {
		req.Phone = "08123456789"
	}
	if req.ProductName == "" {
		req.ProductName = "Test Product"
	}
}
