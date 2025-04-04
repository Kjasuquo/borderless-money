package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"operation-borderless/internal/domain/dto"
	"operation-borderless/internal/domain/services"
	"operation-borderless/pkg/config"
	"operation-borderless/pkg/util"
	"strings"
	"time"
)

type Handler struct {
	ExternalAPIs services.ExternalAPIs
	Wallet       services.Services
	cfg          *config.Config
}

func NewHandler(ForexAPI services.ExternalAPIs, Wallet services.Services, cfg *config.Config) *Handler {
	return &Handler{
		ExternalAPIs: ForexAPI,
		Wallet:       Wallet,
		cfg:          cfg,
	}
}

func (h *Handler) Home() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "welcome to borderless money"})
	}
}

func (h *Handler) AuditLogs(ctx context.Context, userID, userAgent, ipAddress string) {
	country, err := h.ExternalAPIs.GetUserCountry(ctx, ipAddress)
	if err != nil {
		log.Println(fmt.Sprintf("error getting country: %v", err))
	}
	device, browser := util.ParseUserAgent(userAgent)
	auditLog := dto.AuditLog{
		UserID:    userID,
		IPAddress: ipAddress,
		Device:    device,
		Country:   country,
		Browser:   browser,
		Timestamp: time.Now(),
	}

	err = h.Wallet.CreateAuditLogs(ctx, auditLog)
	if err != nil {
		log.Println(fmt.Sprintf("error creating logs: %v", err))
	}

	log.Println(fmt.Sprintf("Audit Logs: %v", auditLog))
}

func (h *Handler) GetClientPublicIP(ctx *gin.Context) string {
	if ip := ctx.GetHeader("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}

	if ip := ctx.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}

	return ctx.ClientIP()
}
