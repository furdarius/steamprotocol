package messages

import (
	"encoding/binary"
	"io"

	"github.com/furdarius/steamprotocol"
)

type Header struct {
	Type        steamprotocol.EMsg
	TargetJobID uint64
	SourceJobID uint64
}

func NewHeader(msgType steamprotocol.EMsg, targetJobID uint64, sourceJobID uint64) *Header {
	return &Header{
		Type:        msgType,
		TargetJobID: targetJobID,
		SourceJobID: sourceJobID,
	}
}

func (m *Header) Serialize(w io.Writer) error {
	err := binary.Write(w, binary.LittleEndian, m.Type)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, m.TargetJobID)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, m.SourceJobID)
	if err != nil {
		return err
	}

	return nil
}

func (m *Header) Deserialize(r io.Reader) error {
	var t int32
	err := binary.Read(r, binary.LittleEndian, &t)
	if err != nil {
		return err
	}

	m.Type = steamprotocol.EMsg(t)

	err = binary.Read(r, binary.LittleEndian, &m.TargetJobID)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.LittleEndian, &m.SourceJobID)
	if err != nil {
		return err
	}

	return nil
}
