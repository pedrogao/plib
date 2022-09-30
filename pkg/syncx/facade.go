package syncx

type (
	// Set interface for concurrent or not
	Set interface {
		// Find the key. return `false` if not found
		Find(key Ordered) (Ordered, bool)

		// Contains Returns `true` if the set contains the key.
		Contains(key Ordered) bool

		// Insert a key to the set. If the set already has the key, return `false`
		Insert(key Ordered) bool

		// Remove the key from the set and return it.
		Remove(key Ordered) Ordered
	}

	Ordered interface {
		Equal(Ordered) bool

		Less(Ordered) bool
	}
)
