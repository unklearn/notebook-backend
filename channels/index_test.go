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
