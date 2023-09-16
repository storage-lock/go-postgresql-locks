package postgresql_locks

import (
	"context"
	"database/sql"
	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

var (
	globalUriLockFactory = storage_lock_factory.NewStorageLockFactoryBeanFactory[string, *sql.DB]()
)

// NewLockByDSN 从DSN创建锁
func NewLockByDSN(ctx context.Context, dsn string, lockId string) (*storage_lock.StorageLock, error) {
	init, err := GetLockFactoryByDSN(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return init.CreateLock(lockId)
}

func NewLockByDSNWithOptions(ctx context.Context, dsn string, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	init, err := GetLockFactoryByDSN(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return init.CreateLockWithOptions(options)
}

func GetLockFactoryByDSN(ctx context.Context, dsn string) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
	return globalUriLockFactory.GetOrInit(ctx, dsn, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
		connectionManager := postgresql_storage.NewPostgresqlConnectionGetterFromDSN(dsn)
		options := postgresql_storage.NewPostgresqlStorageOptions().SetConnectionManager(connectionManager)
		storage, err := postgresql_storage.NewPostgresqlStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		return storage_lock_factory.NewStorageLockFactory[*sql.DB](storage, connectionManager), nil
	})
}
