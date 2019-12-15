package di

// Lifetime defines the lifetime of a resolved dependency.
type Lifetime int

// Singleton is the largest lifetime.  Dependencies that are registered with
// Singleton will only ever be resolved once.  After the initial resolution, the
// same instance will be returned in all further resolutions, even across
// containers.
const Singleton Lifetime = 0

// InstancePerContainer behaves in much the same way as the Singleton lifetime,
// with the exception of forked containers.  When a container is forked,
// registrations made with the Singleton lifetime will persist in the new
// container.  Registrations made with InstancePerContainer will not.  This
// means that each container will resolve InstancePerContainer registrations
// up to one time each.
const InstancePerContainer Lifetime = 1

// InstancePerDependency is the smallest, and simplest, lifetime.  Every call to
// Resolve will resolve a new instance from the container.
const InstancePerDependency Lifetime = 2
