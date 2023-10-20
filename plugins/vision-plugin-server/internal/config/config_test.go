package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/atos-digital/NHSS-scigateway/internal/config"
)

func TestConfigDefaults(t *testing.T) {
	os.Clearenv()
	c := config.New()
	assert.Equal(t, "localhost", c.Host)
	assert.Equal(t, "8080", c.Port)
	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", c.DatabaseURL)
	assert.Equal(t, []string{"http://localhost:8080", "https://localhost:8080"}, c.AllowedOrigins)
}

func TestConfigEnv(t *testing.T) {
	os.Setenv("HOST", "example.com")
	os.Setenv("PORT", "1234")
	os.Setenv("DATABASE_URL", "postgres://env:postgres@localhost:5432/postgres?sslmode=disable")
	os.Setenv("ALLOWED_ORIGINS", "http://example.com:1234,https://example.com:1234")

	c := config.New()
	assert.Equal(t, "example.com", c.Host)
	assert.Equal(t, "1234", c.Port)
	assert.Equal(t, "postgres://env:postgres@localhost:5432/postgres?sslmode=disable", c.DatabaseURL)
	assert.Equal(t, []string{"http://example.com:1234", "https://example.com:1234"}, c.AllowedOrigins)
}
