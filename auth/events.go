package auth

import (
	"time"

	"github.com/furdarius/steamprotocol"
)

// SuccessfullyAuthenticatedEvent is fired when successful CMsgClientLogonResponse is received.
type SuccessfullyAuthenticatedEvent struct {
	Heartbeat time.Duration
	SteamID   uint64
	SessionID int32
}

// AuthenticationFailedEvent is fired when failed CMsgClientLogonResponse is received.
type AuthenticationFailedEvent struct {
	Result steamprotocol.EResult
}

// LoggedOffEvent is fired when steam log off authenticated client.
type LoggedOffEvent struct {
	Result steamprotocol.EResult
}

// NewLoginKeyAcceptedEvent is fired when CMsgClientNewLoginKey is received.
type NewLoginKeyAcceptedEvent struct {
	UniqID uint32
	Key    string
}

// MachineAuthUpdateEvent is fired when CMsgClientUpdateMachineAuth is received.
type MachineAuthUpdateEvent struct {
	Hash []byte
}
