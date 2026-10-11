package postgresql_locks

import (
	"context"
	"database/sql"
	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
	storage_lock "github.com/storage-lock/go-storage-lock"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

var (
	globalSqlDbLockFactory = storage_lock_factory.NewStorageLockFactoryBeanFactory[*sql.DB, *sql.DB]()
)

// newPostgresqlStorage 指向 postgresql_storage.NewPostgresqlStorage，包级变量以便单测注入构造失败，
// 覆盖 GetOrInit initFunc 的错误分支。
var newPostgresqlStorage = postgresql_storage.NewPostgresqlStorage

func NewLockBySqlDB(ctx context.Context, db *sql.DB, lockId string) (*storage_lock.StorageLock, error) {
	init, err := GetLockFactoryBySqlDB(ctx, db)
	if err != nil {
		return nil, err
	}
	return init.CreateLock(lockId)
}

func NewLockBySqlDBWithOptions(ctx context.Context, db *sql.DB, options *storage_lock.StorageLockOptions) (*storage_lock.StorageLock, error) {
	init, err := GetLockFactoryBySqlDB(ctx, db)
	if err != nil {
		return nil, err
	}
	return init.CreateLockWithOptions(options)
}

func GetLockFactoryBySqlDB(ctx context.Context, db *sql.DB) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
	return globalSqlDbLockFactory.GetOrInit(ctx, db, func(ctx context.Context) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {
		connectionManager := postgresql_storage.NewPostgresqlConnectionGetterFromSqlDb(db)
		options := postgresql_storage.NewPostgresqlStorageOptions().SetConnectionManager(connectionManager)
		storage, err := newPostgresqlStorage(ctx, options)
		if err != nil {
			return nil, err
		}
		return storage_lock_factory.NewStorageLockFactory[*sql.DB](storage, connectionManager), nil
	})
}
