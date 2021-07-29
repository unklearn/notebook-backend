package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerCommandCreateIntent(t *testing.T) {
	c := ContainerCreateCommandIntent{}
	b := []byte(`{"container_id": "foo", "network_options": {"ports": ["8000"]}, "env": ["k:v"]}`)
	e := c.Parse("chan", b)
	if e != nil {
		t.Errorf("Expected error to be nil, except got error %s", e.Error())
	}
	assert.Equal(t, c.ChannelId, "chan")
	assert.Equal(t, c.ContainerId, "foo")
	assert.Equal(t, len(c.NetworkOptions.Ports), 1)
	assert.Equal(t, c.NetworkOptions.Ports[0], "8000")
	assert.Equal(t, c.EnvVars[0], "k:v")
	// Try with error
	e = c.Parse("chan", []byte(`{"container_id": 1234}`))
	assert.NotEqual(t, e, nil, "Must fail with error")
}
