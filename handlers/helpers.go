package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"midtrans-tester/config"
	"midtrans-tester/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

// initCoreEnv sets up a Core API environment from config.
// func initCoreEnv(cfg *config.Config) midtrans.ServerKey {
// 	env := midtrans.Sandbox
// 	if cfg.MidtransProduction {
// 		env = midtrans.Production
// 	}
// 	_ = env
// 	return midtrans.ServerKey(cfg.MidtransServerKey)
// }

// coreEnv returns the midtrans.EnvironmentType.
func coreEnv(cfg *config.Config) midtrans.EnvironmentType {
	if cfg.MidtransProduction {
		return midtrans.Production
	}
	return midtrans.Sandbox
}

// buildChargeReq builds the common fields of a CoreAPI charge request.
func buildChargeReq(req models.PaymentRequest) *coreapi.ChargeReq {
	orderID := req.OrderID
	if orderID == "" {
		orderID = generateOrderID()
	}
	log.Printf("[Charge] OrderID: %s, Amount: %d", orderID, req.Amount)
	return &coreapi.ChargeReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(req.Amount),
		},
		CustomerDetails: &midtrans.CustomerDetails{
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
	}
}

func generateOrderID() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("ORDER-%d-%04d", time.Now().Unix(), rand.Intn(9999))
}

func successResp(c *gin.Context, paymentType string, orderID string, txID string, data interface{}) {
	c.JSON(http.StatusOK, models.PaymentResponse{
		Success:       true,
		Message:       "Charge created successfully",
		PaymentType:   paymentType,
		OrderID:       orderID,
		TransactionID: txID,
		Data:          data,
	})
}

func errResp(c *gin.Context, msg string, err error) {
	detail := msg
	if err != nil {
		detail = fmt.Sprintf("%s: %v", msg, err)
	}
	c.JSON(http.StatusBadRequest, models.PaymentResponse{
		Success: false,
		Message: detail,
	})
}
