package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	app "github.com/andreyAKor/otus_arch_project/internal/app/pending"
	configs "github.com/andreyAKor/otus_arch_project/internal/configs/pending"
	"github.com/andreyAKor/otus_arch_project/internal/logging"
	"github.com/andreyAKor/otus_arch_project/internal/repository/psql"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "pending",
	Short: "Pending service application",
	Long:  "The Pending service for checking pending bids",
	RunE:  run,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&cfgFile, "config", "", "config file")
	if err := cobra.MarkFlagRequired(pf, "config"); err != nil {
		//nolint:forbidigo
		fmt.Println(err)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		//nolint:forbidigo
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init config
	c := new(configs.Config)
	if err := c.Init(cfgFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	// Init logger
	l := logging.New(c.Logging.File, c.Logging.Level)
	if err := l.Init(); err != nil {
		return errors.Wrap(err, "init logging failed")
	}

	log.Info().Msg("Starting...")

	repo, err := psql.New()
	if err != nil {
		return errors.Wrap(err, "can't initialize postgress repository")
	}
	if err := repo.Connect(ctx, c.Database.DSN); err != nil {
		return errors.Wrap(err, "postgress connection error")
	}

	// Init and run app
	a, err := app.New(repo, c)
	if err != nil {
		return errors.Wrap(err, "can't initialize app")
	}
	if err := a.Run(ctx); err != nil {
		return errors.Wrap(err, "app runnign fail")
	}

	log.Info().Msg("Running...")

	// Graceful shutdown
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)
	<-interruptCh

	log.Info().Msg("Stopping...")

	if err := a.Close(); err != nil {
		return errors.Wrap(err, "app closing fail")
	}

	log.Info().Msg("Stopped")

	if err := l.Close(); err != nil {
		return errors.Wrap(err, "logger closing fail")
	}

	return nil
}
