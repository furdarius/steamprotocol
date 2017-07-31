package auth

import (
	"time"

	"github.com/furdarius/steamprotocol"
)

type SuccessfullyAuthenticatedEvent struct {
	Heartbeat time.Duration
	SteamID   uint64
	SessionID int32
}

type AuthenticationFailedEvent struct {
	Result steamprotocol.EResult
}

type LoggedOffEvent struct {
	Result steamprotocol.EResult
}

type NewLoginKeyAcceptedEvent struct {
	UniqID uint32
	Key    string
}

type MachineAuthUpdateEvent struct {
	Hash []byte
}
