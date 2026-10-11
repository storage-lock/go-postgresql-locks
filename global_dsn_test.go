package postgresql_locks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalUriLockFactory_InitiallyInitialized(t *testing.T) {
	assert.NotNil(t, globalUriLockFactory)
}

func TestNewLockByDSN_InvalidDsn(t *testing.T) {
	_, err := NewLockByDSN(context.Background(), "invalid-dsn", "test-lock")
	assert.Error(t, err)
}

func TestGetLockFactoryByDSN_InvalidDsn(t *testing.T) {
	_, err := GetLockFactoryByDSN(context.Background(), "invalid-dsn")
	assert.Error(t, err)
}

func TestNewLockByDSNWithOptions_InvalidDsn(t *testing.T) {
	_, err := NewLockByDSNWithOptions(context.Background(), "invalid-dsn", nil)
	assert.Error(t, err)
}

