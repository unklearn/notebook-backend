package commands

type ContainerStatusResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type ContainerCommandStatusResponse struct {
	ExecId string `json:"exec_id"`
	CellId string `json:"cell_id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}
