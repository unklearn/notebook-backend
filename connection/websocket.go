package connection

import (
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// A multiplexed websocket connection that is capable of writing logs and command outputs to
// a single websocket connection
type MxedWebsocketConn struct {
	// The underlying original websocket connection
	Conn *websocket.Conn
	// The Delimiter to use for writing the message/reading message
	Delimiter string
	// Map of registered channels
	ChannelMap map[string]Channel
}

// Override the default write message to multiplex the message over a channelId. Messages
// sent over this channel will only reach the corresponding mxed websocket listening on
// this channel
func (mx *MxedWebsocketConn) WriteMessage(messageType int, channelId string, message []byte) {
	// Augment message by adding channelId
	message = append([]byte(channelId+mx.Delimiter), message...)
	mx.Conn.WriteMessage(messageType, message)
}

// parse channelId and return channel id and actual message
func (mx *MxedWebsocketConn) parseChannelId(message []byte) (string, []byte) {
	if len(mx.Delimiter) == 0 {
		return "", message
	}
	currByte := mx.Delimiter[0]
	len_delim := len(mx.Delimiter)
	for i, b := range message {
		if b == currByte {
			if mx.Delimiter == string(message[i:len_delim+i]) {
				return string(message[0:i]), message[(i + len_delim):]
			}
		}
	}
	return "", message
}

// Read message and dispatch to appropriate channels
func (mx *MxedWebsocketConn) ReadMessage() error {
	_, message, err := mx.Conn.ReadMessage()
	if err != nil {
		return err
	}
	channelId, msg := mx.parseChannelId(message)
	if mx.ChannelMap == nil {
		return errors.New("no channels have been registered yet")
	}
	// Check if channel exists
	ch, ok := mx.ChannelMap[channelId]
	if !ok {
		// Log warning
		log.Printf("Missing channelId %s", channelId)
		keys := make([]string, len(mx.ChannelMap))

		i := 0
		for k := range mx.ChannelMap {
			keys[i] = k
			i++
		}
		log.Printf("keys %v", keys)
		return nil
	}
	_, err = ch.Write(msg)
	return err
}

// Register a new channel if another does not exist already
func (mx *MxedWebsocketConn) RegisterChannel(channelId string, channel Channel) error {

	// Check if channel exists
	_, ok := mx.ChannelMap[channelId]
	if !ok {
		mx.ChannelMap[channelId] = channel
		keys := make([]string, len(mx.ChannelMap))
		i := 0
		for k := range mx.ChannelMap {
			keys[i] = k
			i++
		}
		log.Printf("keys %v", keys)
		return nil
	} else {
		return fmt.Errorf("another channel exists for channelId %s", channelId)
	}
}
