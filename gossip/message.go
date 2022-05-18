package gossip

type Message struct {
	content     []byte
	retransmits int
}
