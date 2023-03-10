package state

import "event-carried-state-transfer/schema"

type State struct {
	Name string
	Data schema.Person
}

type StatefulComponent interface {
	Dispatch(event schema.Event)
	AddListener(listener func(schema.Event))
	GetState() State
	GetListeners() []func(schema.Event)
}

type statefulComponent struct {
	state     State
	listeners []func(schema.Event)
	events    chan schema.Event
}

func NewStatefulComponent() StatefulComponent {

	component := &statefulComponent{
		events: make(chan schema.Event),
	}
	go component.handleEvents()
	return component
}

func (s *statefulComponent) Dispatch(event schema.Event) {

	s.events <- event
}

func (s *statefulComponent) AddListener(listener func(schema.Event)) {
	s.listeners = append(s.listeners, listener)
}

func (s *statefulComponent) GetState() State {
	return s.state
}

func (s *statefulComponent) GetListeners() []func(schema.Event) {
	return s.listeners
}

func (s *statefulComponent) handleEvents() {
	defer close(s.events)
	for event := range s.events {
		s.state = State{Name: event.Name, Data: event.Data}
		for _, listener := range s.listeners {
			listener(event)
		}
	}
}
