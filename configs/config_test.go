package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDefault(t *testing.T) {
	// Must change directory as the app will start from the base container folder
	// Config load will also be tested in main, but in more general way (just no errors probably)
	err := os.Chdir("../")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	cfg, err := LoadDefault()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	if !assert.NotEmpty(t, cfg) {
		t.FailNow()
	}

	assert.NotEmpty(t, cfg.Server)
	assert.NotEmpty(t, cfg.Database)
	assert.NotEmpty(t, cfg.Logger)

	// Should not change because app will be executed in container environments only
	assert.Equal(t, cfg.Server.ListenAddr, "0.0.0.0")
}
