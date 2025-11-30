package audit

type Subscriber interface {
	Update(*Event) error
	GetID() string
	Close() error
}

type Audit struct {
	subscribers map[string]Subscriber
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
