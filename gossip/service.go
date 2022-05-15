package gossip

import (
	"crypto/md5"
	"fmt"
	"sync"
)

type Config struct {
	Fanout    int
	Transport Transport
}

type Service struct {
	t          Transport
	fanout     int
	Neighbours []*Node
	Messages   sync.Map
}

func StartService(config *Config) (*Service, error) {
	s := &Service{t: config.Transport, fanout: config.Fanout}

	go s.handlePackets()

	return s, nil
}

func (s *Service) handlePackets() {
	for {
		select {
		case packet := <-s.t.PacketCh():
			go s.echo(packet)
		}
	}
}

func (s *Service) echo(packet *Packet) {
	from := packet.Source
	msg := packet.Content

	if err := s.t.WriteTo(from, msg); err != nil {
		fmt.Printf("Error sending message: %v", err)
	} else {
		fmt.Println("sent msg")
	}
}

func (s *Service) HandleMessage(msg []byte) []byte {
	hash := md5.New()
	checksum := string(hash.Sum(msg))
	if _, loaded := s.Messages.LoadOrStore(checksum, &Message{
		Checksum: checksum,
		State:    Infective,
		Content:  string(msg),
	}); !loaded {
		return []byte("NEW")
	} else {
		return []byte("OLD")
	}
}
