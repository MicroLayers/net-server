package module_test

import (
	"net-server/module"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadModuleWillLoadAValidModule(t *testing.T) {
	pwd, _ := os.Getwd()
	modulePath := filepath.Join(pwd, "examples", "echo-module.so")

	mod, err := module.LoadModule(modulePath)

	assert.NoError(t, err, "The loader should load the module without errors")
	assert.NotNil(t, mod, "The loader should succesfully load the module")
}

func TestLoadModuleWillNotLoadModuleNotExposingNetServerModule(t *testing.T) {
	pwd, _ := os.Getwd()
	modulePath := filepath.Join(pwd, "examples", "invalid-module-no-symbol.so")

	mod, err := module.LoadModule(modulePath)
	assert.Nil(t, mod, "It should not be possible to load a module not containing the NetServerModule symbol")
	assert.Error(t, err, "It should not be possible to load a module not containing the NetServerModule symbol")
}

func TestLoadModuleWillNotLoadModuleHavingAnInvalidNetServerModuleSymbol(t *testing.T) {
	pwd, _ := os.Getwd()
	modulePath := filepath.Join(pwd, "examples", "invalid-module-wrong-symbol-type.so")

	mod, err := module.LoadModule(modulePath)
	assert.Nil(t, mod, "It should not be possible to load a module with wrong symbol type")
	assert.Error(t, err, "It should not be possible to load a module with wrong symbol type")
}

func TestLoadModuleWillNotLoadANonExistingModule(t *testing.T) {
	mod, err := module.LoadModule(filepath.Join(os.TempDir(), "net_server_non_existing_module_file"))
	assert.Nil(t, mod, "It should not be possible to load an invalid module file")
	assert.Error(t, err, "It should not be possible to load an invalid module file")
}
