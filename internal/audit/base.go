package audit

// Subscriber represents a component that receives audit events.
// Each subscriber has a unique ID and can process events via Update.
// Subscribers are expected to handle errors internally and return an error
// only if they should be unregistered from the Audit dispatcher.
type Subscriber interface {
	// Update processes an audit event.
	// Returns an error if the subscriber encountered a fatal issue and should be unregistered.
	Update(*Event) error

	// GetID returns a unique identifier for the subscriber.
	// This ID is used to register and unregister the subscriber from the Audit.
	GetID() string

	// Close allows the subscriber to clean up resources (e.g., flush buffers, close files).
	// Returns an error if cleanup fails.
	Close() error
}

// Audit is a publish-subscribe dispatcher for audit events.
// It manages a collection of subscribers and broadcasts events to them.
// If a subscriber returns an error during event delivery, it is automatically unregistered.
type Audit struct {
	subscribers map[string]Subscriber // subscribers holds registered subscribers keyed by their ID.
}

func New() *Audit {
	return &Audit{
		subscribers: make(map[string]Subscriber),
	}
}

func (a *Audit) Register(s Subscriber) {
	a.subscribers[s.GetID()] = s
}

func (a *Audit) Unregister(s Subscriber) {
	delete(a.subscribers, s.GetID())
}

func (a *Audit) Publish(e *Event) {
	var closedSubs []Subscriber
	for _, s := range a.subscribers {
		if err := s.Update(e); err != nil {
			closedSubs = append(closedSubs, s)
		}
	}
	for _, s := range closedSubs {
		a.Unregister(s)
	}
}
