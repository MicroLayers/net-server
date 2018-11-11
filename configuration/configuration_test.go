package configuration_test

import (
	"io/ioutil"
	"net-server/configuration"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

func createTempTestConfiguration(t *testing.T, path string) (configuration.Configuration, yaml.MapSlice) {
	config := configuration.Configuration{}
	config.Server.Module = "mymodule.so"
	config.Server.Listen.Unix.Enabled = false
	config.Server.Listen.Unix.Protocol = "whatever"
	config.Server.Listen.Unix.Socket = "/dir/my.sock"
	config.Server.Listen.Tcp.Enabled = true
	config.Server.Listen.Tcp.Protocol = "vdm"
	config.Server.Listen.Tcp.Port = uint16(62546)

	confBytes, err := yaml.Marshal(config)
	assert.NoError(t, err, "YAML package should correctly marshal the configuration")
	var confMapSlice yaml.MapSlice
	err = yaml.Unmarshal(confBytes, &confMapSlice)
	assert.NoError(t, err, "YAML package should correctly unmarshal the configuration into a MapSlice")

	tmpFile, err := os.Create(path)
	assert.NoError(t, err, "It should be possible to create the configuration file")
	defer tmpFile.Close()
	_, err = tmpFile.Write(confBytes)
	assert.NoError(t, err, "It should be possible to write the temporary file")
	tmpFile.Close()

	return config, confMapSlice
}

func TestReadConfWillReadAnExistingFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "net_server_test_")
	assert.NoError(t, err, "OS should succesfully create a temporary file")
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath)

	config, confMapSlice := createTempTestConfiguration(t, tmpFilePath)

	readConfig, mapSlice := configuration.ReadConf(tmpFilePath)
	assert.Equal(t, config, readConfig, "The written and read configurations should match")
	assert.Equal(t, confMapSlice, mapSlice, "The written and read MapSlice configurations should match")
}

func TestReadConfWillReadTheConfigurationFromTheDefaultPath(t *testing.T) {
	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, "configuration.yml")

	config, confMapSlice := createTempTestConfiguration(t, path)
	readConfig, mapSlice := configuration.ReadConf(path)

	assert.Equal(t, config, readConfig, "The written and read configurations should match")
	assert.Equal(t, confMapSlice, mapSlice, "The written and read MapSlice configurations should match")
}

func TestReadConfWillReturnDefaultConfigurationIfFileDoesNotExist(t *testing.T) {
	defaultConfiguration := configuration.DefaultConfiguration()

	readConfig, mapSlice := configuration.ReadConf(filepath.Join(os.TempDir(), "non-existing-file.whatever"))
	assert.Equal(t, defaultConfiguration, readConfig, "Reading a non-existing configuration should return the default one")
	assert.Nil(t, mapSlice, "Reading a non-existing configuration should return the default map slice")
}
