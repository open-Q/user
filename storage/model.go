package storage

// User represents user storage model.
type User struct {
	ID     string
	Status string
	Meta   map[string][]byte
}
