package glesys

import (
	"github.com/docker/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfigFromFlags(t *testing.T) {
	driver := NewDriver("default", "path")

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"glesys-api-key":       "TOKEN",
			"glesys-bandwidth":     100,
			"glesys-campaign-code": "CODE",
			"glesys-cpu":           2,
			"glesys-data-center":   "Falkenberg",
			"glesys-storage":       40,
			"glesys-memory":        2048,
			"glesys-project":       "cl12345",
			"glesys-root-password": "secret-password",
			"glesys-template":      "Debian 8 64-bit",
			"glesys-ssh-key-path":  "~/.ssh/id_rsa",
		},
		CreateFlags: driver.GetCreateFlags(),
	}

	err := driver.SetConfigFromFlags(checkFlags)

	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)
}

func TestConfigRequiresProject(t *testing.T) {
	driver := NewDriver("default", "path")

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"glesys-api-key": "TOKEN",
		},
		CreateFlags: driver.GetCreateFlags(),
	}

	err := driver.SetConfigFromFlags(checkFlags)
	assert.EqualError(t, err, "glesys driver requires the --glesys-project option")
}

func TestConfigRequiresApiKey(t *testing.T) {
	driver := NewDriver("default", "path")

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"glesys-project": "cl12345",
		},
		CreateFlags: driver.GetCreateFlags(),
	}

	err := driver.SetConfigFromFlags(checkFlags)
	assert.EqualError(t, err, "glesys driver require the --glesys-api-key option")
}

func TestStringToEnvVar(t *testing.T) {
	assert.Equal(t, "GLESYS_API_KEY", stringFlagToEnvVar("glesys-api-key"), "return value is correct")
	assert.Equal(t, "GLESYS_PROJECT", stringFlagToEnvVar("glesys-project"), "return value is correct")
}

func TestDriverName(t *testing.T) {
	driver := NewDriver("default", "path")

	assert.Equal(t, "glesys", driver.DriverName(), "Driver name is correct")
}

func TestSSHPort(t *testing.T) {
	driver := NewDriver("default", "path")
	port, err := driver.GetSSHPort()

	assert.NoError(t, err, "there is no error")
	assert.Equal(t, 22, port, "SSH port is correct")
}

func TestSSHUsername(t *testing.T) {
	driver := NewDriver("default", "path")

	assert.Equal(t, "root", driver.GetSSHUsername(), "SSH username is correct")
}

func TestGeneratePassword(t *testing.T) {
	password := generatePassword(64)

	assert.Len(t, password, 64, "password has correct length")
}

func TestDriver_PreCreateCheck(t *testing.T) {
	driver := NewDriver("default", "path")
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
}
