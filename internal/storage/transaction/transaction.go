package transaction

import "context"

// Manager executes a block of business logic within a single, atomic transaction.
// The context passed to fn should be used for all underlying database operations.
type Manager interface {
	// Do executes given function in a transaction.
	// The transaction will be committed if fn returns nil, or rolled back otherwise.
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
