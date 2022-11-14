package data

import "io"

// Store is generic definition of key-value storage adapter
type Store[T any] interface {
	// GetAll returns all elements from the store
	GetAll() map[string]T
	// Get returns the element with given ID from the store
	Get(string) (T, error)
	// Set sets the contents of the store
	Set(map[string]T) error
	// Insert inserts the element with given ID into the store
	Insert(string, T) error
	// Update updates the element with given ID in the store
	Update(string, T) error
	// Upsert inserts or updates the element with given ID
	// into/in the store
	Upsert(string, T) error
	// Delete deletes the element with given ID from the store
	Delete(string) error
	// DeleteAll deletes all elements from the store
	DeleteAll() error
	// Count returns the number of elements in the store
	Count() int
}

// Subscribable is generic definition of subscribable object
type Subscribable interface {
	// Subscribe subscribes to the channel which notifies about
	// changes being made to the object with particular
	// subscriber's ID
	Subscribe(string) <-chan struct{}
	// Unsubscribe removes the subscription for notifications
	// of particular receiver
	Unsubscribe(string)
}

// SubscribableStore is a definition of store that can indicate
// whether the contents of the store has changed
type SubscribableStore[T any] interface {
	Store[T]
	Subscribable
}

// FilesStore is a definition of store that is capable of storing
// files and defines Close method for cleanup
type FilesStore[T any] interface {
	Store[T]
	io.Closer
	// InsertFile inserts the file into store with random UUID
	// as a key
	InsertFile(string, io.ReadCloser) (string, error)
	// UpdateFile updates the file item with given ID in the store
	UpdateFile(string, string, io.ReadCloser) error
	// UpsertFile inserts or updates the file item with given ID
	// into/in the store
	UpsertFile(string, string, io.ReadCloser) error
}
