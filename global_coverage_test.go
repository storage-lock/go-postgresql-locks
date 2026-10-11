package postgresql_locks

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	"github.com/stretchr/testify/assert"
)

func TestNewLockBySqlDB(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()

	// Init 先切 schema 再建表
	mock.ExpectExec("SET search_path").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
		WillReturnResult(sqlmock.NewResult(0, 0))

	ctx := context.Background()
	factory, err := GetLockFactoryBySqlDB(ctx, mockDB)
	assert.Nil(t, err)
	assert.NotNil(t, factory)

	// 缓存命中
	factory2, err := GetLockFactoryBySqlDB(ctx, mockDB)
	assert.Nil(t, err)
	assert.Equal(t, factory, factory2)

	lock, err := NewLockBySqlDB(ctx, mockDB, "test-pg-lock-sqldb")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-pg-lock-sqldb-opt")
	lock2, err := NewLockBySqlDBWithOptions(ctx, mockDB, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestNewLockBySqlDB_Errors(t *testing.T) {
	old := newPostgresqlStorage
	newPostgresqlStorage = func(ctx context.Context, options *postgresql_storage.PostgresqlStorageOptions) (*postgresql_storage.PostgresqlStorage, error) {
		return nil, errors.New("boom")
	}
	defer func() { newPostgresqlStorage = old }()

	mockDB, _, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()

	ctx := context.Background()
	_, err = GetLockFactoryBySqlDB(ctx, mockDB)
	assert.NotNil(t, err)
	_, err = NewLockBySqlDB(ctx, mockDB, "x")
	assert.NotNil(t, err)
	_, err = NewLockBySqlDBWithOptions(ctx, mockDB, storage_lock.NewStorageLockOptionsWithLockId("x"))
	assert.NotNil(t, err)
}

func TestNewLockByDSN_Local(t *testing.T) {
	dsn := "postgres://postgres:pgpass@127.0.0.1:15432/storage_lock_test?sslmode=disable"
	ctx := context.Background()

	probe, err := GetLockFactoryByDSN(ctx, dsn)
	if err != nil {
		t.Skipf("本机 sl-pg 不可用，跳过: %v", err)
	}
	assert.NotNil(t, probe)

	// 缓存命中
	factory2, err := GetLockFactoryByDSN(ctx, dsn)
	assert.Nil(t, err)
	assert.Equal(t, probe, factory2)

	lock, err := NewLockByDSN(ctx, dsn, "test-pg-lock-dsn")
	assert.Nil(t, err)
	assert.NotNil(t, lock)

	options := storage_lock.NewStorageLockOptionsWithLockId("test-pg-lock-dsn-opt")
	lock2, err := NewLockByDSNWithOptions(ctx, dsn, options)
	assert.Nil(t, err)
	assert.NotNil(t, lock2)

	assert.Nil(t, lock.Lock(ctx, "owner-1"))
	assert.Nil(t, lock.UnLock(ctx, "owner-1"))
}
