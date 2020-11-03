package model

// User represents user storage model.
type User struct {
	ID     string
	Status string
	Meta   map[string]interface{}
}

// UserFindFilter represents filter model for finding users.
type UserFindFilter struct {
	IDs          []string
	Statuses     []string
	MetaPatterns map[string]string
	Limit        *int64
	Offset       *int64
}
