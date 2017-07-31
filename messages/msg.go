package messages

//import (
//	"io"
//
//	"github.com/furdarius/steamprotocol"
//)

//
//import (
//	"io"
//
//	"github.com/golang/protobuf/proto"
//)
//

//type Serializer interface {
//	Serialize(io.Writer) error
//}
//
//type Deserializer interface {
//	Deserialize(io.Reader) error
//}
//
//type MessageBody interface {
//	Serializer
//	Deserializer
//	Type() steamprotocol.EMsg
//}

//
//// Interface for all messages, typically outgoing. They can also be created by
//// using the Read* methods in a PacketMsg.
//type IMsg interface {
//	Serializer
//	IsProto() bool
//	GetMsgType() EMsg
//	GetTargetJobId() uint64
//	SetTargetJobId(uint64)
//	GetSourceJobId() uint64
//	SetSourceJobId(uint64)
//}
//
//// Interface for client messages, i.e. messages that are sent after logging in.
//// ClientMsgProtobuf and ClientMsg implement this.
//type IClientMsg interface {
//	IMsg
//	GetSessionId() int32
//	SetSessionId(int32)
//	GetSteamId() uint64
//	SetSteamId(uint64)
//}
//
//// Represents a protobuf backed client message with session data.
//type ClientMsgProtobuf struct {
//	Header *MsgHdrProtoBuf
//	Body   proto.Message
//}
//
//func NewClientMsgProtobuf(eMsg EMsg, body proto.Message) *ClientMsgProtobuf {
//	hdr := NewMsgHdrProtoBuf()
//	hdr.Msg = eMsg
//	return &ClientMsgProtobuf{
//		Header: hdr,
//		Body:   body,
//	}
//}
//
//func (c *ClientMsgProtobuf) IsProto() bool {
//	return true
//}
//
//func (c *ClientMsgProtobuf) GetMsgType() EMsg {
//	return EMsg(uint32(c.Header.Msg) & EMsgMask)
//}
//
//func (c *ClientMsgProtobuf) GetSessionId() int32 {
//	return c.Header.Proto.GetClientSessionid()
//}
//
//func (c *ClientMsgProtobuf) SetSessionId(session int32) {
//	c.Header.Proto.ClientSessionid = &session
//}
//
//func (c *ClientMsgProtobuf) GetSteamId() uint64 {
//	return c.Header.Proto.GetSteamid()
//}
//
//func (c *ClientMsgProtobuf) SetSteamId(s uint64) {
//	c.Header.Proto.Steamid = proto.Uint64(s)
//}
//
//func (c *ClientMsgProtobuf) GetTargetJobId() uint64 {
//	return c.Header.Proto.GetJobidTarget()
//}
//
//func (c *ClientMsgProtobuf) SetTargetJobId(job uint64) {
//	c.Header.Proto.JobidTarget = proto.Uint64(job)
//}
//
//func (c *ClientMsgProtobuf) GetSourceJobId() uint64 {
//	return c.Header.Proto.GetJobidSource()
//}
//
//func (c *ClientMsgProtobuf) SetSourceJobId(job uint64) {
//	c.Header.Proto.JobidSource = proto.Uint64(job)
//}
//
//func (c *ClientMsgProtobuf) Serialize(w io.Writer) error {
//	err := c.Header.Serialize(w)
//	if err != nil {
//		return err
//	}
//
//	body, err := proto.Marshal(c.Body)
//	if err != nil {
//		return err
//	}
//
//	_, err = w.Write(body)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

//type Msg struct {
//	Header  *Header
//	Body    MessageBody
//	Payload []byte
//}
//
//func NewMsg(body MessageBody, payload []byte) *Msg {
//	hdr := NewHeader(steamprotocol.EMsg_Invalid, ^uint64(0), ^uint64(0))
//	hdr.Type = body.Type()
//	return &Msg{
//		Header:  hdr,
//		Body:    body,
//		Payload: payload,
//	}
//}
//
//func (m *Msg) GetMsgType() steamprotocol.EMsg {
//	return m.Header.Type
//}
//
//func (m *Msg) IsProto() bool {
//	return false
//}
//
//func (m *Msg) GetTargetJobId() uint64 {
//	return m.Header.TargetJobID
//}
//
//func (m *Msg) SetTargetJobId(job uint64) {
//	m.Header.TargetJobID = job
//}
//
//func (m *Msg) GetSourceJobId() uint64 {
//	return m.Header.SourceJobID
//}
//
//func (m *Msg) SetSourceJobId(job uint64) {
//	m.Header.SourceJobID = job
//}
//
//func (m *Msg) Serialize(w io.Writer) error {
//	err := m.Header.Serialize(w)
//	if err != nil {
//		return err
//	}
//
//	err = m.Body.Serialize(w)
//	if err != nil {
//		return err
//	}
//
//	_, err = w.Write(m.Payload)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
