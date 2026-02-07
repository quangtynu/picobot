package bus

import "time"

// InboundMessage represents an incoming message to the agent.
type InboundMessage struct {
	Channel   string
	SenderID  string
	ChatID    string
	Content   string
	Timestamp time.Time
	Media     []string
	Metadata  map[string]interface{}
}

// OutboundMessage represents a message produced by the agent.
type OutboundMessage struct {
	Channel  string
	ChatID   string
	Content  string
	ReplyTo  string
	Media    []string
	Metadata map[string]interface{}
}

// MessageBus provides simple buffered channels for inbound/outbound messages.
type MessageBus struct {
	Inbound  chan InboundMessage
	Outbound chan OutboundMessage
}

// NewMessageBus constructs a new bus with the given buffer size.
func NewMessageBus(buffer int) *MessageBus {
	return &MessageBus{
		Inbound:  make(chan InboundMessage, buffer),
		Outbound: make(chan OutboundMessage, buffer),
	}
}

// Close closes the channels.
func (m *MessageBus) Close() {
	close(m.Inbound)
	close(m.Outbound)
}
