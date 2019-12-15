package di

import "reflect"

// Scope defines the methods that any lifetime scope must implement.  While
// three Lifetime's exist, there are only two actual implementations of this
// interface - InstancePerContainer and InstancePerDependency.  This is because
// the Singleton lifetime simply directs container Forks to copy over the
// existing InstancePerContainer scopes to the new container.
type Scope interface {
	Resolve(c Container) (reflect.Value, error)
	GetUnderlyingResolver() *Resolver
	GetLifetime() Lifetime
}
