package dbx

import (
	"context"

	"gorm.io/gorm"
)

var transactionKey struct{}

// WithTransaction 中间件：将tx存入context
func WithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, transactionKey, tx)
}

// TransactionFromContext 从context获取tx
func TransactionFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionKey).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// InTransaction 在事务中执行函数
func InTransaction(ctx context.Context, db *gorm.DB, fn func(context.Context) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		txCtx := WithTransaction(ctx, tx)
		return fn(txCtx)
	})
}
