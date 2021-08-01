package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Intent objects
type ActionIntent interface {
	// Parse message payload
	GetIntentName() string
	ToString() string
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

// Print string version of intent
func (i ContainerCreateCommandIntent) ToString() string {
	return fmt.Sprintf("%#v", i)
}

// Factory method that supplies new intent or error, when supplied with JSON
// representation of body
func NewContainerCreateCommandIntent(channelId string, payload []byte) (ContainerCreateCommandIntent, error) {
	i := ContainerCreateCommandIntent{
		ChannelId: channelId,
	}
	err := json.Unmarshal(payload, &i)
	if err != nil {
		log.Printf("Error while unmarshalling container input: %s", err.Error())
		return i, fmt.Errorf("invalid input supplied for creating container")
	}
	errors := []string{}
	// Check for other parameters
	if i.Name == "" {
		errors = append(errors, "`name` is a required field")
	}
	if i.Image == "" {
		errors = append(errors, "`image` is a required field")
	}
	if i.ImageTag == "" {
		errors = append(errors, "`tag` is a required field")
	}
	if len(i.Command) == 0 {
		errors = append(errors, "`command` cannot be empty")
	}
	if len(errors) > 0 {
		return i, fmt.Errorf(strings.Join(errors, "\n"))
	}
	return i, nil
}

// Ensure that the provided image and tag exists on the system
type ImagePullCommandIntent struct {
	Image   string
	Tag     string
	RepoUrl string
}

func (i ImagePullCommandIntent) GetIntentName() string {
	return "ImagePullCommandIntent"
}

// Wait for a container to be started, sometimes images do not exist, and images
// must be pulled
type ContainerWaitCommandIntent struct {
	ContainerId string
	// Timeout in seconds
	Timeout int
}

func (i ContainerWaitCommandIntent) GetIntentName() string {
	return "ContainerWaitCommandIntent"
}

func (i ContainerWaitCommandIntent) ToString() string {
	return fmt.Sprintf("%#v", i)
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
