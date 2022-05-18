package gossip

import (
	"time"
)

type Transport interface {
	PacketCh() <-chan *Packet
	WriteTo(addr string, message []byte) error
}

type Packet struct {
	Source    string
	Content   []byte
	Timestamp time.Time
}
