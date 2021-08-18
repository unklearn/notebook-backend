package notebooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeNotebookId(t *testing.T) {
	assert.Equal(t, sanitizeNotebookId("../abcd-def"), "abcd-def")
	assert.Equal(t, sanitizeNotebookId("../../etc/pwd/abcd-123-def"), "etcpwdabcd-123-def")
	assert.Equal(t, sanitizeNotebookId("../../etc/pwd/abcd-def"), "etcpwdabcd-def")
}
