package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	mockLogger "log-receiver/mock/pkg/logger"
)

func TestValidatorValidate(t *testing.T) {
	ctx := context.Background()

	t.Run("empty productCode", func(t *testing.T) {
		mLogger := mockLogger.NewLogger(t)
		mLogger.On("WithContext", ctx).Return(mLogger).Once()
		mLogger.On("WarnF", "productCode is empty").Once()

		v := NewValidator(mLogger)
		ok, err := v.Validate(ctx, "")
		assert.NoError(t, err)
		assert.False(t, ok)
		mLogger.AssertExpectations(t)
	})

	t.Run("unsupported product", func(t *testing.T) {
		mLogger := mockLogger.NewLogger(t)
		mLogger.On("WithContext", ctx).Return(mLogger).Once()
		mLogger.On("WarnF", "productCode %v is not supported", "foo").Once()

		v := NewValidator(mLogger)
		ok, err := v.Validate(ctx, "foo")
		assert.NoError(t, err)
		assert.False(t, ok)
		mLogger.AssertExpectations(t)
	})

	t.Run("supported product", func(t *testing.T) {
		mLogger := mockLogger.NewLogger(t)

		v := NewValidator(mLogger)
		ok, err := v.Validate(ctx, "sao")
		assert.NoError(t, err)
		assert.True(t, ok)
	})
}

func TestIsSupportedProduct(t *testing.T) {
	assert.True(t, isSupportedProduct("sao"))
	assert.False(t, isSupportedProduct("foo"))
}
