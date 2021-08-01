package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerCommandCreateIntent(t *testing.T) {
	b := []byte(`{"container_id": "foo", "name": "name", "image": "python", "tag": "3.6", "command": ["sh"], "network_options": {"ports": ["8000"]}, "env": ["k:v"]}`)
	c, e := NewContainerCreateCommandIntent("chan", b)
	if e != nil {
		t.Errorf("Expected error to be nil, except got error %s", e.Error())
	} else {
		assert.Equal(t, c.ChannelId, "chan")
		assert.Equal(t, c.ContainerId, "foo")
		assert.Equal(t, len(c.NetworkOptions.Ports), 1)
		assert.Equal(t, c.NetworkOptions.Ports[0], "8000")
		assert.Equal(t, c.EnvVars[0], "k:v")
		assert.Equal(t, c.Image, "python")
		assert.Equal(t, c.ImageTag, "3.6")
		assert.Equal(t, c.Command, []string{"sh"})
	}

	// Try with wrong type
	_, e = NewContainerCreateCommandIntent("chan", []byte(`{"container_id": 123, "name": "name", "image": "python", "tag": "3.6", "command": ["sh"], "network_options": {"ports": ["8000"]}, "env": ["k:v"]}`))
	assert.Equal(t, e.Error(), "invalid input supplied for creating container")

	// Try with missing value
	_, e = NewContainerCreateCommandIntent("chan", []byte(`{"image": "python", "tag": "3.6", "command": ["sh"], "network_options": {"ports": ["8000"]}, "env": ["k:v"]}`))
	assert.Equal(t, e.Error(), "`name` is a required field")
	_, e = NewContainerCreateCommandIntent("chan", []byte(`{"name": "name", "tag": "3.6", "command": ["sh"], "network_options": {"ports": ["8000"]}, "env": ["k:v"]}`))
	assert.Equal(t, e.Error(), "`image` is a required field")
	_, e = NewContainerCreateCommandIntent("chan", []byte(`{"name": "python", "image": "python", "command": ["sh"], "network_options": {"ports": ["8000"]}, "env": ["k:v"]}`))
	assert.Equal(t, e.Error(), "`tag` is a required field")
	_, e = NewContainerCreateCommandIntent("chan", []byte(`{"name": "python", "tag": "3.6", "image": "python", "network_options": {"ports": ["8000"]}, "env": ["k:v"]}`))
	assert.Equal(t, e.Error(), "`command` cannot be empty")
	_, e = NewContainerCreateCommandIntent("chan", []byte(`{"image":"python","tag":"3.6","network_options":{"ports":["8000"]},"name":"django","command":["sleep","infinity"],"env":[]}`))
	assert.Equal(t, e, nil)
}
