package test

// RouterWrapperTester is a purely test interface, so we can generate a mocked implementation using mockery which then
// we can use to check whether it was called during execution in the wrapped handler.
type RouterWrapperTester interface {
	CalledIt() string
}
