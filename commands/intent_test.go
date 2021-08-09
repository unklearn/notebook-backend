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

func TestContainerExecuteCommandIntent(t *testing.T) {
	i, e := NewContainerExecuteCommandIntent("foo", []byte(`{"interactive": false, "use_tty": true, "timeout": 123, "command": ["cat", ">", "oyo.py"]}`))
	assert.Equal(t, e, nil)
	assert.Equal(t, i.ContainerId, "foo")
	assert.Equal(t, i.Interactive, false)
	assert.Equal(t, i.UseTty, true)
	assert.Equal(t, i.Timeout, 123)
	assert.Equal(t, i.Command, []string{"cat", ">", "oyo.py"})

	// Try with missing values
	i, e = NewContainerExecuteCommandIntent("foo", []byte(`{"command": ["bash"]}`))
	assert.Equal(t, e, nil)
	assert.Equal(t, i.ContainerId, "foo")
	assert.Equal(t, i.Interactive, false)
	assert.Equal(t, i.UseTty, false)
	assert.Equal(t, i.Command, []string{"bash"})

	// Try with wrong args
	_, e = NewContainerExecuteCommandIntent("foo", []byte(`{}`))
	assert.Equal(t, e.Error(), "command cannot be empty")
}
