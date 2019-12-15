package di

import (
	"fmt"
	"reflect"
	"sync"
)

// DefaultContainer is the default (and currently, only) implementation of
// Container.
type DefaultContainer struct {
	resolverScopes map[reflect.Type]Scope
	mx             *sync.RWMutex
}

var _ Container = &DefaultContainer{}

// NewContainer creates and returns a new, empty container.
func NewContainer() Container {
	return &DefaultContainer{
		resolverScopes: make(map[reflect.Type]Scope),
		mx:             &sync.RWMutex{},
	}
}

// Register registers a resolver function into the container with the provided
// lifetime.
func (c *DefaultContainer) Register(lifetime Lifetime, f interface{}) {
	resolver, err := NewResolver(f)
	if err != nil {
		panic(err)
	}

	var scope Scope
	if lifetime == InstancePerDependency {
		scope = NewInstancePerDependencyScope(resolver)
	} else {
		scope = NewInstancePerContainerScope(resolver, lifetime)
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	returnedInterfaceType := resolver.GetReturnedInterfaceType()
	c.resolverScopes[returnedInterfaceType] = scope
}

// Resolve resolves from the container into the provided interface pointers.
// Resolve will return an errors that arise will resolving dependencies.  If an
// error ocurrs, no guarantees are made around what is and is not resolved.
func (c *DefaultContainer) Resolve(dependencies ...interface{}) error {
	for _, dependency := range dependencies {
		err := c.resolveSingle(dependency)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *DefaultContainer) resolveSingle(into interface{}) error {
	interfacePointerType := reflect.TypeOf(into)
	if interfacePointerType == nil || interfacePointerType.Kind() != reflect.Ptr || interfacePointerType.Elem().Kind() != reflect.Interface {
		return fmt.Errorf("expected parameter like `*T` where T is an interface, but was `%v`", interfacePointerType)
	}

	interfaceType := interfacePointerType.Elem()

	c.mx.RLock()
	defer c.mx.RUnlock()

	scope, ok := c.resolverScopes[interfaceType]
	if !ok {
		return fmt.Errorf("the type `%v` does not have a resolver in this container", interfaceType)
	}

	v, err := scope.Resolve(c)
	if err != nil {
		return err
	}

	reflect.ValueOf(into).Elem().Set(v)
	return nil
}

// Fork forks the container into a new container lifetime scope.
func (c *DefaultContainer) Fork() Container {
	newContainer := &DefaultContainer{
		resolverScopes: make(map[reflect.Type]Scope),
		mx:             &sync.RWMutex{},
	}

	instancePerContainerScopes := []Scope{}

	for interfaceType, resolverScope := range c.resolverScopes {
		if resolverScope.GetLifetime() == InstancePerContainer {
			instancePerContainerScopes = append(instancePerContainerScopes, resolverScope)
			continue
		}

		newContainer.resolverScopes[interfaceType] = resolverScope
	}

	for _, instancePerContainerScope := range instancePerContainerScopes {
		resolver := instancePerContainerScope.GetUnderlyingResolver()
		interfaceType := resolver.GetReturnedInterfaceType()
		newInstancePerContainerScope := NewInstancePerContainerScope(resolver, InstancePerContainer)

		newContainer.resolverScopes[interfaceType] = newInstancePerContainerScope
	}

	return newContainer
}
