package bid

import (
	"context"
	"io"

	"github.com/micro/go-micro"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/otus_arch_project/internal/blockchains"
	"github.com/andreyAKor/otus_arch_project/internal/blockchains/btc"
	"github.com/andreyAKor/otus_arch_project/internal/blockchains/eth"
	configs "github.com/andreyAKor/otus_arch_project/internal/configs/bid"
	"github.com/andreyAKor/otus_arch_project/internal/consts"
	grpcServersBid "github.com/andreyAKor/otus_arch_project/internal/grpc/servers/bid"
	"github.com/andreyAKor/otus_arch_project/internal/repository"
)

var _ io.Closer = (*App)(nil)

type App struct {
	microService  micro.Service
	grpcBidServer *grpcServersBid.Bid
}

func New(microService micro.Service, rsql repository.DBBidsRepo, conf *configs.Config) (*App, error) {
	e, err := eth.New(conf)
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize eth")
	}

	b, err := btc.New(conf)
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize btc")
	}

	// Init grpc bid server
	grpcBidServer, err := grpcServersBid.New(
		microService.Server(),
		rsql,
		conf,
		map[uint64]blockchains.Blockchainer{
			consts.CoinTypeETH: e,
			consts.CoinTypeBTC: b,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize grpc bid server")
	}

	return &App{microService, grpcBidServer}, nil
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.microService.Run(); err != nil {
			log.Fatal().Err(err).Msg("micro service run fail")
		}
	}()

	if err := a.grpcBidServer.Run(ctx); err != nil {
		return errors.Wrap(err, "grpc bid server run fail")
	}

	return nil
}

func (a *App) Close() error {
	return nil
}
