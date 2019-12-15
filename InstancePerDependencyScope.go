package di

import "reflect"

// InstancePerDependencyScope calls the embedded resolver every time Resolve
// is called.
type InstancePerDependencyScope struct {
	resolver *Resolver
}

var _ Scope = &InstancePerDependencyScope{}

// NewInstancePerDependencyScope creates a new InstancePerDependency lifetime
// scope using the provided resolver.
func NewInstancePerDependencyScope(resolver *Resolver) *InstancePerDependencyScope {
	return &InstancePerDependencyScope{
		resolver: resolver,
	}
}

// Resolve calls the embedded resolver exactly once each time it is called.
func (s *InstancePerDependencyScope) Resolve(c Container) (reflect.Value, error) {
	return s.resolver.Resolve(c)
}

// GetUnderlyingResolver simply returns the underlying resolver for this scope.
func (s *InstancePerDependencyScope) GetUnderlyingResolver() *Resolver {
	return s.resolver
}

// GetLifetime simply returns the lifetime for this scope.
func (s *InstancePerDependencyScope) GetLifetime() Lifetime {
	return InstancePerDependency
}
