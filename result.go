package noor

type Result[T any] struct {
	Ok  T
	Err error
}

func Ok[T any](value T) Result[T] {
	return Result[T]{Ok: value}
}

func Err[T any](err error) Result[T] {
	return Result[T]{Err: err}
}

func (r Result[T]) Unwrap() (T, error) {
	return r.Ok, r.Err
}

func (r Result[T]) UnwrapOr(def T) T {
	if r.Err != nil {
		return def
	}
	return r.Ok
}

func (r Result[T]) UnwrapOrElse(f func() T) T {
	if r.Err != nil {
		return f()
	}
	return r.Ok
}

func (r Result[T]) UnwrapOrPanic() T {
	if r.Err != nil {
		panic(r.Err)
	}
	return r.Ok
}

func (r Result[T]) IsOk() bool {
	return r.Err == nil
}

func (r Result[T]) IsErr() bool {
	return r.Err != nil
}
