package gateway

import (
	"context"
	"io"

	"github.com/micro/go-micro"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	grpcClientsBid "github.com/andreyAKor/otus_arch_project/internal/grpc/clients/bid"
	"github.com/andreyAKor/otus_arch_project/internal/http/server"
)

var _ io.Closer = (*App)(nil)

type App struct {
	httpSrv      *server.Server
	microService micro.Service
}

func New(microService micro.Service, host string, port, bodyLimit int) (*App, error) {
	// Init grpc bid client
	grpcBidClient, err := grpcClientsBid.New(microService.Client())
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize grpc bid client")
	}

	// Init http-server
	httpSrv, err := server.New(host, port, bodyLimit, grpcBidClient)
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize http-server")
	}

	return &App{httpSrv, microService}, nil
}

// Run application.
func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.httpSrv.Run(ctx); err != nil {
			log.Fatal().Err(err).Msg("http-server run fail")
		}
	}()
	go func() {
		if err := a.microService.Run(); err != nil {
			log.Fatal().Err(err).Msg("micro service run fail")
		}
	}()

	return nil
}

// Close application.
func (a *App) Close() error {
	return a.httpSrv.Close()
}
