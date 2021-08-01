package connection

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// Multiplexed websocket encoder writes mulitplexed messages onto underlying socket for a mxed websocket
// This creates a websocket subprotocol for client to use

// The protocol is faily simple
//
//  |----4 bytes----||----4 bytes----|-------channelId---------|--------eventName--------|payload
//
//  The first 4 bytes store length of channelId as uint32, and second 4 bytes store length of event name
//  The parsed lengths are used to read in channelId and eventName and payload
type MxedWebsocketSubprotocol struct {
	subProtocolName string
}

// Holder for decoded messages
type DecodedMxWebsocketResponse struct {
	ChannelId string
	EventName string
	Payload   []byte
}

// Encode writes to a given channel and eventName by encoding the lengths,
func (e MxedWebsocketSubprotocol) Encode(channelId string, eventName string, message []byte) []byte {
	dst := make([]byte, 0)
	buf := bytes.NewBuffer(dst)
	channelIdLen := make([]byte, 4)
	eventNameLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(channelIdLen, uint32(len(channelId)))
	binary.LittleEndian.PutUint32(eventNameLen, uint32(len(eventName)))
	buf.Write(channelIdLen)
	buf.Write([]byte(eventNameLen))
	buf.Write([]byte(channelId))
	buf.Write([]byte(eventName))
	buf.Write(message)
	return buf.Bytes()
}

func (e MxedWebsocketSubprotocol) Decode(message []byte) (DecodedMxWebsocketResponse, error) {
	r := DecodedMxWebsocketResponse{}
	// Read first 4 bytes, then next 4
	channelIdSize := int(binary.LittleEndian.Uint32(message[0:4]))
	eventNameSize := int(binary.LittleEndian.Uint32(message[4:8]))

	// Truncate sizes to a max so that we don't read out of bounds
	if channelIdSize > 256 || eventNameSize > 256 {
		return r, errors.New("channel id or event name must be less than 256 characters")
	}

	// Parse channelId and eventName
	r.ChannelId = string(message[8 : channelIdSize+8])
	r.EventName = string(message[channelIdSize+8 : channelIdSize+8+eventNameSize])
	r.Payload = message[channelIdSize+eventNameSize+8:]

	if r.ChannelId == "" {
		return DecodedMxWebsocketResponse{}, fmt.Errorf("ECODE::enc-dec-bad-channel-id::Missing channel Id")
	}
	if r.EventName == "" {
		return DecodedMxWebsocketResponse{}, fmt.Errorf("ECODE::enc-dec-bad-event-name::Missing event name")
	}
	return r, nil
}

// Return the subprotocol name for the encoder
func (e MxedWebsocketSubprotocol) GetSubprotocol() string {
	return e.subProtocolName
}

func NewMxedWebsocketSubprotocol() *MxedWebsocketSubprotocol {
	return &MxedWebsocketSubprotocol{subProtocolName: "unk"}
}
