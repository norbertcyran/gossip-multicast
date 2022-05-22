package gossip

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Config struct {
	Fanout               int
	Transport            Transport
	GossipInterval       int
	RetransmitMultiplier int
	Neighbours           []string
	Tracer               Tracer
}

type Service struct {
	t              Transport
	fanout         int
	gossipInterval time.Duration

	neighbours []string

	messages             map[string]*Message
	retransmitMultiplier int
	mLock                sync.Mutex
	mCodec               MessageCodec

	tracer Tracer
}

func StartService(config *Config) (*Service, error) {
	interval := time.Duration(config.GossipInterval) * time.Millisecond
	s := &Service{
		t:                    config.Transport,
		fanout:               config.Fanout,
		gossipInterval:       interval,
		retransmitMultiplier: config.RetransmitMultiplier,
		neighbours:           config.Neighbours,
		mCodec:               &B64Codec{},
		messages:             make(map[string]*Message),
	}
	go s.handlePackets()
	go s.handleGossip()

	s.maybeTrace(ServiceStarted)
	return s, nil
}

func (s *Service) handlePackets() {
	for {
		select {
		case packet := <-s.t.PacketCh():
			encMsgs := bytes.Split(packet.Content, []byte(" "))
			for _, msg := range encMsgs {
				content, err := s.mCodec.Decode(msg)
				if err != nil {
					fmt.Printf("Invalid message format: %v, dropping...\n", err)
					continue
				} else {
					fmt.Printf("Received message: %q from: %q\n", string(content), packet.Source)
				}
				sum := md5.Sum(content)
				enc := hex.EncodeToString(sum[:])
				if _, ok := s.messages[enc]; !ok {
					s.messages[enc] = &Message{
						content:     content,
						retransmits: 0,
					}
					s.maybeTrace(ReceivedMessage)
				} else {
					s.maybeTrace(DuplicatedMessage)
				}
			}
		}
	}
}

func (s *Service) handleGossip() {
	for {
		select {
		case <-time.After(s.gossipInterval):
			s.gossipRound()
		}
	}
}

func (s *Service) gossipRound() {
	nodes := randomSample(s.fanout, s.neighbours)
	msgs := s.selectToSend()
	payload := s.generatePayload(msgs)
	if len(payload) == 0 {
		return
	}
	for _, node := range nodes {
		fmt.Printf("Sending gossip to: %q\n", node)
		if err := s.t.WriteTo(node, payload); err != nil {
			fmt.Printf("Sending the message to %q failed: %v\n", node, err)
		}
	}
}

func (s *Service) selectToSend() []*Message {
	msgs := make([]*Message, 0, len(s.messages))
	limit := s.retransmitLimit()
	for _, msg := range s.messages {
		if msg.retransmits < limit {
			msgs = append(msgs, msg)
			msg.retransmits++
		}
	}
	return msgs
}

func (s *Service) retransmitLimit() int {
	return retransmitLimit(s.retransmitMultiplier, s.fanout, len(s.neighbours))
}

func (s *Service) generatePayload(msgs []*Message) []byte {
	enc := make([][]byte, 0, len(msgs))
	for _, msg := range msgs {
		enc = append(enc, s.mCodec.Encode(msg.content))
	}
	return bytes.Join(enc, []byte(" "))
}

func (s *Service) maybeTrace(evt EventType) {
	if s.tracer != nil {
		s.tracer.Trace(evt)
	}
}
