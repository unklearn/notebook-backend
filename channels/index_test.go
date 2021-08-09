package channels

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unklearn/notebook-backend/commands"
)

func TestGetId(t *testing.T) {
	rc := NewRootChannel("foo")
	assert.Equal(t, rc.GetId(), "foo")
}

func TestHandleMessageContainerStart(t *testing.T) {
	rc := NewRootChannel("chan")
	its, err := rc.HandleMessage(string(ContainerStartEventName), []byte(`{"container_id": "foo"}`))
	assert.Equal(t, err, nil)
	assert.Equal(t, len(its), 1)
	assert.IsType(t, &commands.ContainerCreateCommandIntent{}, its[0])
}

func TestHandleMessageUnknownType(t *testing.T) {
	rc := NewRootChannel("chan")
	_, err := rc.HandleMessage("many", []byte(`{"container_id": "foo"}`))
	assert.NotEqual(t, err, nil, "Should only handle known messages")
}

func TestContainerChannelExecuteCommand(t *testing.T) {
	cc := NewContainerChannel("foo")
	payload := []byte(`{"command": ["bash"]}`)
	intents, e := cc.HandleMessage(string(ContainerExecuteCommandEventName), payload)
	assert.Equal(t, e, nil)
	c, _ := commands.NewContainerExecuteCommandIntent("foo", payload)
	assert.Equal(t, intents[0], c)
}
