package pedding

import (
	"github.com/pkg/errors"

	"github.com/andreyAKor/otus_arch_project/internal/configs"
)

var (
	ErrCoinTypeNotFound = errors.New("coin type not found")

	_ configs.Configer = (*Config)(nil)
)

type Config struct {
	App       configs.App
	Logging   configs.Logging
	Nodes     configs.Nodes
	Database  configs.Database
	Addresses configs.Addresses
}

func (c *Config) Init(file string) error {
	cfg, err := configs.Init(file, c)

	_, ok := cfg.(*Config)
	if !ok {
		return errors.Wrap(err, "init config failed")
	}

	return nil
}

func (c *Config) GetAddress(coinType uint64) (string, error) {
	addr, ok := c.Addresses[coinType]
	if !ok {
		//nolint:wrapcheck
		return "", ErrCoinTypeNotFound
	}

	return addr, nil
}

func (c *Config) GetNodes() configs.Nodes {
	return c.Nodes
}
