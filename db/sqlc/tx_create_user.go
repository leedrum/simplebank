package db

import (
	"context"
)

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// CreateUserTxResult is the result of the CreateUser transaction
type CreateUserTxResult struct {
	User User
}

// CreateUserTx performs a money CreateUser from one account to the other.
// It creates the CreateUser, add account entries, and update accounts' balance within a database transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
