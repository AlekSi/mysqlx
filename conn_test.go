package mysqlx

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthData(t *testing.T) {
	decode := func(s string) []byte {
		b, err := hex.DecodeString(s)
		require.NoError(t, err)
		return b
	}

	assert.Equal(t, "world_x\x00root\x00",
		string(authData("world_x", "root", "", decode("434169533f3569721167252e59117a645a681500"))))
	assert.Equal(t, "world_x\x00my_user\x00*c8b66adb21e1e674249869852aebc573db7a5639",
		string(authData("world_x", "my_user", "my_password", decode("434169533f3569721167252e59117a645a681500"))))
}
