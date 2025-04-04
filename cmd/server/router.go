package server

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"operation-borderless/internal/api"
)

func DefineRoutes(handler *api.Handler) *gin.Engine {
	log.Println("Routes defined")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	{
		router.GET("/", handler.Home())
	}

	r := router.Group("/api/v1")

	{
		r.POST("/create-wallet", handler.CreateWallet())
		r.POST("/deposit/:userID", handler.Deposit())
		r.POST("/convert/:userID", handler.Swap())
		r.POST("/transfer/:userID", handler.Transfer())
		r.GET("/wallet/:userID", handler.GetUserWallets())
		r.GET("transactions/:userID", handler.GetUserTransactions())
	}

	return router
}

func SetupRouter(h *api.Handler) *gin.Engine {
	log.Println("Router setup")
	r := DefineRoutes(h)

	return r
}
