package di

import (
	"fmt"
	"io"
	"sync"
	"testing"

	"github.com/ljpx/test"
)

func TestDefaultContainerSuccess(t *testing.T) {
	// Arrange.
	c := getPreparedContainerWithLifetime(InstancePerDependency)

	// Act.
	var x testInterface
	err := c.Resolve(&x)
	test.That(t, err).IsNil()

	// Assert.
	test.That(t, x.Greeting()).IsEqualTo("Hello, World!")
}

func TestDefaultContainerSuccessConcurrent(t *testing.T) {
	// Arrange.
	c := getPreparedContainerWithLifetime(Singleton)

	wg := &sync.WaitGroup{}
	wg.Add(3)

	closure := func() {
		var x testInterface
		err := c.Resolve(&x)
		test.That(t, err).IsNil()
		wg.Done()
	}

	// Act and Assert.
	go closure()
	go closure()
	go closure()

	wg.Wait()
}

func TestDefaultContainerResolveIntoNonPointerToInterface(t *testing.T) {
	// Arrange.
	c := getPreparedContainerWithLifetime(InstancePerDependency)

	// Act.
	var x testInterface
	err := c.Resolve(x)

	// Assert.
	test.That(t, err).IsNotNil()
}

func TestDefaultContainerUnknownInterfaceForResolve(t *testing.T) {
	// Arrange.
	c := getPreparedContainerWithLifetime(InstancePerDependency)

	// Act.
	var x io.Writer
	err := c.Resolve(&x)

	// Assert.
	test.That(t, err).IsNotNil()
}

func TestDefaultContainerResolveError(t *testing.T) {
	// Arrange.
	c := NewContainer()
	c.Register(InstancePerDependency, func(c Container) (testInterface, error) {
		return nil, fmt.Errorf("an error")
	})

	// Act.
	var x testInterface
	err := c.Resolve(&x)

	// Assert.
	test.That(t, err).IsNotNil()
	test.That(t, err.Error()).IsEqualTo("di resolve failure: an error")
}

func TestDefaultContainerFork(t *testing.T) {
	// Arrange.
	c1 := getPreparedContainerWithLifetime(InstancePerContainer)

	c2 := c1.Fork()

	// Act.
	var x1 testInterface
	err := c1.Resolve(&x1)
	test.That(t, err).IsNil()

	var x2 testInterface
	err = c1.Resolve(&x2)
	test.That(t, err).IsNil()

	var x3 testInterface
	err = c2.Resolve(&x3)
	test.That(t, err).IsNil()

	// Assert.
	test.That(t, x2).IsEqualTo(x1)
	test.That(t, x2).IsNotEqualTo(x3)
	test.That(t, x1).IsNotEqualTo(x3)
}

func getPreparedContainerWithLifetime(lifetime Lifetime) Container {
	c := NewContainer()
	c.Register(lifetime, func(c Container) (testInterface, error) {
		return newtestStruct("Hello, World!"), nil
	})

	return c
}
