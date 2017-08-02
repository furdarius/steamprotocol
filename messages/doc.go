// Package messages has list of protocol messages structs
//
// Steam client messages normally consist of a general header,
// followed by a message structure, and end with a payload blob that is used in only a few messages.
// Messages are identified by integer constants known as an EMsg.

// There are three currently used message headers:
// MsgHdr is used before the client has an assigned session id or steamid. Used during the crypto handshake.
// ExtendedClientMsgHdr is used when the client has been assigned a session id or steamid.
// This is generally used after the crypto handshake.
// MsgHdrProtoBuf was recently introduced when Valve upgraded the Steam protocol to support protocol buffers.
// This is similar to the ExtendedClientMsgHdr, but is used when the client message is protobuf'd.

// Every client message has a structure which is normally the entire contents of the message,
// and is easily interpreted.
// The payload is only used in certain circumstances, such as the message data when sending MsgClientFriendMsg.
// There's a certain message of interest called a MsgMulti.
// This message is a wrapper message that encompasses multiple messages inside itself,
// often times compressed using PKZIP. The header of this message contains an int 'unzipped size',
// if this size is greater than 0, the message payload is compressed.
// After decompressing, the data is a stream of size and data pairs, with the data being a new client message.
//
// Read more: https://bitbucket.org/robingchan/steam/wiki/Networking/Protocol-level_messages.wiki
package messages
