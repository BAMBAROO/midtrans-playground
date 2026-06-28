package main

import (
	"log"
	"midtrans-tester/config"
	"midtrans-tester/handlers"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Open or create log file
	logFile, err := os.OpenFile("logs/dev.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Redirect standard logger
	log.SetOutput(logFile)

	// Redirect Gin logger
	gin.DefaultWriter = logFile
	gin.DefaultErrorWriter = logFile

	// Redirect os.Stdout and os.Stderr
	os.Stdout = logFile
	os.Stderr = logFile

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, reading from environment")
	}

	cfg := config.Load()

	r := gin.Default()

	// Serve static HTML tester
	r.StaticFile("/", "./static/index.html")
	// r.StaticFile("/testing", "./static/testing.html")

	// ─── Snap ─────────────────────────────────────────────────────────────────
	snap := r.Group("/snap")
	{
		snap.POST("/create", handlers.SnapCreate(cfg))
	}

	// ─── Core API ─────────────────────────────────────────────────────────────
	core := r.Group("/core")
	{
		// Bank Transfer
		core.POST("/bank-transfer/bca", handlers.BankTransferBCA(cfg))
		core.POST("/bank-transfer/bni", handlers.BankTransferBNI(cfg))
		core.POST("/bank-transfer/bri", handlers.BankTransferBRI(cfg))
		core.POST("/bank-transfer/mandiri", handlers.BankTransferMandiri(cfg))
		core.POST("/bank-transfer/permata", handlers.BankTransferPermata(cfg))
		core.POST("/bank-transfer/cimb", handlers.BankTransferCIMB(cfg))

		// E-Wallet
		core.POST("/ewallet/gopay", handlers.EWalletGoPay(cfg))
		core.POST("/ewallet/shopeepay", handlers.EWalletShopeePay(cfg))
		core.POST("/ewallet/dana", handlers.EWalletDANA(cfg))
		core.POST("/ewallet/ovo", handlers.EWalletOVO(cfg))

		// QRIS
		core.POST("/qris", handlers.QRIS(cfg))

		// Credit Card
		core.POST("/credit-card", handlers.CreditCard(cfg))

		// Convenience Store (Indomaret, Alfamart)
		core.POST("/cstore/indomaret", handlers.CStoreIndomaret(cfg))
		core.POST("/cstore/alfamart", handlers.CStoreAlfamart(cfg))

		// Akulaku
		core.POST("/akulaku", handlers.Akulaku(cfg))

		// Kredivo
		core.POST("/kredivo", handlers.Kredivo(cfg))

		// UOB EzPay
		core.POST("/uob-ezpay", handlers.UOBEzPay(cfg))
	}

	// ─── Notification / Webhook ───────────────────────────────────────────────
	r.POST("/notification", handlers.Notification(cfg))
	r.POST("/", handlers.Notification(cfg)) // fallback if merchant dashboard sets root URL

	// ─── Transaction Status ───────────────────────────────────────────────────
	r.GET("/status/:order_id", handlers.TransactionStatus(cfg))
	r.POST("/cancel/:order_id", handlers.CancelTransaction(cfg))

	log.Printf("Server running on :8080")
	log.Printf("Open http://localhost:8080 to start testing")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
