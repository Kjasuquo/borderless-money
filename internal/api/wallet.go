package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"operation-borderless/internal/domain/dto"
	"operation-borderless/pkg/config"
	"operation-borderless/pkg/util"
)

func (h *Handler) CreateWallet() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := h.GetClientPublicIP(ctx)
		userAgent := ctx.GetHeader("User-Agent")

		var request dto.User
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("bad payload err: %v", err)})
			return
		}

		err := request.ConfirmEmailFormat()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("invalid email format err: %v", err)})
			return
		}

		var userID string

		user, err := h.Wallet.GetUserByEmail(ctx, strings.TrimSpace(request.Email))
		if err != nil {
			userID, err = h.Wallet.CreateWallet(ctx, strings.TrimSpace(request.Email))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("failed to create wallet: %v", err)})
				return
			}
		}

		userID = user.ID

		go h.AuditLogs(ctx, userID, userAgent, ip)

		ctx.JSON(http.StatusOK, gin.H{"userID": userID})
	}
}

func (h *Handler) Deposit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("userID")

		ip := h.GetClientPublicIP(ctx)
		userAgent := ctx.GetHeader("User-Agent")

		var request dto.DepositRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("bad payload err: %v", err)})
			return
		}

		if !dto.IsValidCurrency(request.Currency) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("valid currencies are : %s, %s, %s, %s", config.USD, config.EUR, config.NGN, config.XAF)})
			return
		}
		_, err := h.Wallet.GetUserWalletByCurrency(ctx, userID, request.Currency)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("failed to fetch user wallet err: %v", err)})
			return
		}

		// todo: this is meant to be pending status and changed and deposited in the webhook but that is not part of this scope

		trxID, err := h.Wallet.DepositFunds(ctx, dto.Transaction{
			SenderID:       userID,
			ReceiverID:     userID,
			FromAmount:     request.Amount,
			ToAmount:       request.Amount,
			FromCurrency:   request.Currency,
			ToCurrency:     request.Currency,
			Type:           string(config.Deposit),
			ConversionRate: 1,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Deposit err: %v", err)})
			return
		}

		go h.AuditLogs(ctx, userID, userAgent, ip)

		ctx.JSON(http.StatusOK, gin.H{"transactionID": trxID})
	}
}

func (h *Handler) Swap() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("userID")
		ip := h.GetClientPublicIP(ctx)
		userAgent := ctx.GetHeader("User-Agent")

		var request dto.TransferRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("bad payload err: %v", err)})
			return
		}

		if !dto.IsValidCurrency(request.FromCurrency) || !dto.IsValidCurrency(request.ToCurrency) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("valid currencies are : %s, %s, %s, %s", config.USD, config.EUR, config.NGN, config.XAF)})
			return
		}

		exchangeRate := 0.1

		if request.FromCurrency != request.ToCurrency {
			fromCurrency := strings.TrimPrefix(request.FromCurrency, "c")
			toCurrency := strings.TrimPrefix(request.ToCurrency, "c")

			rate, err := h.ExternalAPIs.GetPairExchangeRate(ctx, fromCurrency, toCurrency)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("valid currencies are : %s, %s, %s, %s", config.USD, config.EUR, config.NGN, config.XAF)})
				return
			}
			exchangeRate = rate.ConversionRate
		}

		if (request.FromAmount == 0 && request.ToAmount == 0) || (request.FromAmount != 0 && request.ToAmount != 0) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid amounts... specify either fromAmount only or toAmount only"})
			return
		}

		fromAmount, toAmount := request.FromAmount, request.ToAmount
		if request.FromAmount != 0 {
			toAmount = util.RoundToTwoDecimalPlaces(request.FromAmount * exchangeRate)
		} else {
			fromAmount = util.RoundToTwoDecimalPlaces(request.ToAmount / exchangeRate)
		}

		wallet, err := h.Wallet.GetUserWalletByCurrency(ctx, userID, request.FromCurrency)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("failed to fetch user wallet err: %v", err)})
			return
		}

		if wallet.Balance < fromAmount {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "insufficient balance"})
			return
		}

		trxID, err := h.Wallet.Transfer(ctx, dto.Transaction{
			SenderID:       userID,
			ReceiverID:     userID,
			FromAmount:     fromAmount,
			ToAmount:       toAmount,
			FromCurrency:   request.FromCurrency,
			ToCurrency:     request.ToCurrency,
			Type:           string(config.Swap),
			ConversionRate: exchangeRate,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("swap err: %v", err)})
			return
		}
		go h.AuditLogs(ctx, userID, userAgent, ip)

		ctx.JSON(http.StatusOK, gin.H{"transactionID": trxID})
	}
}

