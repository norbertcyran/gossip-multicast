package gossip

type EventType int8

const (
	ReceivedMessage EventType = iota
	DuplicatedMessage
	ServiceStarted
)

type Tracer interface {
	Trace(evt EventType)
}
