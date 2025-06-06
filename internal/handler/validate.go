package handler

import (
	"github.com/gin-gonic/gin"
	"log-receiver/internal/usecase"
	"log-receiver/pkg/logger"
	"log-receiver/pkg/middleware"
)

type validateService struct {
	logger logger.Logger

	usecaseValidator usecase.Validator
	isTestPem        bool
}

func NewValidateService(logger logger.Logger, route gin.IRouter, usecaseValidator usecase.Validator, isTestPem bool) {
	handler := validateService{
		logger:           logger,
		usecaseValidator: usecaseValidator,
		isTestPem:        isTestPem,
	}
	handler.handleRoute(route)
	return
}

func (h validateService) handleRoute(route gin.IRouter) {
	route.POST("/validate_token/:productCode", h.validateToken)
}

func (h validateService) validateToken(c *gin.Context) {
	payload, code, err := validateIDPToken(h.logger, c, h.usecaseValidator, h.isTestPem)
	if err != nil {
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.Set(middleware.CtxKeyIDPTokenPayload, payload)
	c.JSON(code, gin.H{"payload": payload})
}