func (h *Handler) Transfer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("userID")

		ip := h.GetClientPublicIP(ctx)
		userAgent := ctx.GetHeader("User-Agent")

		var request dto.TransferRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("bad payload err: %v", err)})
			return
		}

		if !dto.IsValidCurrency(request.FromCurrency) || !dto.IsValidCurrency(request.ToCurrency) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("valid currencies are : %s, %s, %s, %s", config.USD, config.EUR, config.NGN, config.XAF)})
			return
		}

		receiver, err := h.Wallet.GetUserByEmail(ctx, strings.TrimSpace(request.ReceiverEmail))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("failed to fetch receiver wallet err: %v", err)})
			return
		}

		exchangeRate := 1.0
		if request.FromCurrency != request.ToCurrency {
			fromCurrency := strings.TrimPrefix(request.FromCurrency, "c")
			toCurrency := strings.TrimPrefix(request.ToCurrency, "c")

			rate, err := h.ExternalAPIs.GetPairExchangeRate(ctx, fromCurrency, toCurrency)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("valid currencies are : %s, %s, %s, %s", config.USD, config.EUR, config.NGN, config.XAF)})
				return
			}
			exchangeRate = rate.ConversionRate
		}

		if (request.FromAmount == 0 && request.ToAmount == 0) || (request.FromAmount != 0 && request.ToAmount != 0) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid amounts... specify either fromAmount only or toAmount only"})
			return
		}

		fromAmount, toAmount := request.FromAmount, request.ToAmount
		if request.FromAmount != 0 {
			toAmount = util.RoundToTwoDecimalPlaces(request.FromAmount * exchangeRate)
		} else {
			fromAmount = util.RoundToTwoDecimalPlaces(request.ToAmount / exchangeRate)
		}

		wallet, err := h.Wallet.GetUserWalletByCurrency(ctx, userID, request.FromCurrency)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("failed to fetch user wallet err: %v", err)})
			return
		}

		if wallet.Balance < fromAmount {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "insufficient balance"})
			return
		}

		trxID, err := h.Wallet.Transfer(ctx, dto.Transaction{
			SenderID:       userID,
			ReceiverID:     receiver.ID,
			FromAmount:     fromAmount,
			ToAmount:       toAmount,
			FromCurrency:   request.FromCurrency,
			ToCurrency:     request.ToCurrency,
			Type:           string(config.Transfer),
			ConversionRate: exchangeRate,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("transfer err: %v", err)})
			return
		}

		go h.AuditLogs(ctx, userID, userAgent, ip)

		ctx.JSON(http.StatusOK, gin.H{"transactionID": trxID})
	}
}

func (h *Handler) GetUserWallets() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("userID")
		ip := h.GetClientPublicIP(ctx)
		userAgent := ctx.GetHeader("User-Agent")

		// Use Goroutines for parallel execution
		var (
			user  dto.User
			rates dto.ExchangeRate
			err   error
		)
		errChan := make(chan error, 2)

		go func() {
			user, err = h.Wallet.GetUserByID(ctx, userID)
			errChan <- err
		}()

		go func() {
			rates, err = h.ExternalAPIs.GetUSDExchangeRate(ctx)
			errChan <- err
		}()

		for i := 0; i < 2; i++ {
			if e := <-errChan; e != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("fetch error: %v", e)})
				return
			}
		}

		currentRate := map[string]float64{
			string(config.USD): rates.ConversionRates.USD,
			string(config.NGN): rates.ConversionRates.NGN,
			string(config.EUR): rates.ConversionRates.EUR,
			string(config.XAF): rates.ConversionRates.XAF,
		}

		var totalUSDBalance float64
		for _, wallet := range user.Wallets {
			totalUSDBalance += wallet.Balance / currentRate[wallet.Currency]
		}

		user.TotalUSDBalance = util.RoundToTwoDecimalPlaces(totalUSDBalance)

		go h.AuditLogs(ctx, userID, userAgent, ip)

		ctx.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func (h *Handler) GetUserTransactions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("userID")

		ip := h.GetClientPublicIP(ctx)
		userAgent := ctx.GetHeader("User-Agent")

		transactions, err := h.Wallet.GetUserTransactions(ctx, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("fetch user err: %v", err)})
			return
		}

		go h.AuditLogs(ctx, userID, userAgent, ip)

		ctx.JSON(http.StatusOK, gin.H{"transactions": transactions})
	}
}
