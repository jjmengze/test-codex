package handler

import (
	"github.com/gin-gonic/gin"
	"log-receiver/internal/usecase"
	"log-receiver/pkg/logger"
	"log-receiver/pkg/middleware"
)

type httpHandler struct {
	logger logger.Logger
	app    *gin.Engine
}

func NewHttpHandler(logger logger.Logger, app *gin.Engine, receiver usecase.Receiver, validator usecase.Validator) *httpHandler {
	h := &httpHandler{
		logger: logger,
		app:    app,
	}
	h.app.GET("/health", func(c *gin.Context) {})

	api := h.app.Group("/api")
	api.Use(middleware.AssignTraceID())
	v2 := api.Group("/v2")

	NewReceiverService(logger, v2, receiver, validator)
	NewValidateService(logger, v2, validator)
	return h
}
