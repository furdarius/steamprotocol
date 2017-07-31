package messages

import (
	"encoding/binary"
	"io"

	"github.com/furdarius/steamprotocol"
)

type EncryptResult struct {
	Result steamprotocol.EResult
}

func NewEncryptResult() *EncryptResult {
	return &EncryptResult{
		Result: steamprotocol.EResult_Invalid,
	}
}

func (m *EncryptResult) Type() steamprotocol.EMsg {
	return steamprotocol.EMsg_ChannelEncryptResult
}

func (m *EncryptResult) Serialize(w io.Writer) error {
	err := binary.Write(w, binary.LittleEndian, m.Result)
	if err != nil {
		return err
	}

	return nil
}

func (m *EncryptResult) Deserialize(r io.Reader) error {
	var tmp int32
	err := binary.Read(r, binary.LittleEndian, &tmp)
	if err != nil {
		return err
	}

	m.Result = steamprotocol.EResult(tmp)

	return nil
}
