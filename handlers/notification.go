package handlers

import (
	"log"
	"midtrans-tester/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go/coreapi"
)

// ─── Notification (Webhook) ───────────────────────────────────────────────────

// POST /notification
// Midtrans will POST here when a transaction status changes.
// Configure this URL in: Midtrans Dashboard > Settings > Configuration > Payment Notification URL
func Notification(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var notifPayload map[string]interface{}
		if err := c.ShouldBindJSON(&notifPayload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		var client coreapi.Client
		client.New(cfg.MidtransServerKey, coreEnv(cfg))

		// Verify notification signature
		orderID, exists := notifPayload["order_id"].(string)
		if !exists || orderID == "" {
			// Check if it is a recurring / subscription notification
			if subID, subExists := notifPayload["id"].(string); subExists && subID != "" {
				subName, _ := notifPayload["name"].(string)
				subStatus, _ := notifPayload["status"].(string)
				log.Printf("[Notification] Received recurring payment notification: sub_id=%s sub_name=%s status=%s", subID, subName, subStatus)
				c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "recurring notification received"})
				return
			}

			// Check if it is a GoPay account linking notification
			if accountID, acctExists := notifPayload["account_id"].(string); acctExists && accountID != "" {
				acctStatus, _ := notifPayload["account_status"].(string)
				paymentType, _ := notifPayload["payment_type"].(string)
				log.Printf("[Notification] Received account linking notification: account_id=%s payment_type=%s status=%s", accountID, paymentType, acctStatus)
				c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "account linking notification received"})
				return
			}

			c.JSON(http.StatusBadRequest, gin.H{"error": "missing order_id"})
			return
		}

		// Handle mock dashboard test notification
		if strings.HasPrefix(orderID, "payment_notif_test_") {
			log.Printf("[Notification] Received mock test notification from Midtrans dashboard for order %s. Responding with 200 OK.", orderID)
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "mock notification received"})
			return
		}

		// Re-fetch status from Midtrans to avoid spoofing
		transactionStatus, err := client.CheckTransaction(orderID)
		if err != nil {
			log.Printf("[Notification] Failed to verify order %s: %v", orderID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "status check failed"})
			return
		}

		log.Printf("[Notification] order_id=%s status=%s fraud=%s payment_type=%s",
			transactionStatus.OrderID,
			transactionStatus.TransactionStatus,
			transactionStatus.FraudStatus,
			transactionStatus.PaymentType,
		)

		// ── Business logic by status ─────────────────────────────────────────
		switch transactionStatus.TransactionStatus {
		case "capture":
			if transactionStatus.FraudStatus == "accept" {
				log.Printf("[Notification] ✅ PAID (capture+accept): %s", orderID)
				// TODO: mark order as paid in your DB
			} else if transactionStatus.FraudStatus == "challenge" {
				log.Printf("[Notification] ⚠️  CHALLENGE (manual review needed): %s", orderID)
				// TODO: set order to pending-review
			}

		case "settlement":
			log.Printf("[Notification] ✅ SETTLED: %s", orderID)
			// TODO: mark order as paid in your DB

		case "deny":
			log.Printf("[Notification] ❌ DENIED: %s", orderID)
			// TODO: mark order as failed

		case "cancel", "expire":
			log.Printf("[Notification] ❌ CANCELLED/EXPIRED: %s", orderID)
			// TODO: mark order as cancelled, release reserved stock

		case "pending":
			log.Printf("[Notification] ⏳ PENDING: %s", orderID)
			// TODO: mark order as awaiting payment

		case "refund", "partial_refund":
			log.Printf("[Notification] 💸 REFUNDED: %s", orderID)
			// TODO: handle refund
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

// ─── Transaction Status ───────────────────────────────────────────────────────

// GET /status/:order_id
func TransactionStatus(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		if orderID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
			return
		}

		var client coreapi.Client
		client.New(cfg.MidtransServerKey, coreEnv(cfg))

		resp, err := client.CheckTransaction(orderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    resp,
		})
	}
}

// ─── Cancel Transaction ───────────────────────────────────────────────────────

// POST /cancel/:order_id
func CancelTransaction(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		if orderID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})
			return
		}

		var client coreapi.Client
		client.New(cfg.MidtransServerKey, coreEnv(cfg))

		resp, err := client.CancelTransaction(orderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Transaction cancelled",
			"data":    resp,
		})
	}
}
