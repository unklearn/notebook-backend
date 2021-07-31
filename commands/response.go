package commands

type ContainerStatusResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}
