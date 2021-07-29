package connection

type IWebsocketConn interface {
	WriteMessage(messageType int, payload []byte) error
	ReadMessage() (messageType int, message []byte, err error)
}

// A multiplexed websocket connection that is capable of writing logs and command outputs to
// a single websocket connection
type MxedWebsocketConn struct {
	// The underlying original websocket connection
	conn     IWebsocketConn
	protocol *MxedWebsocketSubprotocol
}

func NewMxedWebsocketConn(conn IWebsocketConn) *MxedWebsocketConn {
	return &MxedWebsocketConn{conn: conn, protocol: NewMxedWebsocketSubprotocol()}
}

// Override the default write message to multiplex the message over a channelId. Messages
// sent over this channel will only reach the corresponding mxed websocket listening on
// this channel
// messageType indicates type of websocket message: text or binary.
// channelId & eventName can be used to target specific channels and actions.
func (mx *MxedWebsocketConn) WriteMessage(channelId string, eventName string, message []byte) {
	// Encode channelId and eventName and bytes with encoder
	output := mx.protocol.Encode(channelId, eventName, message)
	mx.conn.WriteMessage(2, output)
}

// Read message and return the appropriate channelId, eventName etc
func (mx *MxedWebsocketConn) ReadMessage() (DecodedMxWebsocketResponse, error) {
	_, message, err := mx.conn.ReadMessage()
	if err != nil {
		return DecodedMxWebsocketResponse{}, err
	}
	decoded, err := mx.protocol.Decode(message)
	if err != nil {
		return DecodedMxWebsocketResponse{}, err
	}
	return decoded, nil
}
