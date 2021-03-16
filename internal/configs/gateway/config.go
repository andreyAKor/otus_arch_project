package gateway

import (
	"github.com/pkg/errors"

	"github.com/andreyAKor/otus_arch_project/internal/configs"
)

var _ configs.Configer = (*Config)(nil)

type Config struct {
	App     configs.App
	Micro   configs.Micro
	Logging configs.Logging

	// HTTP-server settings
	HTTP struct {
		// Host
		Host string

		// Port
		Port int

		// Maximum content size limit
		BodyLimit int
	}
}

func (c *Config) Init(file string) error {
	cfg, err := configs.Init(file, c)

	_, ok := cfg.(*Config)
	if !ok {
		return errors.Wrap(err, "init config failed")
	}

	return nil
}

func (c *Config) GetAddress(coinType uint64) (_ string, _ error) {
	return
}

func (c *Config) GetNodes() (_ configs.Nodes) {
	return
}
