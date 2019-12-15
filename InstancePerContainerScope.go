package di

import (
	"reflect"
	"sync"
)

// InstancePerContainerScope calls the embedded resolver the first time Resolve
// is called.  Any subsequent calls to Resolve will simply return the instance
// that was first resolved.  InstancePerContainerScope is thread-safe.
type InstancePerContainerScope struct {
	resolver *Resolver
	mx       *sync.RWMutex
	value    *reflect.Value
	lifetime Lifetime
}

var _ Scope = &InstancePerContainerScope{}

// NewInstancePerContainerScope creates a new InstancePerContainer lifetime
// scope using the provided resolver.
func NewInstancePerContainerScope(resolver *Resolver, lifetime Lifetime) *InstancePerContainerScope {
	return &InstancePerContainerScope{
		resolver: resolver,
		mx:       &sync.RWMutex{},
		lifetime: lifetime,
	}
}

// Resolve calls the embedded resolver only once before it begins returning the
// same initially resolved value for each call.
func (s *InstancePerContainerScope) Resolve(c Container) (reflect.Value, error) {
	value := func() *reflect.Value {
		s.mx.RLock()
		defer s.mx.RUnlock()

		return s.value
	}()

	if value != nil {
		return *value, nil
	}

	s.mx.Lock()
	defer s.mx.Unlock()

	// If multiple routines are competing, s.value may already have a
	// new value.
	if s.value != nil {
		return *s.value, nil
	}

	v, err := s.resolver.Resolve(c)
	if err != nil {
		return reflect.Value{}, err
	}

	s.value = &v
	return v, nil
}

// GetUnderlyingResolver simply returns the underlying resolver for this scope.
func (s *InstancePerContainerScope) GetUnderlyingResolver() *Resolver {
	return s.resolver
}

// GetLifetime simply returns the lifetime for this scope.
func (s *InstancePerContainerScope) GetLifetime() Lifetime {
	return s.lifetime
}
