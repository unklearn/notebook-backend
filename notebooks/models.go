package notebooks

import "time"

// A struct representing notebook row. Updates are made directly to the JSON row.
// Cells can be loaded and paginated through
type NotebookCell struct {
	Id        string    `json:"id"`         // id of the cell
	Status    string    `json:"status"`     // status if cell is executable
	Type      string    `json:"type"`       // type of the cell
	Contents  string    `json:"contents"`   // cell contents
	Output    string    `json:"output"`     // output of the cell if any
	CreatedAt time.Time `json:"created_at"` // Time the cell was created at
	PodId     string    `json:"pod_id"`     // id of the pod cell is connected to
}

type Notebook struct {
	Id        string         `json:"id"`         // id of the notebook
	Name      string         `json:"name"`       // name of the notebook
	CreatedBy string         `json:"created_by"` // id of the user who created the notebook
	Cells     []NotebookCell `json:"cells"`      // Array of cells in the notebook
}
