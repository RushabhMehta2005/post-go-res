package store

// InMemStore defines the interface for a simple in-memory key-value store.
type InMemStore interface {
	// Get retrieves the value for the given key.
	// It returns the value and a boolean indicating whether the key was found.
	Get(key string) (string, bool)

	// Set stores the given value for the given key.
	Set(key, value string)

	// Delete removes the key from the store.
	// If the key does not exist, this should be a no-op.
	Delete(key string)
}
