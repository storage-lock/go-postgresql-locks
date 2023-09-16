package postgresql_locks

import (
	"context"
	storage_lock_test_helper "github.com/storage-lock/go-storage-lock-test-helper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewLockByDSN(t *testing.T) {
	envName := "STORAGE_LOCK_POSTGRESQL_DSN"
	dsn := os.Getenv(envName)
	factory, err := GetLockFactoryByDSN(context.Background(), dsn)
	assert.Nil(t, err)

	storage_lock_test_helper.PlayerNum = 50
	storage_lock_test_helper.EveryOnePlayTimes = 10
	storage_lock_test_helper.TestStorageLock(t, factory)

}
