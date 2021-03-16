package configs

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Configer interface {
	Init(file string) error
	GetAddress(coinType uint64) (string, error)
	GetNodes() Nodes
}

// Application information.
type App struct {
	Name string // Service name
}

// Data to initialize go-micro.
type Micro struct {
	// Data for working with the registry
	Registry        string
	RegistryAddress string

	// Data for working with a broker
	Broker        string
	BrokerAddress string

	// Data for working with transport
	Transport        string
	TransportAddress string
}

// Logging settings.
type Logging struct {
	// Path to the log file.
	File string

	// Logging level, variants levels:
	//  - debug - defines debug log level
	//  - info - defines info log level
	//  - warn - defines warn log level
	//  - error - defines error log level
	//  - fatal - defines fatal log level
	//  - panic - defines panic log level
	//  - no - defines an absent log level
	//  - disabled - disables the logger
	//  - trace - defines trace log level.
	Level string
}

// Blockchains nodes params.
type Nodes struct {
	ETH struct {
		Host string
		Port int64
	}
	BTC struct {
		Host string
		Port int64
		User string
		Pass string
	}
}

// Database settings.
type Database struct {
	// DSN string for database connection.
	DSN string
}

// Bid addresses.
type Addresses map[uint64]string

// Init is using to initialize the current config instance.
func Init(file string, c Configer) (Configer, error) {
	// read in environment variables that match
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "open config file failed")
	}

	if err := viper.Unmarshal(c); err != nil {
		return nil, errors.Wrap(err, "unmarshal config file failed")
	}

	return c, nil
}
