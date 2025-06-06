package usecase

import (
	"context"

	"log-receiver/internal/domain/entity"
	"log-receiver/pkg/logger"
)

type Validator interface {
	Validate(ctx context.Context, productCode string) (bool, error)
}

type validator struct {
	logger logger.Logger
}

func NewValidator(logger logger.Logger) Validator {
	return validator{logger: logger}
}

func (v validator) Validate(ctx context.Context, productCode string) (bool, error) {
	if productCode == "" {
		v.logger.WithContext(ctx).WarnF("productCode is empty")
		return false, nil
	}
	ok := isSupportedProduct(productCode)
	if !ok {
		v.logger.WithContext(ctx).WarnF("productCode %v is not supported", productCode)
		return false, nil
	}
	return true, nil
}

func isSupportedProduct(code string) bool {
	_, ok := entity.SupportedProductCode[code]
	return ok
}
