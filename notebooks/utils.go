package notebooks

import "regexp"

var matcher = regexp.MustCompile(`[^A-Za-z\d-]+`).ReplaceAllString

// sanitize notebook id to prevent path escapes. Only alpha-number + "-" character is allowed
func sanitizeNotebookId(notebookId string) string {
	return matcher(notebookId, "")
}
