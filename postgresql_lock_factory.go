package postgresql_locks

import (
	"context"
	"database/sql"
	postgresql_storage "github.com/storage-lock/go-postgresql-storage"
	"github.com/storage-lock/go-storage"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

func NewPostgresqlFactory(ctx context.Context, dsn string, schema ...string) (*storage_lock_factory.StorageLockFactory[*sql.DB], error) {

	if len(schema) == 0 {
		schema[0] = postgresql_storage.DefaultPostgreSQLStorageSchema
	}

	connectionManager := postgresql_storage.NewPostgreSQLConnectionGetterFromDSN(dsn)
	storageOptions := &postgresql_storage.PostgreSQLStorageOptions{
		ConnectionManager: connectionManager,
		TableName:         storage.DefaultStorageTableName,
		Schema:            schema[0],
	}
	storage, err := postgresql_storage.NewPostgreSQLStorage(ctx, storageOptions)
	if err != nil {
		return nil, err
	}

	return storage_lock_factory.NewStorageLockFactory[*sql.DB](storage, connectionManager), nil
}
