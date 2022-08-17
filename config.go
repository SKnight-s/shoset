package shoset

import (
	"os"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config: collection of configuration information for a shoset.
type Config struct {
	baseDirectory string
	fileName      string
	viper         *viper.Viper
	mu            sync.Mutex
}

// NewConfig returns a *Config object.
// Initialize home directory and viper.
func NewConfig(name string) *Config {
	homeDirectory := "."
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Error().Msg("couldn't get home directory folder: " + err.Error())
		return nil
	}
	cfg := &Config{
		baseDirectory: homeDirectory + "/.shoset/" + name + "/",
		viper:         viper.New(),
	}
	if err := mkdir(cfg.baseDirectory); err != nil {
		log.Error().Msg("couldn't create folder: " + err.Error())
		return nil
	}
	return cfg
}

// GetBaseDirectory returns baseDirectory from config, baseDirectory corresponds to homeDirectory + shoset repertory.
func (cfg *Config) GetBaseDirectory() string { return cfg.baseDirectory }

// GetFileName returns fileName from config.
func (cfg *Config) GetFileName() string { return cfg.fileName }

// SetFileName sets fileName for a config.
func (cfg *Config) SetFileName(fileName string) { cfg.fileName = fileName }

// InitFolders creates following config folders if needed: <base>/<name>/{cert, config}.
func (cfg *Config) InitFolders(name string) (string, error) {
	if err := mkdir(cfg.baseDirectory + name + "/"); err != nil {
		return VOID, err
	}
	if err := mkdir(cfg.baseDirectory + name + PATH_CONFIG_DIRECTORY); err != nil {
		return VOID, err
	}
	if err := mkdir(cfg.baseDirectory + name + PATH_CERT_DIRECTORY); err != nil {
		return VOID, err
	}
	return cfg.baseDirectory, nil
}

// ReadConfig will load the configuration file from disk for a given fileName.
// Initialize viper parameters before reading.
func (cfg *Config) ReadConfig(fileName string) error {

	cfg.viper.AddConfigPath(cfg.baseDirectory + fileName + PATH_CONFIG_DIRECTORY)
	cfg.viper.SetConfigName(CONFIG_FILE)
	cfg.viper.SetConfigType(CONFIG_TYPE)

	err := cfg.viper.ReadInConfig()

	return err
}

// WriteConfig writes current configuration to a given shoset.
func (cfg *Config) WriteConfig(fileName string) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return cfg.viper.WriteConfigAs(cfg.baseDirectory + fileName + PATH_CONFIG_DIRECTORY + CONFIG_FILE)
}

// AppendToKey sets the value for a key for viper config.
func (cfg *Config) AppendToKey(key string, values []string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	// Avoid duplicate
	valueSlice := cfg.GetSlice(key)
	for _, a := range values {
		if !contains(valueSlice, a) {
			valueSlice = append(valueSlice, a)
		}
	}
	cfg.viper.Set(key, valueSlice)
	cfg.viper.WriteConfig()
}

// DeleteFromKey deletes the vales from the list from the key.
func (cfg *Config) DeleteFromKey(key string, value []string) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	valueCfg := cfg.GetSlice(key)
	valueCfgOut := []string{}
	for _, a := range valueCfg {
		if !contains(value, a) {
			valueCfgOut = append(valueCfgOut, a)
		}
	}
	cfg.viper.Set(key, valueCfgOut)
	cfg.viper.WriteConfig()
}

// GetSlice returns the viper config for a specific protocol.
func (cfg *Config) GetSlice(protocol string) []string {
	return cfg.viper.GetStringSlice(protocol)
}
