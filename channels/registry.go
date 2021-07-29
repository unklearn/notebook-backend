package channels

import (
	"errors"
	"fmt"
)

type ChannelRegistry struct {
	// Internal store for mapping channelId to channel
	channelMap map[string]IChannel
}

// RegisterChannel registers a channel against a channelId.
// If a channel exists for given channelId, it returns an error
func (cr *ChannelRegistry) RegisterChannel(channelId string, channel IChannel) error {
	if cr.channelMap == nil {
		cr.channelMap = make(map[string]IChannel)
	}
	_, ok := cr.channelMap[channelId]
	// If another channel exists, return error
	if ok {
		return fmt.Errorf("ECODE::dup-channel::There exists another channel for channelId %s", channelId)
	}
	cr.channelMap[channelId] = channel
	return nil
}

// Deregister channel removes a channel from the store if it exists,
// otherwise returns error
func (cr *ChannelRegistry) DeregisterChannel(channelId string) (IChannel, error) {
	if cr.channelMap == nil {
		return nil, errors.New("ECODE::missing-map::Registry has not been initialized")
	}
	ch, ok := cr.channelMap[channelId]
	if ok {
		delete(cr.channelMap, channelId)
		return ch, nil
	}
	return nil, fmt.Errorf("ECODE::missing-channel::There exists no channel with channelId %s", channelId)
}

// Return a channel by id if it exists, otherwise return error
func (cr *ChannelRegistry) GetChannelById(channelId string) (IChannel, error) {
	if cr.channelMap == nil {
		return nil, errors.New("ECODE::missing-map::Registry has not been initialized")
	}
	ch, ok := cr.channelMap[channelId]
	if ok {
		return ch, nil
	}
	return nil, fmt.Errorf("ECODE::missing-channel::There exists no channel with channelId %s", channelId)
}
