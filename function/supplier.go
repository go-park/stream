package function

type Supplier[T any] func() T

func (fn Supplier[T]) Get() T {
	return fn()
}
