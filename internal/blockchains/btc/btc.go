package btc

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/otus_arch_project/internal/blockchains"
	"github.com/andreyAKor/otus_arch_project/internal/configs"
)

const (
	avgTransPrice = 4096 // Стоимость транзакции в Сатоши. Упрощаем задачу.
	avgTxSize     = 1024 // Размер транзакции. Упрощаем задачу.
)

var ErrValueBIEmpty = errors.New("Value BI is empty")

var _ blockchains.Blockchainer = (*btc)(nil)

type btc struct {
	conf configs.Configer
}

func New(conf configs.Configer) (blockchains.Blockchainer, error) {
	return &btc{conf}, nil
}

func (b *btc) CalcTransaction(_ context.Context, value string) (blockchains.TransactionPrice, error) {
	valueBI, _ := new(big.Int).SetString(value, 10)
	if valueBI == nil {
		return blockchains.TransactionPrice{}, errors.Wrap(ErrValueBIEmpty, "set string fail")
	}

	avgTransPriceBI := new(big.Int).SetUint64(avgTransPrice)

	return blockchains.TransactionPrice{
		Value:     new(big.Int).Sub(valueBI, avgTransPriceBI).String(),
		Fee:       fmt.Sprintf("%d", avgTransPrice),
		AvgTxSize: avgTxSize,
	}, nil
}

func (b *btc) Send(
	_ context.Context,
	_, walletTo, value string,
) (string, error) {
	log.Info().
		Str("walletTo", walletTo).
		Str("value", value).
		Msg("send btc coins")

	// Устанавливаем комсу для кошелька
	if err := b.client().SetTxFee(avgTransPrice); err != nil {
		return "", errors.Wrap(err, "set tx fee fail")
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return "", errors.Wrap(err, "a to i fail")
	}

	// Отправляет сумму на указанный адрес
	txid, err := b.client().SendToAddress(walletTo, int64(val))
	if err != nil {
		return "", errors.Wrap(err, "send to address fail")
	}

	log.Info().
		Str("walletTo", walletTo).
		Str("value", value).
		Str("txid", txid).
		Msg("btc coins is sending")

	return txid, nil
}

func (b *btc) SummaryBalance(_ context.Context) (string, error) {
	amount, err := b.client().GetBalanceMinConf("*", 1)
	if err != nil {
		return "", errors.Wrap(err, "get balance min confirmation fail")
	}

	return fmt.Sprintf("%d", amount), nil
}

func (b *btc) client() *Bitcoin {
	return NewBitcoin(
		b.conf.GetNodes().BTC.Host,
		b.conf.GetNodes().BTC.Port,
		b.conf.GetNodes().BTC.User,
		b.conf.GetNodes().BTC.Pass,
	)
}
