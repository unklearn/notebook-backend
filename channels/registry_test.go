package channels

import (
	"testing"
)

type dummyChannel struct {
	IChannel
}

func TestRegisterChannel(t *testing.T) {
	cr := ChannelRegistry{}
	// Adding non-existing channel
	var ch IChannel = dummyChannel{}
	e := cr.RegisterChannel("dummy", ch)
	if e != nil {
		t.Errorf("Registering a channel for the first time must not raise error")
	}
	// dup
	e = cr.RegisterChannel("dummy", dummyChannel{})
	if e == nil {
		t.Error("Expected error from duplicate channel")
	}
}

func TestDeregisterChannel(t *testing.T) {
	cr := ChannelRegistry{}
	// Adding non-existing channel
	var ch IChannel = dummyChannel{}

	// Missing init channel
	_, e := cr.DeregisterChannel("foo")
	if e == nil {
		t.Error("Deregister channel must return error if not initialzed")
	}

	cr.RegisterChannel("ma", ch)
	_, e = cr.DeregisterChannel("foo2")
	if e == nil {
		t.Error("Deregister channel must return error if channel cannot be found")
	}

	cr.RegisterChannel("dummy", ch)
	rc, e := cr.DeregisterChannel("dummy")
	if e != nil || rc != ch {
		t.Error("Channel registry must successfully deregister channel and return the stored channel")
	}

	// Try to register again, no error should appear
	e = cr.RegisterChannel("dummy", ch)
	if e != nil {
		t.Error(e.Error())
	}

}

func TestGetChannelById(t *testing.T) {
	cr := ChannelRegistry{}
	// Adding non-existing channel
	var ch IChannel = dummyChannel{}

	_, e := cr.GetChannelById("foo")
	if e == nil {
		t.Error("If channel map is not initialized, GetChannelById must return error")
	}
	cr.RegisterChannel("maa", ch)
	_, e = cr.GetChannelById("foo")
	if e == nil {
		t.Error("If channel does not exist, GetChannelById must return error")
	}
	cr.RegisterChannel("dummy", ch)
	ct, _ := cr.GetChannelById("dummy")
	if ct != ch {
		t.Error("Must return correct channel when id is passed correctly")
	}
}
