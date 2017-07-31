package messages

import (
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
	"github.com/furdarius/steamprotocol"
	"github.com/furdarius/steamprotocol/protobuf"
)

type HeaderProto struct {
	Type         steamprotocol.EMsg
	HeaderLength int32
	Data         *protobuf.CMsgProtoBufHeader
}

func NewHeaderProto(msgType steamprotocol.EMsg) *HeaderProto {
	return &HeaderProto{
		Type: msgType,
		Data: &protobuf.CMsgProtoBufHeader{},
	}
}

func (m *HeaderProto) Serialize(w io.Writer) error {
	buf, err := proto.Marshal(m.Data)
	if err != nil {
		return err
	}

	msgType := steamprotocol.EMsg(uint32(m.Type) | steamprotocol.ProtoMask)

	err = binary.Write(w, binary.LittleEndian, msgType)
	if err != nil {
		return err
	}

	headerLen := int32(len(buf))

	err = binary.Write(w, binary.LittleEndian, headerLen)
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (m *HeaderProto) Deserialize(r io.Reader) error {
	var t int32
	err := binary.Read(r, binary.LittleEndian, &t)
	if err != nil {
		return err
	}

	m.Type = steamprotocol.EMsg(uint32(t) & steamprotocol.EMsgMask)

	err = binary.Read(r, binary.LittleEndian, &m.HeaderLength)
	if err != nil {
		return err
	}

	buf := make([]byte, m.HeaderLength)

	_, err = io.ReadFull(r, buf)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(buf, m.Data)
	if err != nil {
		return err
	}

	return nil
}
