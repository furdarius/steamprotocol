package messages

import (
	"encoding/binary"
	"io"

	"github.com/furdarius/steamprotocol"
)

type EncryptResponse struct {
	ProtocolVersion uint32
	KeySize         uint32
}

func NewEncryptResponse(protocolVersion uint32, keySize uint32) *EncryptResponse {
	return &EncryptResponse{
		ProtocolVersion: protocolVersion,
		KeySize:         keySize,
	}
}

func (m *EncryptResponse) Type() steamprotocol.EMsg {
	return steamprotocol.EMsg_ChannelEncryptResponse
}

func (m *EncryptResponse) Serialize(w io.Writer) error {
	err := binary.Write(w, binary.LittleEndian, m.ProtocolVersion)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, m.KeySize)
	if err != nil {
		return err
	}

	return nil
}

func (m *EncryptResponse) Deserialize(r io.Reader) error {
	err := binary.Read(r, binary.LittleEndian, &m.ProtocolVersion)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.LittleEndian, &m.KeySize)
	if err != nil {
		return err
	}

	return nil
}
