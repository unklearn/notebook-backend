package connection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeWebsocketConn struct {
	intake    []byte
	intBuffer []byte
}

func (f *fakeWebsocketConn) WriteMessage(messageType int, payload []byte) error {
	// Write to internal buffer
	f.intBuffer = payload
	return nil
}

func (f fakeWebsocketConn) ReadMessage() (messageType int, payload []byte, err error) {
	return 0, f.intake, nil
}

func TestMxedWebsocketWriteMessage(t *testing.T) {
	f := &fakeWebsocketConn{intake: []byte("wohoo")}
	sub := NewMxedWebsocketSubprotocol()
	mx := NewMxedWebsocketConn(f)
	mx.WriteMessage("chan", "some-event", []byte("woohoo"))
	assert.Equal(t, f.intBuffer, sub.Encode("chan", "some-event", []byte("woohoo")))
}

func TestMxedWebsocketRead(t *testing.T) {
	sub := NewMxedWebsocketSubprotocol()
	f := &fakeWebsocketConn{intake: sub.Encode("chan", "some-event", []byte("woohoo"))}
	mx := NewMxedWebsocketConn(f)
	d, err := mx.ReadMessage()
	assert.Equal(t, err, nil)
	assert.Equal(t, d.ChannelId, "chan")
	assert.Equal(t, d.EventName, "some-event")
	assert.Equal(t, d.Payload, []byte("woohoo"))
}
