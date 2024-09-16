package model

import "context"

type Model interface {
	GetID() string
}

type Repository[T Model] interface {
	// Fetch a single item by its ID
	GetByID(ctx context.Context, id string) (*T, error)

	// Fetch multiple items by their IDs
	GetByIDList(ctx context.Context, ids []string) ([]*T, error)

	// Create a new item
	Create(ctx context.Context, item *T) (*string, error)

	// Save (create or update) an item
	Save(ctx context.Context, item *T) error

	// Update an existing item
	Update(ctx context.Context, item *T) error

	// Delete an item by its ID
	Delete(ctx context.Context, id string) error

	// Delete multiple items by their IDs
	DeleteList(ctx context.Context, ids []string) error
}
