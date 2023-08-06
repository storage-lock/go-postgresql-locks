package postgresql_locks

import (
	"context"
	"database/sql"
	storage_lock_factory "github.com/storage-lock/go-storage-lock-factory"
)

var DefaultPostgresqlStorageLockFactory *storage_lock_factory.StorageLockFactory[*sql.DB]

func Init(ctx context.Context, dsn string, schema ...string) error {
	factory, err := NewPostgresqlFactory(ctx, dsn, schema...)
	if err != nil {
		return err
	}
	DefaultPostgresqlStorageLockFactory = factory
	return nil
}
