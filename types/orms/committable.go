package orms

type Committable interface {
	TableName() string
}
