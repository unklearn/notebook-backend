package commands

import "encoding/json"

// Intent objects
type ActionIntent interface {
	// Parse message payload
	GetIntentName() string
}

type ContainerNetworkOptions struct {
	Ports []string `json:"ports"`
}

// An intent that is designed to store container
// creation configuration
type ContainerCreateCommandIntent struct {
	// Id of the connection, can be notebookId or userId
	ChannelId string `json:"-"`
	// Optional, can be used to sync with existing container
	// if the notebook sends in the command intent from front-end
	ContainerId string `json:"container_id"`
	// The name of the container
	Name string `json:"name"`
	// The image to use for the container
	Image string `json:"image"`
	// The image tag for container, if not provided `latest` will be used
	ImageTag string `json:"tag"`
	// Repository address to pull image from, if left out, docker hub will be used
	RepoUrl string `json:"repo_url"`
	// Network options
	NetworkOptions ContainerNetworkOptions `json:"network_options"`
	// Env variables in KEY:VALUE format
	EnvVars []string `json:"env"`
	// Start command to use
	Command []string `json:"command"`
}

func (i ContainerCreateCommandIntent) GetIntentName() string {
	return "ContainerCreateCommandIntent"
}

// Parse payload into container create command intent or return error if config is invalid
func (i *ContainerCreateCommandIntent) Parse(channelId string, payload []byte) error {
	err := json.Unmarshal(payload, i)
	i.ChannelId = channelId
	if err != nil {
		return err
	}
	return nil
}

// An intent that can be used to execute a command inside container
type ContainerExecuteCommandIntent struct {
	// Id of the connection, can be notebookId or userId
	ConnectionId string
	// Id of container
	ContainerId string
	// Whether command can accept stdin
	Interactive bool
	// Whether command requires tty
	UseTty bool
	// Used to timeout commands if necessary
	Timeout int64
}
