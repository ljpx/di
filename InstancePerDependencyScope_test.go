package di

import (
	"sync"
	"testing"

	"github.com/ljpx/test"
)

func TestInstancePerDependencyScopeResolvesDifferentInstances(t *testing.T) {
	// Arrange.
	resolver, err := newTestResolver()
	test.That(t, err).IsNil()

	scope := NewInstancePerDependencyScope(resolver)

	// Act.
	v1, err := scope.Resolve(newTestContainer())
	test.That(t, err).IsNil()

	v2, err := scope.Resolve(newTestContainer())
	test.That(t, err).IsNil()

	inst1, ok := v1.Interface().(testInterface)
	test.That(t, ok).IsTrue()

	inst2, ok := v2.Interface().(testInterface)
	test.That(t, ok).IsTrue()

	// Assert.
	test.That(t, inst1).IsNotEqualTo(inst2)
}

func TestInstancePerDependencyScopeResolvesDifferentInstancesConcurrently(t *testing.T) {
	// Arrange.
	resolver, err := newTestResolver()
	test.That(t, err).IsNil()

	scope := NewInstancePerDependencyScope(resolver)

	wg := &sync.WaitGroup{}
	wg.Add(3)

	closure := func() {
		v1, err := scope.Resolve(newTestContainer())
		test.That(t, err).IsNil()

		v2, err := scope.Resolve(newTestContainer())
		test.That(t, err).IsNil()

		inst1, ok := v1.Interface().(testInterface)
		test.That(t, ok).IsTrue()

		inst2, ok := v2.Interface().(testInterface)
		test.That(t, ok).IsTrue()

		test.That(t, inst1).IsNotEqualTo(inst2)
		wg.Done()
	}

	// Act and Assert.
	go closure()
	go closure()
	go closure()

	wg.Wait()
}
