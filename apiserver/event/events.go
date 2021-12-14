package event

import "log"

type Event interface{}

type Failure struct {
	// Original event
	Original Event

	// Listener points to whom failed
	Listener Listener
}

type Listener interface {
	// OnEvent is called when a new event is raised. Returning an error will raise a new event
	OnEvent(e Event) error
}

type Processor struct {
	listeners []Listener
}

func (p *Processor) AddListener(listener Listener) {
	p.listeners = append(p.listeners, listener)
}

// RaiseEvent to send a generic server event to various parts of the server, such as if
// a new user is added.
//
// The processing is done in a separate goroutine
func (p *Processor) RaiseEvent(evt Event) {
	go func() {
		for _, listener := range p.listeners {
			if err := listener.OnEvent(evt); err != nil {
				log.Printf("WARN: could not process event: %v\n", err)
				p.RaiseEvent(&Failure{
					Original: evt,
					Listener: listener,
				})
			}
		}
	}()
}
