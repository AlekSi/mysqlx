package mysqlx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeverityStringer(t *testing.T) {
	assert.Equal(t, "ERROR", SeverityError.String())
	assert.Equal(t, "FATAL", SeverityFatal.String())
	assert.Equal(t, "Severity 42", Severity(42).String())
}
