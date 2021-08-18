package notebooks

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

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
	dir := "/usr/local/var/notebooks"
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

// Update a notebook by saving its new contents
func (nb *NotebookCRUDService) Update(notebookId string, payload map[string]interface{}) (map[string]interface{}, error) {
	notebookFilePath := filepath.Join(nb.rootDir, sanitizeNotebookId(notebookId))
	_, err := nb.fs.Stat(notebookFilePath)
	if err != nil {
		return nil, err
	} else {
		file, err := nb.fs.OpenFile(notebookFilePath, os.O_WRONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}
		contents, _ := json.Marshal(payload)
		_, err = file.Write(contents)
		if err != nil {
			return nil, err
		}
		return payload, nil
	}
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
