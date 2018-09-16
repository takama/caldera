package provider

// Transact defines transact provider type
type Transact interface {
	Commit() error
	Rollback() error
}
