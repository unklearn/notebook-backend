package commands

type ContainerStatusResponse struct {
	Id     string `json:"id"`
	Hash   string `json:"hash"`
	Status string `json:"status"`
}

type ContainerCommandStatusResponse struct {
	ExecId string `json:"exec_id"`
	CellId string `json:"cell_id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type SyncFileResponse struct {
	NotebookId string `json:"notebook_id"`
	FilePath   string `json:"file_path"`
	Content    string `json:"content"`
	CellId     string `json:"cell_id"`
	Error      string `json:"error,omitempty"`
}
