package steamprotocol

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"

	"time"

	"fmt"

	"github.com/pkg/errors"
)

const (
	// Magic contains in all TCP packets, and must be read after packet length bytes.
	Magic     uint32 = 0x31305456 // "VT01"

	// ProtoMask used to check is it protobuf message.
	// Message is proto if expression "rawMsg & ProtoMask > 0" is true.
	ProtoMask uint32 = 0x80000000

	// EMsgMask used to get EMsg by rawMsg: EMsg(rawMsg & EMsgMask).
	// Messages are identified by integer constants known as an EMsg.
	EMsgMask         = ^ProtoMask
)

// Encryptor used to encrypt data on write and decrypt on read.
// It's required after ChannelEncryptResult message gotten.
type Encryptor interface {
	Encrypt(src []byte) ([]byte, error)
	Decrypt(src []byte) []byte
}

// Packet is container of message data.
// It's used to broadcast message to packet handlers.
type Packet struct {
	Type EMsg
	Data []byte
}

// Client implements communication with Steam CM servers.
type Client struct {
	conn         net.Conn
	eventManager *EventManager
	crypto       Encryptor
}

// NewClient initialize new instance of Client.
func NewClient(conn net.Conn, eventManager *EventManager) *Client {
	return &Client{
		conn:         conn,
		eventManager: eventManager,
	}
}

// Listen start to read connection with Steam server.
// It uses endless cycle for reading.
func (c *Client) Listen() error {
	var (
		packetLen   uint32
		packetMagic uint32
		packet      *Packet
	)

	for {
		err := binary.Read(c.conn, binary.LittleEndian, &packetLen)
		if err != nil {
			if err == io.EOF {
				time.Sleep(time.Millisecond * 100)

				continue
			}

			return errors.Wrap(err, "failed to read packet length")
		}

		err = binary.Read(c.conn, binary.LittleEndian, &packetMagic)
		if err != nil {
			return errors.Wrap(err, "failed to read packet magic")
		}

		if packetMagic != Magic {
			return errors.New("invalid connection magic")
		}

		// Used to accumulate packet data and then broadcast it to handlers via EventManager
		buf := make([]byte, packetLen)

		_, err = io.ReadFull(c.conn, buf)
		if err != nil {
			return errors.Wrap(err, "failed to read packet data to buffer")
		}

		if c.crypto != nil {
			buf = c.crypto.Decrypt(buf)
		}

		r := bytes.NewReader(buf)

		var rawMsg uint32
		err = binary.Read(r, binary.LittleEndian, &rawMsg)
		if err != nil {
			return errors.Wrap(err, "failed to read raw msg")
		}

		eMsg := EMsg(rawMsg & EMsgMask)

		// Костыль, для того, что-бы в handlePacket не перечитывать значения заголовка заново
		//startFrom, err := r.Seek(0, io.SeekCurrent)
		//if err != nil {
		//	return errors.Wrap(err, "failed to seek reader position")
		//}

		packet = &Packet{
			Type: eMsg,
			Data: buf,
		}

		err = c.handlePacket(packet)
		if err != nil {
			return errors.Wrap(err, "failed to handle packet")
		}
	}
}

// Write is used to write byte array to Steam connection.
func (c *Client) Write(data []byte) (err error) {
	if c.conn == nil {
		return errors.New("connection is not defined")
	}

	if c.crypto != nil {
		data, err = c.crypto.Encrypt(data)
		if err != nil {
			return errors.Wrap(err, "failed to encrypt data")
		}
	}

	dataLen := uint32(len(data))

	err = binary.Write(c.conn, binary.LittleEndian, dataLen)
	if err != nil {
		return errors.Wrap(err, "failed to write data len")
	}

	err = binary.Write(c.conn, binary.LittleEndian, Magic)
	if err != nil {
		return errors.Wrap(err, "failed to write magic")
	}

	n, err := c.conn.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to write data")
	}

	if uint32(n) != dataLen {
		return errors.Wrap(err, "data wasn't fully sent")
	}

	return nil
}

func (c *Client) handlePacket(p *Packet) error {
	fmt.Println("CLIENT GOT MESSAGE: ", p.Type)

	return c.eventManager.FirePacket(p)
}

// SetEncryptor change data Encryptor in Client instance.
func (c *Client) SetEncryptor(enc Encryptor) {
	c.crypto = enc
}
