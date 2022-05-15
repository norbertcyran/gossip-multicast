package gossip

import (
	"fmt"
	"net"
	"time"
)

const udpBufSize = 64 * 1024

type UDPTransportConfig struct {
	BindAddr string
	BindPort int
}

type UDPTransport struct {
	config   *UDPTransportConfig
	listener *net.UDPConn
	packetCh chan *Packet
}

func NewUDPTransport(config *UDPTransportConfig) (*UDPTransport, error) {
	ip := net.ParseIP(config.BindAddr)
	addr := &net.UDPAddr{IP: ip, Port: config.BindPort}
	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to set up listener on %q port %d: %v",
			config.BindAddr,
			config.BindPort,
			err,
		)
	}
	fmt.Printf("Listening UDP on %s:%d...\n", config.BindAddr, config.BindPort)

	t := &UDPTransport{
		config:   config,
		packetCh: make(chan *Packet),
		listener: listener,
	}
	go t.listen()
	return t, nil
}

func (t *UDPTransport) PacketCh() <-chan *Packet {
	return t.packetCh
}

func (t *UDPTransport) WriteTo(addr net.Addr, message []byte) error {
	n, err := t.listener.WriteTo(message, addr)
	if err != nil {
		return err
	}
	fmt.Printf("Sent %d bytes\n", n)
	return nil
}

func (t *UDPTransport) listen() {
	for {
		buf := make([]byte, udpBufSize)
		n, addr, err := t.listener.ReadFrom(buf)
		if err != nil {
			fmt.Printf("Error reading from UDP socket: %v", err)
		}
		fmt.Printf("Read %d bytes\n", n)
		ts := time.Now()

		t.packetCh <- &Packet{
			Source:    addr,
			Content:   buf[:n],
			Timestamp: ts,
		}
	}
}
