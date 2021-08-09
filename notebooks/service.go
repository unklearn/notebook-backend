package notebooks

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

type NotebookCRUDService struct {
	fs      afero.Fs
	rootDir string
}

func NewNotebookCRUDService(dbDir string) *NotebookCRUDService {
	var AppFs = afero.NewOsFs()
	// dir := afero.GetTempDir(AppFs, "notebooks")
	dir := "/tmp/notebooks"
	log.Printf("Using directory %s\n", dir)
	return &NotebookCRUDService{fs: AppFs, rootDir: dir}
}

// Create a new notebook and return the id and error if any.
func (nb *NotebookCRUDService) Create(payload map[string]interface{}) (map[string]interface{}, error) {
	uid, _ := uuid.NewUUID()
	file, err := nb.fs.Create(nb.rootDir + "/" + uid.String())
	// Write json to the file
	payload["id"] = uid.String()
	payload["containers"] = make([]interface{}, 0)
	payload["cells"] = make([]interface{}, 0)
	contents, _ := json.Marshal(payload)
	file.Write(contents)
	if err != nil {
		return nil, err
	}
	return payload, err
}

// Fetch a notebook by its id and return the notebook contents or error
func (nb *NotebookCRUDService) GetById(docId string) (map[string]interface{}, error) {
	readBack, err := afero.ReadFile(nb.fs, nb.rootDir+"/"+docId)
	var s map[string]interface{}
	json.Unmarshal(readBack, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
