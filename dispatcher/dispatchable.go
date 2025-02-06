package dispatcher

type Dispatchable[T any] interface {
	Id() string
	Dispatch(T) error
	Stop()
}
