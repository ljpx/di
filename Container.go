package di

// Container is a dependency-injection container.  New dependencies can be
// registered into the container using Register.  Dependencies can be resolved
// from the container using Resolve.  The container can be forked into a new
// lifetime using Fork.
type Container interface {
	Register(lifetime Lifetime, f interface{})
	Resolve(dependencies ...interface{}) error
	Fork() Container
}
