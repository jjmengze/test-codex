package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"log-receiver/internal/usecase"
	"log-receiver/pkg/auth"
	"log-receiver/pkg/logger"
	"log-receiver/pkg/middleware"
)

type validateService struct {
	logger logger.Logger

	usecaseValidator usecase.Validator
	isTestPem        bool
}

func NewValidateService(logger logger.Logger, route gin.IRouter, usecaseValidator usecase.Validator, pubKeyAbsPath string, isTestPem bool) {
	handler := validateService{
		logger:           logger,
		usecaseValidator: usecaseValidator,
		isTestPem:        isTestPem,
	}
	handler.handleRoute(route, pubKeyAbsPath, isTestPem)
	return
}

func (h validateService) handleRoute(route gin.IRouter, pubKeyAbsPath string, isTestPem bool) {
	route.POST("/validate_token/:productCode", middleware.ValidateTokenController(h.logger, h.usecaseValidator, pubKeyAbsPath, isTestPem), h.validateToken)
}

func (h validateService) validateToken(c *gin.Context) {
	ctx := c.Request.Context()

	val, ok := c.Get(middleware.CtxKeyIDPTokenPayload)
	if !ok {
		err := errors.New("missing ID P Token Payload")
		h.logger.WithContext(ctx).ErrorF(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	payload, ok := val.(*auth.IDPTokenPayload)
	if !ok {
		err := fmt.Errorf("type %t can't convert to IDP Token Payload", val)
		h.logger.WithContext(ctx).ErrorF(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payload": payload})
}
