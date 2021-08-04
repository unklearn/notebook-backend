package notebooks

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var nbService = NewNotebookCRUDService("/tmp/notebooks")

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func handleNotebookCreate(w http.ResponseWriter, r *http.Request) {
	p := make([]byte, r.ContentLength)
	r.Body.Read(p)
	var payload map[string]interface{}
	json.Unmarshal(p, &payload)
	id, err := nbService.Create(payload)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid notebook payload")
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{"id": id})
}

func handleNotebookGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notebookId, err := strconv.Atoi(vars["notebookId"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid notebookId")
	}
	n, e := nbService.GetById(notebookId)
	if e != nil {
		respondWithError(w, http.StatusNotFound, "Cannot find notebook")
	}
	respondWithJSON(w, http.StatusOK, n)
}

func NotebooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleNotebookCreate(w, r)
	}
}

func NotebookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleNotebookGet(w, r)
	}
}
