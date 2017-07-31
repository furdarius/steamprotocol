package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash/crc32"

	"github.com/pkg/errors"
	"github.com/furdarius/steamprotocol"
	"github.com/furdarius/steamprotocol/messages"
)

type Module struct {
	eventManager *steamprotocol.EventManager
	cl           *steamprotocol.Client
	sessionKey   []byte
}

func NewModule(cl *steamprotocol.Client, eventManager *steamprotocol.EventManager) *Module {
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
	case steamprotocol.EMsg_ChannelEncryptRequest:
		return m.handleChannelEncryptRequest(p)
	case steamprotocol.EMsg_ChannelEncryptResult:
		return m.handleChannelEncryptResult(p)
	default:
		return nil
	}
}

func (m *Module) handleChannelEncryptRequest(p *steamprotocol.Packet) error {
	var (
		header messages.Header
		msg    messages.EncryptRequest
	)

	r := bytes.NewReader(p.Data)

	err := header.Deserialize(r)
	if err != nil {
		return errors.Wrap(err, "failed to decode encrypt request header")
	}

	err = msg.Deserialize(r)
	if err != nil {
		return errors.Wrap(err, "failed to decode encrypt request")
	}

	if msg.Universe != steamprotocol.EUniverse_Public {
		return errors.New("invalid universe")
	}

	if msg.ProtocolVersion != messages.EncryptRequestDefaultProtocol {
		return fmt.Errorf("invalid protocol version %d", msg.ProtocolVersion)
	}

	pub, err := getPublicKey(steamprotocol.EUniverse_Public)
	if err != nil {
		return errors.Wrapf(err, "failed to get public key for universe %v", msg.Universe)
	}

	m.sessionKey = make([]byte, 32)
	rand.Read(m.sessionKey)

	encryptedKey, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, m.sessionKey, nil)
	if err != nil {
		return errors.Wrap(err, "failed to encrypt session key")
	}

	responseMsg := messages.NewEncryptResponse(messages.EncryptRequestDefaultProtocol, 128)
	responseHeader := messages.NewHeader(responseMsg.Type(), header.SourceJobID, header.TargetJobID)

	buf := new(bytes.Buffer)

	err = responseHeader.Serialize(buf)
	if err != nil {
		return errors.Wrap(err, "failed to serialize header")
	}

	err = responseMsg.Serialize(buf)
	if err != nil {
		return errors.Wrap(err, "failed to serialize response msg")
	}

	n, err := buf.Write(encryptedKey)
	if err != nil || len(encryptedKey) != n {
		return errors.Wrap(err, "failed to write encrypted key to result buffer")
	}

	keyHash := crc32.ChecksumIEEE(encryptedKey)

	err = binary.Write(buf, binary.LittleEndian, keyHash)
	if err != nil {
		return errors.Wrap(err, "failed to write key hash to result buffer")
	}

	err = binary.Write(buf, binary.LittleEndian, uint32(0))
	if err != nil {
		return errors.Wrap(err, "failed to write finish bytes to result buffer")
	}

	err = m.cl.Write(buf.Bytes())
	if err != nil {
		return errors.Wrap(err, "failed to write data")
	}

	return nil
}

func (m *Module) handleChannelEncryptResult(p *steamprotocol.Packet) error {
	var (
		header messages.Header
		msg    messages.EncryptResult
	)

	r := bytes.NewReader(p.Data)

	err := header.Deserialize(r)
	if err != nil {
		return errors.Wrap(err, "failed to decode encrypt request header")
	}

	err = msg.Deserialize(r)
	if err != nil {
		return errors.Wrap(err, "failed to decode encrypt request")
	}

	if msg.Result != steamprotocol.EResult_OK {
		return errors.New("encryption failed")
	}

	block, err := aes.NewCipher(m.sessionKey)
	if err != nil {
		return errors.Wrap(err, "failed to create cipher from session key")
	}

	encryptor := NewAes(block)
	m.cl.SetEncryptor(encryptor)

	return m.eventManager.FireEvent(ChannelReadyEvent{})
}
