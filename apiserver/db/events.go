package db

// EventDataChanged is called when data is changed in the database
type EventDataChanged struct {
	// The path to the data which is changed
	Path string
}
