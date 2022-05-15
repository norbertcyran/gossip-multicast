package gossip

import (
	"net"
	"time"
)

type Transport interface {
	PacketCh() <-chan *Packet
	WriteTo(addr net.Addr, message []byte) error
}

type Packet struct {
	Source    net.Addr
	Content   []byte
	Timestamp time.Time
}
