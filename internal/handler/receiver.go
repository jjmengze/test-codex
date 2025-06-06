package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"log-receiver/internal/domain/entity"
	"log-receiver/internal/usecase"
	"log-receiver/pkg/logger"
	"log-receiver/pkg/middleware"
)

type receiverService struct {
	logger           logger.Logger
	usecaseReceiver  usecase.Receiver
	usecaseValidator usecase.Validator
}

func NewReceiverService(logger logger.Logger, route gin.IRouter, usecaseReceiver usecase.Receiver, validator usecase.Validator) {
	handler := receiverService{
		logger:           logger,
		usecaseReceiver:  usecaseReceiver,
		usecaseValidator: validator,
	}
	handler.handleRoute(route)
	return
}

func (h receiverService) handleRoute(route gin.IRouter) {
	r := route.Group("/", ValidateTokenController(h.usecaseValidator), middleware.AssignV2Header(), middleware.CheckHeader(), middleware.CheckRequestBody())
	{
		// TODO: implement the HTTP handler for POST method
		// New employees should add the POST handler here
		activityLog := r.Group("/activity_log")
		activityLog.POST(":productCode", h.handleActivityLog)
	}
}

// TODO: Implement this controller to handle the activity log data
// This function should:
// 1. Extract data from the request
// 2. Create a putDataInput struct with the necessary fields
// 3. Call putData with the input
// 4. Return appropriate response
func (h receiverService) handleActivityLog(c *gin.Context) {
	ctx := c.Request.Context()

	input, err := h.extractContextFields(c)
	if err != nil {
		h.logger.ErrorF("Failed to extract context fields: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.usecaseReceiver.PutData(ctx, input)
	if err != nil {
		h.logger.WithContext(ctx).ErrorF("error putting data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error putting data",
		})
		return
	}
	c.JSON(200, gin.H{
		"message":  "Success",
		"trace_id": input.TraceID,
	})
}

func (h receiverService) extractContextFields(c *gin.Context) (entity.PutDataInput, error) {
	var input entity.PutDataInput
	var err error

	if input.ProductCode, err = h.getStringFromContext(c, "productCode"); err != nil {
		return input, err
	}
	if input.TraceID, err = h.getStringFromContext(c, "traceID"); err != nil {
		return input, err
	}
	if input.CustomerID, err = h.getStringFromContext(c, "customerID"); err != nil {
		return input, err
	}
	if input.Encoding, err = h.getStringFromContext(c, "encoding"); err != nil {
		return input, err
	}
	if input.SubType, err = h.getStringFromContext(c, "subType"); err != nil {
		return input, err
	}
	if input.SourceID, err = h.getStringFromContext(c, "sourceID"); err != nil {
		return input, err
	}
	return input, nil
}

func (h receiverService) getStringFromContext(c *gin.Context, key string) (string, error) {
	val, exists := c.Get(key)
	if !exists {
		err := fmt.Errorf("missing key in context: %s", key)
		h.logger.WarnF("error getting key %s: %v", key, err)
		return "", err
	}
	str, ok := val.(string)
	if !ok {
		err := fmt.Errorf("invalid type for key %s: expected string", key)
		h.logger.WarnF("error getting key %s: %v", key, err)
		return "", err
	}
	return str, nil
}
