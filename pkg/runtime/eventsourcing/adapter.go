package eventsourcing

// Adapter is the interface for message buses.
type Adapter interface {
	Load()
	Apply() error
}
