package heartbeat

import (
	"bytes"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/furdarius/steamprotocol"
	"github.com/furdarius/steamprotocol/auth"
	"github.com/furdarius/steamprotocol/messages"
	"github.com/furdarius/steamprotocol/protobuf"
)

type Module struct {
	eventManager *steamprotocol.EventManager
	cl           *steamprotocol.Client
	steamID      uint64
	sessionID    int32
	errorCh      chan error
	doneCh       chan struct{}
}

func NewModule(cl *steamprotocol.Client, eventManager *steamprotocol.EventManager) *Module {
	return &Module{
		cl:           cl,
		eventManager: eventManager,
		errorCh:      make(chan error),
		doneCh:       make(chan struct{}),
	}
}

func (m *Module) ErrorChannel() <-chan error {
	return m.errorCh
}

func (m *Module) Subscribe() {
	m.eventManager.OnEvent(m.handleEvent)
}

func (m *Module) handleEvent(e interface{}) error {
	switch event := e.(type) {
	case auth.SuccessfullyAuthenticatedEvent:
		m.handleSuccessfullyAuthenticatedEvent(event)
	case auth.LoggedOffEvent:
		m.handleLoggedOffEvent()
	}

	return nil
}

func (m *Module) handleSuccessfullyAuthenticatedEvent(e auth.SuccessfullyAuthenticatedEvent) {
	m.sessionID = e.SessionID
	m.steamID = e.SteamID

	m.eventManager.FireEvent(HeartBeatStartingEvent{
		Timeout: e.Heartbeat,
	})

	ticker := time.NewTicker(e.Heartbeat)

	go m.heartbeatLoop(ticker)
}

func (m *Module) handleLoggedOffEvent() {
	close(m.doneCh)
}

func (m *Module) heartbeatLoop(ticker *time.Ticker) {
	// TODO: Use context for graceful shutdown
	for {
		select {
		case <-ticker.C:
			err := m.doTick()
			if err != nil {
				m.errorCh <- err
			}
		case <-m.doneCh:
			ticker.Stop()
			close(m.errorCh)

			return
		}
	}
}

func (m *Module) doTick() error {
	responseHeader := messages.NewHeaderProto(steamprotocol.EMsg_ClientHeartBeat)

	responseHeader.Data.Steamid = proto.Uint64(m.steamID)
	responseHeader.Data.ClientSessionid = proto.Int32(m.sessionID)

	responseMsg := &protobuf.CMsgClientHeartBeat{}

	buf := new(bytes.Buffer)

	err := responseHeader.Serialize(buf)
	if err != nil {
		return errors.Wrap(err, "failed to serialize header")
	}

	body, err := proto.Marshal(responseMsg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal response msg")
	}

	_, err = buf.Write(body)
	if err != nil {
		return errors.Wrap(err, "failed to append msg to buffer")
	}

	err = m.cl.Write(buf.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to write data")
	}

	err = m.eventManager.FireEvent(HeartBeatTickedEvent{})
	if err != nil {
		return errors.Wrap(err, "failed to fire heartbeat tick event")
	}

	return nil
}
