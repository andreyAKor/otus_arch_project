package pedding

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/otus_arch_project/internal/blockchains"
	"github.com/andreyAKor/otus_arch_project/internal/blockchains/btc"
	"github.com/andreyAKor/otus_arch_project/internal/blockchains/eth"
	configs "github.com/andreyAKor/otus_arch_project/internal/configs/pedding"
	"github.com/andreyAKor/otus_arch_project/internal/consts"
	"github.com/andreyAKor/otus_arch_project/internal/pkg/workers"
	"github.com/andreyAKor/otus_arch_project/internal/repository"
)

var (
	ErrBlockchainNotFound = errors.New("blockchain not found")

	checkPeddingBids = time.Second * 10 // Проверка ожидающих заявок

	_ io.Closer = (*App)(nil)
)

type App struct {
	repo        repository.DBBidsRepo
	conf        *configs.Config
	blockchains map[uint64]blockchains.Blockchainer
	cancel      context.CancelFunc
	processMx   sync.Mutex
}

func New(repo repository.DBBidsRepo, conf *configs.Config) (*App, error) {
	e, err := eth.New(conf)
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize eth")
	}

	b, err := btc.New(conf)
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize btc")
	}

	return &App{
		repo: repo,
		conf: conf,
		blockchains: map[uint64]blockchains.Blockchainer{
			consts.CoinTypeETH: e,
			consts.CoinTypeBTC: b,
		},
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	ctx, a.cancel = context.WithCancel(ctx)
	workers.Timeout(ctx, checkPeddingBids, a.process)

	return nil
}

func (a *App) Close() error {
	a.cancel()

	return nil
}

func (a *App) process(ctx context.Context) error {
	a.processMx.Lock()
	defer a.processMx.Unlock()

	log.Info().
		Str("checkPeddingBids", checkPeddingBids.String()).
		Msg("checking pedding bids")

	bids, err := a.repo.GetListByStatus(ctx, repository.StatusPedding)
	if err != nil {
		return errors.Wrap(err, "get list by status fail")
	}

	for _, bid := range bids {
		if err := a.processBid(ctx, bid); err != nil {
			return errors.Wrap(err, "process bid fail")
		}
	}

	return nil
}

func (a *App) processBid(ctx context.Context, bid repository.Bid) error {
	cl, err := a.getBlockchain(bid.GivenCoinType)
	if err != nil {
		if err := a.repo.UpdateStatus(ctx, bid.ID, repository.StatusError); err != nil {
			return errors.Wrap(err, "update status fail")
		}

		return errors.Wrap(err, "get blockchain fail")
	}

	// Упрощаю задачу. На данном этапе просто верим, что пользователь отправил нам свои монетки и не будем заниматься
	// мониторингом пришедших от пользователь монеток.
	// Для мониторинга прихода монеток от пользователя, проще выдавать пользователю на каждую заявку уникальный кошелек,
	// чтобы пользователь на нее высылал свои монетки, такми образом мониторинг прихода будет более прозначным.
	// Но на данный момент не вижу смысла усложнять проект и выдавать пользователю уникальные адреса.

	var walletFrom string

	// Если монетки получаем btc, то значит слать будет eth и поэтому указываем системный кошелек
	// для списания транзакции. Опять упрощаем задачу.
	if bid.ReceivedCoinType == consts.CoinTypeBTC {
		walletFrom, err = a.conf.GetAddress(consts.CoinTypeETH)
		if err != nil {
			if err := a.repo.UpdateStatus(ctx, bid.ID, repository.StatusError); err != nil {
				return errors.Wrap(err, "update status fail")
			}

			return errors.Wrap(err, "get address fail")
		}
	}

	// Шлем транзакцию
	txid, err := cl.Send(ctx, walletFrom, bid.GivenAddress, bid.GivenValue)
	if err != nil {
		if err := a.repo.UpdateStatus(ctx, bid.ID, repository.StatusError); err != nil {
			return errors.Wrap(err, "update status fail")
		}

		return errors.Wrap(err, "get address fail")
	}

	// Закрываем заявку
	if err := a.repo.UpdateStatus(ctx, bid.ID, repository.StatusClose); err != nil {
		return errors.Wrap(err, "update status fail")
	}

	log.Info().
		Str("txid", txid).
		Msg("sending crypto coins")

	return nil
}

func (a *App) getBlockchain(coinType uint64) (blockchains.Blockchainer, error) {
	cl, ok := a.blockchains[coinType]
	if !ok {
		//nolint:wrapcheck
		return nil, ErrBlockchainNotFound
	}

	return cl, nil
}
