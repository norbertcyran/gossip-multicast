package gossip

type EventType int8

const (
	ReceivedMessage EventType = iota
	DuplicatedMessage
)

type Tracer interface {
	Trace(evt EventType)
}
