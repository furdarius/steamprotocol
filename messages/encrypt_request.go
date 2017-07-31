package messages

import (
	"encoding/binary"
	"io"

	"github.com/furdarius/steamprotocol"
)

const (
	EncryptRequestDefaultProtocol uint32 = 1
)

type EncryptRequest struct {
	ProtocolVersion uint32
	Universe        steamprotocol.EUniverse
}

func (m *EncryptRequest) Type() steamprotocol.EMsg {
	return steamprotocol.EMsg_ChannelEncryptRequest
}

func (m *EncryptRequest) Serialize(w io.Writer) error {
	err := binary.Write(w, binary.LittleEndian, m.ProtocolVersion)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, m.Universe)
	if err != nil {
		return err
	}

	return nil
}

func (m *EncryptRequest) Deserialize(r io.Reader) error {
	err := binary.Read(r, binary.LittleEndian, &m.ProtocolVersion)
	if err != nil {
		return err
	}

	var universe int32
	err = binary.Read(r, binary.LittleEndian, &universe)
	if err != nil {
		return err
	}

	m.Universe = steamprotocol.EUniverse(universe)

	return nil
}
