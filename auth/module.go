package auth

import (
	"bytes"

	"time"

	"crypto/sha1"

	"github.com/furdarius/steamprotocol"
	"github.com/furdarius/steamprotocol/crypto"
	"github.com/furdarius/steamprotocol/messages"
	"github.com/furdarius/steamprotocol/protobuf"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// Log on with the given details. You must always specify username and
// password. For the first login, don't set an authcode or a hash and you'll receive an error
// and Steam will send you an authcode. Then you have to login again, this time with the authcode.
// Shortly after logging in, you'll receive a MachineAuthUpdateEvent with a hash which allows
// you to login without using an authcode in the future.
//
// If you don't use Steam Guard, username and password are enough

// Details used to auth user
type Details struct {
	Username     string
	Password     string
	AuthCode     string
	SharedSecret string
}

// Module used to auth user.
//
// The general gist of a Steam sign on is:
// ChannelEncryptRequest (server->client): The server is asking a user to negotiate
//     encryption on the specified universe.
// ChannelEncryptResponse (client->server): The client generates a key and encrypts
//     it with the universe public key. The specifics are in the implementation.
// ChannelEncryptResult (server->client): EResult of whether the negotiation was successful.
// ClientLogOn (client->server): Log on to service. Described in Steam and implementation.
// Multi (meta message, server->client): Server sent multiple messages, decompress and process, see implementation.
// ClientLogOnResponse (server->client): EResult of logon, as well as the authoritative steam id,
//     session ID (needs to be in every message header like steamid), and the heartbeat interval.
//     Every message after ClientLogOn needs the steamid and sessionID set in the header
// ClientGetAppOwnershipTicket (client->server): Client normally requests winui (appid 7) ownership ticket.
// ClientGetAppOwnershipTicketResponse (server->client): response to the above request,
//     returns the ownership ticket for the request appid , or an error EResult.
// ClientAccountInfo/EmailInfo/VACBanStatus (server->client): Account info
// ClientFriendsList (server->client): List of friends by steam id. The readable info is sent in PersonaStates.
// ClientLicenseList (server->client): List of licenses, which correspond to subs.
// ClientAuthList (client->server): Request a batch of game connect tokens from Steam
// ClientGameConnectTokens (server->client): Game connect tokens (see topicof tokens/tickets)
// ClientSessionToken (server->client): Used for something!
// ClientCMList (server->client): Authoritive list of CM servers,
//     for clients that use the bootstrap list of CM servers. CM servers are Steam servers.
// ClientHeartBeat (client->server): Heartbeat sent to server to signal connection is alive.
//     Sent every few seconds, as specified in the LogOnResponse message.
type Module struct {
	eventManager *steamprotocol.EventManager
	cl           *steamprotocol.Client
	gen          *TOTPGenerator
	details      Details
	sessionKey   []byte
	steamID      uint64
	sessionID    int32
}

// NewModule initialize new instance of auth Module.
func NewModule(
	cl *steamprotocol.Client,
	eventManager *steamprotocol.EventManager,
	gen *TOTPGenerator,
	details Details,
) *Module {
	return &Module{
		cl:           cl,
		eventManager: eventManager,
		gen:          gen,
		details:      details,
	}
}

// Subscribe used to start listen event and packets from eventManager.
func (m *Module) Subscribe() {
	m.eventManager.OnEvent(m.handleEvent)
	m.eventManager.OnPacket(m.handlePacket)
}

func (m *Module) handleEvent(e interface{}) error {
	switch e.(type) {
	case crypto.ChannelReadyEvent:
		return m.handleChannelEncryptedEvent()
	}

	return nil
}

func (m *Module) handlePacket(p *steamprotocol.Packet) error {
	switch p.Type {
	case steamprotocol.EMsg_ClientLogOnResponse:
		return m.handleLogOnResponse(p)
	case steamprotocol.EMsg_ClientLoggedOff:
		return m.handleLoggedOff(p)
	case steamprotocol.EMsg_ClientNewLoginKey:
		return m.handleNewLoginKey(p)
	case steamprotocol.EMsg_ClientUpdateMachineAuth:
		return m.handleUpdateMachineAuth(p)
	case steamprotocol.EMsg_ClientAccountInfo:
		// TODO: return m.handleAccountInfo(cl, p)
	}

	return nil
}

func (m *Module) handleChannelEncryptedEvent() error {
	if len(m.details.Username) == 0 {
		return errors.New("empty username")
	}

	if len(m.details.Password) == 0 {
		return errors.New("empty password")
	}

	responseHeader := messages.NewHeaderProto(steamprotocol.EMsg_ClientLogon)

	steamID := steamprotocol.NewIdAdv(
		0,
		1,
		int32(steamprotocol.EUniverse_Public),
		int32(steamprotocol.EAccountType_Individual),
	)

	// Save for sending with SuccessfullyAuthenticatedEvent event
	m.steamID = uint64(steamID)
	m.sessionID = 0

	responseHeader.Data.Steamid = proto.Uint64(m.steamID)
	responseHeader.Data.ClientSessionid = proto.Int32(m.sessionID)

	responseMsg := &protobuf.CMsgClientLogon{
		AccountName:     &m.details.Username,
		Password:        &m.details.Password,
		ClientLanguage:  proto.String("english"),
		ProtocolVersion: proto.Uint32(messages.ClientLogonCurrentProtocol),
		//ShaSentryfile:   []byte{}, // TODO: Get hash from storage
	}

	if m.details.AuthCode != "" {
		responseMsg.AuthCode = proto.String(m.details.AuthCode)
	}

	if m.details.SharedSecret != "" {
		code, err := m.gen.TwoFactorSynced(m.details.SharedSecret)
		if err != nil {
			return errors.Wrap(err, "failed to fetch two factor code")
		}

		responseMsg.TwoFactorCode = proto.String(code)
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

func (m *Module) handleLogOnResponse(p *steamprotocol.Packet) error {
	var (
		header *messages.HeaderProto = messages.NewHeaderProto(steamprotocol.EMsg_Invalid)
		msg    protobuf.CMsgClientLogonResponse
	)

	dataBuf := bytes.NewBuffer(p.Data)

	err := header.Deserialize(dataBuf)
	if err != nil {
		return errors.Wrap(err, "failed to deserialize logon response header")
	}

	err = proto.Unmarshal(dataBuf.Bytes(), &msg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal logon response msg")
	}

	result := steamprotocol.EResult(msg.GetEresult())

	if result == steamprotocol.EResult_OK {
		return m.eventManager.FireEvent(SuccessfullyAuthenticatedEvent{
			Heartbeat: time.Duration(msg.GetOutOfGameHeartbeatSeconds()) * time.Second,
			SteamID:   m.steamID,
			SessionID: m.sessionID,
		})
	}

	return m.eventManager.FireEvent(AuthenticationFailedEvent{
		Result: result,
	})
}

func (m *Module) handleLoggedOff(p *steamprotocol.Packet) error {
	var (
		header *messages.HeaderProto = messages.NewHeaderProto(steamprotocol.EMsg_Invalid)
		msg    protobuf.CMsgClientLoggedOff
	)

	dataBuf := bytes.NewBuffer(p.Data)

	err := header.Deserialize(dataBuf)
	if err != nil {
		return errors.Wrap(err, "failed to deserialize logged off response header")
	}

	err = proto.Unmarshal(dataBuf.Bytes(), &msg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal logged off response msg")
	}

	result := steamprotocol.EResult(msg.GetEresult())

	return m.eventManager.FireEvent(LoggedOffEvent{
		Result: result,
	})
}

func (m *Module) handleNewLoginKey(p *steamprotocol.Packet) error {
	var (
		header *messages.HeaderProto = messages.NewHeaderProto(steamprotocol.EMsg_Invalid)
		msg    protobuf.CMsgClientNewLoginKey
	)

	dataBuf := bytes.NewBuffer(p.Data)

	err := header.Deserialize(dataBuf)
	if err != nil {
		return errors.Wrap(err, "failed to deserialize new login key response header")
	}

	err = proto.Unmarshal(dataBuf.Bytes(), &msg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal new login key response msg")
	}

	uniqID := msg.GetUniqueId()
	key := msg.GetLoginKey()

	responseHeader := messages.NewHeaderProto(steamprotocol.EMsg_ClientNewLoginKeyAccepted)
	responseHeader.Data.Steamid = proto.Uint64(m.steamID)
	responseHeader.Data.ClientSessionid = proto.Int32(m.sessionID)

	responseMsg := &protobuf.CMsgClientNewLoginKeyAccepted{
		UniqueId: proto.Uint32(uniqID),
	}

	buf := new(bytes.Buffer)

	err = responseHeader.Serialize(buf)
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

	return m.eventManager.FireEvent(NewLoginKeyAcceptedEvent{
		UniqID: uniqID,
		Key:    key,
	})
}

func (m *Module) handleUpdateMachineAuth(p *steamprotocol.Packet) error {
	var (
		header *messages.HeaderProto = messages.NewHeaderProto(steamprotocol.EMsg_Invalid)
		msg    protobuf.CMsgClientUpdateMachineAuth
	)

	dataBuf := bytes.NewBuffer(p.Data)

	err := header.Deserialize(dataBuf)
	if err != nil {
		return errors.Wrap(err, "failed to deserialize new login key response header")
	}

	err = proto.Unmarshal(dataBuf.Bytes(), &msg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal new login key response msg")
	}

	hash := sha1.New()
	hash.Write(msg.Bytes)
	shaHash := hash.Sum(nil)

	responseHeader := messages.NewHeaderProto(steamprotocol.EMsg_ClientNewLoginKeyAccepted)
	responseHeader.Data.Steamid = proto.Uint64(m.steamID)
	responseHeader.Data.ClientSessionid = proto.Int32(m.sessionID)
	responseHeader.Data.JobidTarget = header.Data.JobidSource

	responseMsg := &protobuf.CMsgClientUpdateMachineAuthResponse{
		ShaFile: shaHash,
	}

	buf := new(bytes.Buffer)

	err = responseHeader.Serialize(buf)
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

	return m.eventManager.FireEvent(MachineAuthUpdateEvent{
		Hash: shaHash,
	})
}
