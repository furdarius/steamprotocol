package social

import (
	"bytes"

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
}

func NewModule(cl *steamprotocol.Client, eventManager *steamprotocol.EventManager) *Module {
	return &Module{
		cl:           cl,
		eventManager: eventManager,
	}
}

func (m *Module) Subscribe() {
	m.eventManager.OnEvent(m.handleEvent)
}

func (m *Module) handleEvent(e interface{}) error {
	switch event := e.(type) {
	case auth.SuccessfullyAuthenticatedEvent:
		m.handleSuccessfullyAuthenticatedEvent(event)
	}

	return nil
}

func (m *Module) handleSuccessfullyAuthenticatedEvent(e auth.SuccessfullyAuthenticatedEvent) {
	m.sessionID = e.SessionID
	m.steamID = e.SteamID

	m.SetUserOnline()
}

func (m *Module) SetUserOnline() error {
	responseHeader := messages.NewHeaderProto(steamprotocol.EMsg_ClientChangeStatus)
	responseHeader.Data.Steamid = proto.Uint64(m.steamID)
	responseHeader.Data.ClientSessionid = proto.Int32(m.sessionID)

	responseMsg := &protobuf.CMsgClientChangeStatus{
		PersonaState: proto.Uint32(uint32(steamprotocol.EPersonaState_Online)),
	}

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

	return nil
}
