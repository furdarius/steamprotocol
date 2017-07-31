package multi

import (
	"bytes"

	"compress/gzip"
	"encoding/binary"
	"io"
	"io/ioutil"

	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/furdarius/steamprotocol"
	"github.com/furdarius/steamprotocol/messages"
	"github.com/furdarius/steamprotocol/protobuf"
)

type Module struct {
	eventManager *steamprotocol.EventManager
	cl           *steamprotocol.Client
}

func NewModule(
	cl *steamprotocol.Client,
	eventManager *steamprotocol.EventManager,
) *Module {
	return &Module{
		cl:           cl,
		eventManager: eventManager,
	}
}

func (m *Module) Subscribe() {
	m.eventManager.OnPacket(m.handlePacket)
}

func (m *Module) handlePacket(p *steamprotocol.Packet) error {
	switch p.Type {
	case steamprotocol.EMsg_Multi:
		return m.handleMulti(p)
	}

	return nil
}

func (m *Module) handleMulti(p *steamprotocol.Packet) error {
	var (
		header *messages.HeaderProto = messages.NewHeaderProto(steamprotocol.EMsg_Invalid)
		msg    protobuf.CMsgMulti
	)

	dataBuf := bytes.NewBuffer(p.Data)

	err := header.Deserialize(dataBuf)
	if err != nil {
		return errors.Wrap(err, "failed to deserialize multi header")
	}

	err = proto.Unmarshal(dataBuf.Bytes(), &msg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal multi msg")
	}

	payload := msg.GetMessageBody()

	if msg.GetSizeUnzipped() > 0 {
		r, err := gzip.NewReader(bytes.NewReader(payload))
		if err != nil {
			return errors.Wrap(err, "failed to decompress payload")
		}

		payload, err = ioutil.ReadAll(r)
		if err != nil {
			return errors.Wrap(err, "failed to read all decompressed payload")
		}
	}

	pr := bytes.NewReader(payload)

	// Cycle body same as steamprotocol.Client.Read method,
	// so, it can be better to move it to function,
	// otherwise we had code duplication
	// TODO: Refactor code for packet preparing
	for pr.Len() > 0 {
		var packetLen uint32
		binary.Read(pr, binary.LittleEndian, &packetLen)

		buf := make([]byte, packetLen)

		_, err = io.ReadFull(pr, buf)
		if err != nil {
			return errors.Wrap(err, "failed to read packet data to buffer")
		}

		r := bytes.NewReader(buf)

		var rawMsg uint32
		err = binary.Read(r, binary.LittleEndian, &rawMsg)
		if err != nil {
			return errors.Wrap(err, "failed to read raw msg")
		}

		eMsg := steamprotocol.EMsg(rawMsg & steamprotocol.EMsgMask)

		fmt.Println(eMsg)

		packet := &steamprotocol.Packet{
			Type: eMsg,
			Data: buf,
		}

		err = m.eventManager.FirePacket(packet)
		if err != nil {
			return err
		}
	}

	return nil
}
