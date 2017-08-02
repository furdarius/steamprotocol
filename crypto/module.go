// Package crypto used to encrypt communication channel.
//
// After establishing a connection to a CM server, the server and client go through a handshake process
// that establishes an encrypted connection.
// Client messages are encrypted using AES with a session key that is generated
// by the client during the handshake.
// There exists evidence that a connection can be unencrypted,
// because of the export restriction of strong cryptography from the US,
// but it has not been observed.
//
// Steps:
// 1. Server requests the client to encrypt traffic within the specified universe (normally Public)
// 2. Client generates a 256bit session key.
// 3. This key is encrypted by a 1024bit public RSA key for the specific universe.
// 4. The encrypted key is sent to the server, along with a 32bit crc of the encrypted key.
// 5. The server replies with an unencrypted success/failure message.
// 6. All traffic from here is AES encrypted with the session key.
//
// Symmetric crypto
// * All messages after the handshake are AES encrypted.
// * A random 16 byte IV is generated for every message.
// * This IV is AES encrypted in ECB mode using the session key generated during the handshake.
// * Message data is encrypted with AES using the generated (not encrypted) IV and session key in CBC mode.
// * The encrypted IV and encrypted message data are concatenated together and sent off.
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

	"github.com/furdarius/steamprotocol"
	"github.com/furdarius/steamprotocol/messages"
	"github.com/pkg/errors"
)

// Module used to encrypt communication channel.
type Module struct {
	eventManager *steamprotocol.EventManager
	cl           *steamprotocol.Client
	sessionKey   []byte
}

// NewModule initialize new instance of crypto Module.
func NewModule(cl *steamprotocol.Client, eventManager *steamprotocol.EventManager) *Module {
	return &Module{
		cl:           cl,
		eventManager: eventManager,
	}
}

// Subscribe used to start listen event and packets from eventManager.
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
