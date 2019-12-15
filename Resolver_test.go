package di

import (
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/ljpx/test"
)

func TestResolverSuccess(t *testing.T) {
	// Arrange.
	resolver, err := NewResolver(func(c Container) (testInterface, error) {
		return newtestStruct("Hello, World!"), nil
	})
	test.That(t, err).IsNil()

	// Act.
	val, err := resolver.Resolve(newTestContainer())
	test.That(t, err).IsNil()

	inst, ok := val.Interface().(*testStruct)
	test.That(t, ok).IsTrue()

	// Assert.
	test.That(t, inst.Greeting()).IsEqualTo("Hello, World!")
}

func TestResolverSuccessConcurrently(t *testing.T) {
	// Arrange.
	resolver, err := NewResolver(func(c Container) (testInterface, error) {
		return newtestStruct("Hello, World!"), nil
	})
	test.That(t, err).IsNil()

	wg := &sync.WaitGroup{}
	wg.Add(3)

	closure := func() {
		_, err := resolver.Resolve(newTestContainer())
		test.That(t, err).IsNil()
		wg.Done()
	}

	// Act and Assert.
	go closure()
	go closure()
	go closure()

	wg.Wait()
}

func TestResolverGetReturnedInterfaceType(t *testing.T) {
	// Arrange.
	resolver, err := NewResolver(func(c Container) (testInterface, error) {
		return newtestStruct("Hello, World!"), nil
	})
	test.That(t, err).IsNil()

	// Act.
	x := resolver.GetReturnedInterfaceType()

	// Assert.
	test.That(t, x).IsEqualTo(reflect.TypeOf((*testInterface)(nil)).Elem())
}

func TestResolverResolveError(t *testing.T) {
	// Arrange.
	resolver, err := NewResolver(func(c Container) (testInterface, error) {
		return nil, errors.New("an error")
	})
	test.That(t, err).IsNil()

	// Act.
	_, err = resolver.Resolve(newTestContainer())

	// Assert.
	test.That(t, err).IsNotNil()
	test.That(t, err.Error()).IsEqualTo("di resolve failure: an error")
}

func TestResolverFuncNotFunc(t *testing.T) {
	// Arrange and Act.
	_, err := NewResolver(5)

	// Assert.
	test.That(t, err).IsNotNil()
}

func TestResolverFuncNotExactly1InputArgument(t *testing.T) {
	// Arrange and Act.
	_, err := NewResolver(func(c Container, n int) (testInterface, error) {
		return newtestStruct("Hello, World!"), nil
	})

	// Assert.
	test.That(t, err).IsNotNil()
}

func TestResolverFuncNotExactly2OutputArguments(t *testing.T) {
	// Arrange and Act.
	_, err := NewResolver(func(c Container) (testInterface, int, error) {
		return newtestStruct("Hello, World!"), 5, nil
	})

	// Assert.
	test.That(t, err).IsNotNil()
}

func TestResolverFuncArgumentNotContainerType(t *testing.T) {
	// Arrange and Act.
	_, err := NewResolver(func(n int) (testInterface, error) {
		return newtestStruct("Hello, World!"), nil
	})

	// Assert.
	test.That(t, err).IsNotNil()
}

func TestResolverFuncFirstReturnTypeNotInterface(t *testing.T) {
	// Arrange and Act.
	_, err := NewResolver(func(c Container) (*testStruct, error) {
		return newtestStruct("Hello, World!"), nil
	})

	// Assert.
	test.That(t, err).IsNotNil()
}

func TestResolverFuncSecondReturnTypeNotError(t *testing.T) {
	// Arrange and Act.
	_, err := NewResolver(func(c Container) (testInterface, int) {
		return newtestStruct("Hello, World!"), 5
	})

	// Assert.
	test.That(t, err).IsNotNil()
}

// -----------------------------------------------------------------------------

func newTestResolver() (*Resolver, error) {
	return NewResolver(func(c Container) (testInterface, error) {
		return newtestStruct("Hello, World!"), nil
	})
}

type testInterface interface {
	Greeting() string
}

type testStruct struct {
	greeting string
}

var _ testInterface = &testStruct{}

func newtestStruct(greeting string) *testStruct {
	return &testStruct{
		greeting: greeting,
	}
}

func (r *testStruct) Greeting() string {
	return r.greeting
}

type testContainer struct{}

var _ Container = &testContainer{}

func newTestContainer() *testContainer {
	return &testContainer{}
}

func (c *testContainer) Fork() Container {
	return &testContainer{}
}

func (c *testContainer) Register(lifetime Lifetime, f interface{}) {}

func (c *testContainer) Resolve(dependencies ...interface{}) error {
	return nil
}
