package data

import (
	"fmt"
	"sync"

	"github.com/Icikowski/GoosyMock/model"
	"github.com/rs/zerolog"
)

// RoutesStore is a store implementation for routes
type RoutesStore struct {
	log         zerolog.Logger
	content     map[string]model.Route
	subscribers map[string](chan struct{})

	mux sync.Mutex
}

// NewRoutesStore creates new RoutesStore for runtime
func NewRoutesStore(log zerolog.Logger) *RoutesStore {
	return &RoutesStore{
		log:         log,
		content:     map[string]model.Route{},
		subscribers: map[string](chan struct{}){},
		mux:         sync.Mutex{},
	}
}

func (s *RoutesStore) notify() {
	s.log.Debug().Int("count", len(s.subscribers)).Msg("notifying subscribers about routes update")
	for _, subscriber := range s.subscribers {
		subscriber := subscriber
		go func() {
			subscriber <- struct{}{}
		}()
	}
}

// GetAll implements SubscribableStore
func (s *RoutesStore) GetAll() map[string]model.Route {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.content
}

// Get implements SubscribableStore
func (s *RoutesStore) Get(path string) (model.Route, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	var err error
	route, ok := s.content[path]
	if !ok {
		err = fmt.Errorf("route with given path does not exist: %s", path)
	}

	return route, err
}

// Set implements SubscribableStore
func (s *RoutesStore) Set(routes map[string]model.Route) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	defer s.notify()
	s.content = routes
	return nil
}

// Insert implements SubscribableStore
func (s *RoutesStore) Insert(path string, route model.Route) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, exists := s.content[path]; exists {
		return fmt.Errorf("route with given path already exists: %s", path)
	}
	s.content[path] = route
	defer s.notify()
	return nil
}

// Update implements SubscribableStore
func (s *RoutesStore) Update(path string, route model.Route) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, exists := s.content[path]; !exists {
		return fmt.Errorf("route with given path does not exist: %s", path)
	}
	s.content[path] = route
	defer s.notify()
	return nil
}

// Upsert implements SubscribableStore
func (s *RoutesStore) Upsert(path string, route model.Route) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.content[path] = route
	defer s.notify()
	return nil
}

// Delete implements SubscribableStore
func (s *RoutesStore) Delete(path string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if _, exists := s.content[path]; !exists {
		return fmt.Errorf("route with given path does not exist: %s", path)
	}
	delete(s.content, path)
	defer s.notify()
	return nil
}

// DeleteAll implements SubscribableStore
func (s *RoutesStore) DeleteAll() error {
	s.mux.Unlock()
	defer s.mux.Lock()

	s.content = map[string]model.Route{}
	defer s.notify()
	return nil
}

// Count implements SubscribableStore
func (s *RoutesStore) Count() int {
	return len(s.content)
}

// Subscribe implements SubscribableStore
func (s *RoutesStore) Subscribe(subscriberId string) <-chan struct{} {
	notify := make(chan struct{})
	s.subscribers[subscriberId] = notify
	return notify
}

// Unsubscribe implements SubscribableStore
func (s *RoutesStore) Unsubscribe(subscriberId string) {
	delete(s.subscribers, subscriberId)
}

var _ SubscribableStore[model.Route] = &RoutesStore{}
