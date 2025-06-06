package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"log-receiver/internal/usecase"
	"log-receiver/pkg/logger"
	"log-receiver/pkg/middleware"
)

type httpHandler struct {
	logger logger.Logger
	app    *gin.Engine
}

func NewHttpHandler(logger logger.Logger, app *gin.Engine, receiver usecase.Receiver, validator usecase.Validator, isTestPem bool) *httpHandler {
	h := &httpHandler{
		logger: logger,
		app:    app,
	}
	h.app.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
		return
	})

	api := h.app.Group("/api")
	api.Use(middleware.AssignTraceID())
	v2 := api.Group("/v2")

	NewReceiverService(logger, v2, receiver, validator, isTestPem)
	NewValidateService(logger, v2, validator, isTestPem)
	return h
}
