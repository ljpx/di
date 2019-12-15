package di

import (
	"sync"
	"testing"

	"github.com/ljpx/test"
)

func TestInstancePerContainerScopeResolvesSameInstances(t *testing.T) {
	// Arrange.
	resolver, err := newTestResolver()
	test.That(t, err).IsNil()

	scope := NewInstancePerContainerScope(resolver, InstancePerContainer)

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
	test.That(t, inst1).IsEqualTo(inst2)
}

func TestInstancePerContainerScopeResolvesSameInstancesConcurrently(t *testing.T) {
	// Arrange.
	resolver, err := newTestResolver()
	test.That(t, err).IsNil()

	scope := NewInstancePerContainerScope(resolver, InstancePerContainer)

	instc := make(chan testInterface)
	wg := &sync.WaitGroup{}
	wg.Add(5)

	closure := func() {
		v, err := scope.Resolve(newTestContainer())
		test.That(t, err).IsNil()

		inst, ok := v.Interface().(testInterface)
		test.That(t, ok).IsTrue()

		instc <- inst
		wg.Done()
	}

	// Act.
	go closure()
	go closure()
	go closure()
	go closure()
	go closure()

	go func() {
		wg.Wait()
		close(instc)
	}()

	vs := []testInterface{}
	for inst := range instc {
		vs = append(vs, inst)
	}

	// Assert.
	for i := 0; i < len(vs); i++ {
		for j := 0; j < len(vs); j++ {
			if i == j {
				continue
			}

			test.That(t, vs[i]).IsEqualTo(vs[j])
		}
	}
}
