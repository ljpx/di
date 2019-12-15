package di

import (
	"fmt"
	"reflect"
)

// Resolver validates and embeds a registration function.
type Resolver struct {
	resolverFunc          reflect.Value
	returnedInterfaceType reflect.Type
}

var containerReflectionType = reflect.TypeOf((*Container)(nil)).Elem()
var errorReflectionType = reflect.TypeOf((*error)(nil)).Elem()

// NewResolver creates a new resolver with the resolver function f.  It is
// required that f looks like `func (c di.Container) (T, error)` where `T` is an
// interface type.  An error will be returned if f is not a valid resolver
// function.
func NewResolver(f interface{}) (*Resolver, error) {
	t := reflect.TypeOf(f)

	if t.Kind() != reflect.Func || t.NumIn() != 1 || t.NumOut() != 2 {
		return nil, errorForInvalidResolverFunction(t)
	}

	if t.In(0) != containerReflectionType {
		return nil, errorForInvalidResolverFunction(t)
	}

	if t.Out(0).Kind() != reflect.Interface {
		return nil, errorForInvalidResolverFunction(t)
	}

	if t.Out(1) != errorReflectionType {
		return nil, errorForInvalidResolverFunction(t)
	}

	return &Resolver{
		resolverFunc:          reflect.ValueOf(f),
		returnedInterfaceType: t.Out(0),
	}, nil
}

// Resolve calls the resolver function, returning the resolved value as a
// reflect.Value, and the error returned from calling the resolver function.
// Exactly one of these return values will have a non-zero value.
func (r *Resolver) Resolve(c Container) (reflect.Value, error) {
	out := r.resolverFunc.Call([]reflect.Value{reflect.ValueOf(c)})

	if out[1].IsValid() && !out[1].IsNil() {
		err := out[1].Interface().(error)
		return out[0], fmt.Errorf("di resolve failure: %w", err)
	}

	return out[0], nil
}

// GetReturnedInterfaceType returns the type of the interface that is returned
// by the resolver function.
func (r *Resolver) GetReturnedInterfaceType() reflect.Type {
	return r.returnedInterfaceType
}

func errorForInvalidResolverFunction(t reflect.Type) error {
	return fmt.Errorf("expected a function that looked like `func (di.Container) (T, error)` but instead had `%v`", t)
}
