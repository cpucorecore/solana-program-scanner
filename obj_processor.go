package main

type ObjProcessor[T any] interface {
	id() string
	process(T) error
	done()
}
