package notebooks

import (
	"github.com/HouzuoGuo/tiedot/db"
)

type NotebookCRUDService struct {
	database   *db.DB
	collection *db.Col
}

func NewNotebookCRUDService(dbDir string) *NotebookCRUDService {

	myDB, err := db.OpenDB(dbDir)
	if err != nil {
		// Panic and exit immediately
		panic(err)
	}
	var collection *db.Col
	// Check if collection exists, else create it
	if !myDB.ColExists("notebooks") {
		myDB.Create("notebooks")
		collection = myDB.Use("notebooks")
	}
	return &NotebookCRUDService{database: myDB, collection: collection}
}

// Create a new notebook and return the id and error if any.
func (nb *NotebookCRUDService) Create(payload map[string]interface{}) (int, error) {
	docId, err := nb.collection.Insert(payload)
	if err != nil {
		return 0, err
	}
	return docId, err
}

// Fetch a notebook by its id and return the notebook contents or error
func (nb *NotebookCRUDService) GetById(docId int) (map[string]interface{}, error) {
	readBack, err := nb.collection.Read(docId)
	if err != nil {
		return nil, err
	}
	return readBack, nil
}
