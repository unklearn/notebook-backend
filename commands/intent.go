package commands

// Intent objects
type ActionIntent interface {
	GetIntentName() string
}

type ContainerNetworkOptions struct{}

// An intent that is designed to store container
// creation configuration
type ContainerCreateCommandIntent struct {
	// Id of the connection, can be notebookId or userId
	ConnectionId string
	// Optional, can be used to sync with existing container
	// if the notebook sends in the command intent from front-end
	ContainerId string
	// The name of the container
	Name string
	// The image to use for the container
	Image string
	// The image tag for container, if not provided `latest` will be used
	ImageTag string
	// Repository address to pull image from, if left out, docker hub will be used
	RepoUrl string
	// Network options
	NetworkOptions ContainerNetworkOptions
	// Env variables in KEY:VALUE format
	EnvVars []string
	// Start command to use
	Command []string
}

func (i ContainerCreateCommandIntent) GetIntentName() string {
	return "ContainerCreateCommandIntent"
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

func (i ContainerExecuteCommandIntent) GetIntentName() string {
	return "ContainerExecuteCommandIntent"
}
