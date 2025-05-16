package config

import (
	"context"
	"errors"
	"flag"
	"os"
	"uiren/pkg/logger"

	vlt "github.com/hashicorp/vault/api"

	"gopkg.in/yaml.v3"
)

type vault interface {
	GetValueV1(ctx context.Context, mountPath, key string) (map[string]interface{}, error)
}

type config struct {
	Configs   map[string]interface{} `yaml:"configs"`
	Vault     vault
	mountPath string
}

func init() {
	cfg = config{
		Configs: map[string]interface{}{},
	}

	err := ParseConfig()
	if err != nil {
		logger.Fatal("config error: ", err)
	}
}

var (
	cfg config
)

// WithVault add's VaultClient to config.
// GetValue will search in vault at first, and if there is no content,
// it will check config file.
// If VaultClient is not initialized, GetValue will search only in config file
func WithVault(v vault, mountPath string) {
	cfg.Vault = v
	cfg.mountPath = mountPath
}

// ParseConfig The function fills the configuration structure from a file,
// the path to which can be passed as a flag when starting the program executable file.
// If no file path is provided, the default path will be used.
// The function returns an error if the process fails.
func ParseConfig() error {
	cfgFilePath := flag.String("config", "test", "config")
	flag.Parse()

	bytes, err := os.ReadFile("./config/config_" + *cfgFilePath + ".yml")
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return err
	}
	return nil
}

// ParseFile The function fills the configuration structure from a file,
// the path to which must be passed as an argument.
// The function returns an error if the process fails.
func ParseFile(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(bytes, &cfg); err != nil {
		return err
	}
	return nil
}

// GetValue retrieves the value of the config variable named by the key.
// It returns the value, which will be empty if the variable is not present.
// To distinguish between an empty value and an unset value, use LookupValue.
func GetValue(key string) Value {
	if cfg.Vault != nil {
		val, err := cfg.Vault.GetValueV1(context.Background(), cfg.mountPath, key)
		if err != nil {
			if errors.Is(err, vlt.ErrSecretNotFound) {
				v, _ := LookupValue(key)
				return v
			}
			logger.Error("pkg.config.GetValue.GetValueV1 error ", err)
			return value{
				exists: false,
			}
		}
		v, ok := val["value"]
		return value{
			value:  v,
			exists: ok,
		}
	}
	v, _ := LookupValue(key)
	return v
}

// LookupValue retrieves the value of the config variable
// named by the key. If the variable is present in the config
// the value (which may be empty) is returned and the boolean
// is true. Otherwise the returned value will be empty and the
// boolean will be false.
func LookupValue(key string) (Value, bool) {
	res, ok := cfg.Configs[key]
	return value{
		value:  res,
		exists: ok,
	}, ok
}
