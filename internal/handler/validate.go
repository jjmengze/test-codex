package handler

import (
	"github.com/gin-gonic/gin"
	"log-receiver/internal/usecase"
	"log-receiver/pkg/logger"
)

type validateService struct {
	logger logger.Logger

	usecaseValidator usecase.Validator
}

func NewValidateService(logger logger.Logger, route gin.IRouter, usecaseValidator usecase.Validator) {
	handler := validateService{
		logger:           logger,
		usecaseValidator: usecaseValidator,
	}
	handler.handleRoute(route)
	return
}

func (h validateService) handleRoute(route gin.IRouter) {
	route.POST("/validate_token/:productCode", h.validateToken)
}

func (h validateService) validateToken(c *gin.Context) {
	payload, code, err := validateIDPToken(c, h.usecaseValidator)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.Set(ctxKeyIDPTokenPayload, payload)
	c.JSON(code, gin.H{"payload": payload})
}
