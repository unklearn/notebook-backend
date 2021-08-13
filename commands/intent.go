package commands

import (
	"encoding/json"
	"errors"
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
	// Hash for tracking which request corresponds to failure
	Hash string `json:"hash"`
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
	if i.Hash == "" {
		errors = append(errors, "`hash` is a required field")
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
	// Id of container
	ContainerId string `json:"-"`
	// The id of the cell to correlate command execution inside notebook.
	// without cell id, it is hard to identify which execId corresponds
	// to which cell id
	CellId string `json:"cell_id"`
	// Whether command can accept stdin
	Interactive bool `json:"interactive,omitempty"`
	// Whether command requires tty
	UseTty bool `json:"use_tty,omitempty"`
	// Used to timeout commands if necessary
	Timeout int `json:"timeout,omitempty"`
	// The command to execute along with args
	Command []string `json:"command"`
}

func (i ContainerExecuteCommandIntent) GetIntentName() string {
	return "ContainerExecuteCommandIntent"
}

func (i ContainerExecuteCommandIntent) ToString() string {
	return fmt.Sprintf("%#v", i)
}

// Function to create a new container execute command intent
func NewContainerExecuteCommandIntent(containerId string, payload []byte) (ContainerExecuteCommandIntent, error) {
	c := ContainerExecuteCommandIntent{ContainerId: containerId}
	e := json.Unmarshal(payload, &c)
	if len(c.Command) == 0 {
		e = errors.New("command cannot be empty")
	}
	if e != nil {
		return ContainerExecuteCommandIntent{}, e
	}
	return c, nil
}

// SyncFileIntent syncs the file from server onto the client
type SyncFileIntent struct {
	// Id of the container
	ContainerId string `json:"-"`
	// Path to file including extension
	FilePath string `json:"file_path"`
	// Id of the cell
	CellId string `json:"cell_id"`
	// Optional content if notebook is writing in
	Content string `json:"content,omitempty"`
}

func (i SyncFileIntent) GetIntentName() string {
	return "SyncFileIntent"
}

func (i SyncFileIntent) ToString() string {
	return fmt.Sprintf("%#v", i)
}

// Constructor function for sync file intent
func NewSyncFileIntent(containerId string, payload []byte) (SyncFileIntent, error) {
	si := SyncFileIntent{}
	err := json.Unmarshal(payload, &si)
	if err != nil {
		return si, err
	}
	errors := []string{}
	if len(si.FilePath) == 0 {
		errors = append(errors, "`file_path` cannot be empty")
	}
	if si.CellId == "" {
		errors = append(errors, "`cell_id` is a required field")
	}
	if len(errors) > 0 {
		return si, fmt.Errorf(strings.Join(errors, "\n"))
	}
	si.ContainerId = containerId
	return si, nil
}
