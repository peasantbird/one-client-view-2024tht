package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-api/internal/config"
)

func TestNewConfig(t *testing.T) {
	c := config.NewConfig()

	assert.NotNil(t, c)
}

func TestLoad(t *testing.T) {
	t.Cleanup(func() {
		os.Remove(".env")
	})

	err := os.WriteFile(".env", []byte(`
		PORT=8080
		DB_HOST=localhost
		DB_USER=test
		DB_NAME=test
		DB_PASSWORD=test
		DB_PORT=5432
		DB_SSLMODE=disable
		DB_TIMEZONE=UTC
	`), 0644)
	require.NoError(t, err)

	c := config.NewConfig()
	c.Load()

	assert.Equal(t, "8080", c.Port)
	assert.Equal(t, "host=localhost user=test dbname=test password=test port=5432 sslmode=disable TimeZone=UTC", c.DB.DSN)
}
