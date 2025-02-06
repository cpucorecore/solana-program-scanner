package types

type Event string

const (
	Buy    Event = "buy"
	Sell   Event = "sell"
	Add    Event = "add"
	Remove Event = "remove"
)
