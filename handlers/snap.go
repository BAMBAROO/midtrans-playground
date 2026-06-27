package handlers

import (
	"midtrans-tester/config"
	"midtrans-tester/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// SnapCreate creates a Snap transaction token (all payment methods via Snap UI).
// POST /snap/create
func SnapCreate(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errResp(c, "Invalid request body", err)
			return
		}

		if req.Amount <= 0 {
			req.Amount = 10000
		}

		orderID := req.OrderID
		if orderID == "" {
			orderID = generateOrderID()
		}

		var s snap.Client
		s.New(cfg.MidtransServerKey, coreEnv(cfg))

		snapReq := &snap.Request{
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: int64(req.Amount),
			},
			CustomerDetail: &midtrans.CustomerDetails{
				FName: req.FirstName,
				LName: req.LastName,
				Email: req.Email,
				Phone: req.Phone,
			},
			Items: &[]midtrans.ItemDetails{
				{
					ID:    "ITEM-01",
					Name:  req.ProductName,
					Price: int64(req.Amount),
					Qty:   1,
				},
			},
			// Enabled payment methods — omit to allow all
			EnabledPayments: snap.AllSnapPaymentType,
		}

		snapResp, err := s.CreateTransaction(snapReq)
		if err != nil {
			errResp(c, "Snap create transaction failed", err)
			return
		}

		c.JSON(http.StatusOK, models.PaymentResponse{
			Success:     true,
			Message:     "Snap transaction created",
			OrderID:     orderID,
			PaymentType: "snap",
			Data: gin.H{
				"token":        snapResp.Token,
				"redirect_url": snapResp.RedirectURL,
				"client_key":   cfg.MidtransClientKey,
			},
		})
	}
}
