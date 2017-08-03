package mysqlx

import (
	"encoding/hex"
	"net/url"
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

func TestSetDefaults(t *testing.T) {
	u, err := url.Parse("")
	require.NoError(t, err)
	require.NoError(t, setDefaults(u))
	assert.Equal(t, "tcp://localhost:33060", u.String())

	u, err = url.Parse("mysql:mysql")
	require.NoError(t, err)
	assert.EqualError(t, setDefaults(u), "invalid data source: mysql:mysql")
}
