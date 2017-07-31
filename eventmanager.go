package steamprotocol

// EventHandler used to catch and handle event.
// It takes event object as argument
type EventHandler func(interface{}) error

// PacketHandler used to catch and handle Packet.
// It takes Packet as argument
type PacketHandler func(*Packet) error

// NewEventManager initialize new instance of EventManager.
func NewEventManager() *EventManager {
	return &EventManager{}
}

// EventManager is used to fire and broadcast events and packets to handlers.
type EventManager struct {
	eventHandlers  []EventHandler
	packetHandlers []PacketHandler
}

// OnEvent add Event handler.
// Any event will be received by all handlers.
// Use select via interface casting to catch events you are interested in.
func (m *EventManager) OnEvent(h EventHandler) {
	m.eventHandlers = append(m.eventHandlers, h)
}

// OnPacket add Packet handler.
// Any packet will be received by all handlers.
// Use select via EMsg to catch packets you are interested in.
func (m *EventManager) OnPacket(h PacketHandler) {
	m.packetHandlers = append(m.packetHandlers, h)
}

// FireEvent broadcast event to handlers
func (m *EventManager) FireEvent(e interface{}) error {
	for _, handler := range m.eventHandlers {
		err := handler(e)
		if err != nil {
			return err
		}
	}

	return nil
}

// FirePacket broadcast packet to handlers
func (m *EventManager) FirePacket(p *Packet) error {
	for _, handler := range m.packetHandlers {
		err := handler(p)
		if err != nil {
			return err
		}
	}

	return nil
}
