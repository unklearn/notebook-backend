package connection

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMxedWebsocketEncoder(t *testing.T) {
	e := NewMxedWebsocketSubprotocol()
	res := e.Encode("chan", "event-name", []byte("oyo koyo muntoyo"))
	target := []byte{4, 0, 0, 0, 10, 0, 0, 0, 99, 104, 97, 110, 101, 118, 101, 110, 116, 45, 110, 97, 109, 101, 111, 121, 111, 32, 107, 111, 121, 111, 32, 109, 117, 110, 116, 111, 121, 111}
	if !bytes.Equal(res, target) {
		t.Errorf("Expected encode response to match target array, instead got %v", res)
	}
}

func TestMxedWebsocketDecoderLargeSizeLimit(t *testing.T) {
	e := NewMxedWebsocketSubprotocol()
	_, err := e.Decode([]byte{0, 0, 0, 4, 10, 0, 0, 0, 99, 104, 97, 110, 101, 118, 101, 110, 116, 45, 110, 97, 109, 101, 111, 121, 111, 32, 107, 111, 121, 111, 32, 109, 117, 110, 116, 111, 121, 111})
	assert.NotEqual(t, err, nil)
	_, err = e.Decode([]byte{4, 0, 0, 0, 0, 0, 0, 10, 99, 104, 97, 110, 101, 118, 101, 110, 116, 45, 110, 97, 109, 101, 111, 121, 111, 32, 107, 111, 121, 111, 32, 109, 117, 110, 116, 111, 121, 111})
	assert.NotEqual(t, err, nil)
}

func TestMxedWebsocketDecoder(t *testing.T) {
	e := NewMxedWebsocketSubprotocol()
	res, _ := e.Decode([]byte{4, 0, 0, 0, 10, 0, 0, 0, 99, 104, 97, 110, 101, 118, 101, 110, 116, 45, 110, 97, 109, 101, 111, 121, 111, 32, 107, 111, 121, 111, 32, 109, 117, 110, 116, 111, 121, 111})
	if res.ChannelId != "chan" {
		t.Errorf("Expected channelId to be chan instead got %s", res.ChannelId)
	}
	if res.EventName != "event-name" {
		t.Errorf("Expected event name to be event-name instead got %s", res.EventName)
	}
	if !bytes.Equal([]byte{111, 121, 111, 32, 107, 111, 121, 111, 32, 109, 117, 110, 116, 111, 121, 111}, res.Payload) {
		t.Errorf("Expected decoded payload to match %v", res.Payload)
	}
}
