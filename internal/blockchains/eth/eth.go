package eth

import (
	"context"
	"math/big"

	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/otus_arch_project/internal/blockchains"
	"github.com/andreyAKor/otus_arch_project/internal/configs"
	"github.com/andreyAKor/otus_arch_project/internal/consts"
)

//nolint:gosec
const (
	estimateGas   = 21000             // Объем затрачиваемого газа на транзакцию
	passphrase    = "S3CsQmbaJd"      // Пароль для системного кошелька. Упрощаем задачу.
	avgTransPrice = "210000000000000" // Стоимость газа для совершения транзакции. Упрощаем задачу.
)

var (
	ErrAvgBFEmpty   = errors.New("Avg BF is empty")
	ErrFeeBFEmpty   = errors.New("Fee BF is empty")
	ErrValueBIEmpty = errors.New("Value BI is empty")

	_ blockchains.Blockchainer = (*eth)(nil)
)

type eth struct {
	conf configs.Configer
}

func New(conf configs.Configer) (blockchains.Blockchainer, error) {
	return &eth{conf}, nil
}

func (e *eth) CalcTransaction(_ context.Context, value string) (blockchains.TransactionPrice, error) {
	valueBI, _ := new(big.Int).SetString(value, 10)
	if valueBI == nil {
		return blockchains.TransactionPrice{}, errors.Wrap(ErrValueBIEmpty, "set string fail")
	}

	estimateGasBI := new(big.Int).SetUint64(estimateGas)

	return blockchains.TransactionPrice{
		Value:     new(big.Int).Sub(valueBI, estimateGasBI).String(),
		Fee:       avgTransPrice,
		AvgTxSize: estimateGas,
	}, nil
}

func (e *eth) Send(
	_ context.Context,
	walletFrom, walletTo, value string,
) (string, error) {
	log.Info().
		Str("walletFrom", walletFrom).
		Str("walletTo", walletTo).
		Str("value", value).
		Msg("send eth coins")

	// Преобразуем string в big.Float
	feeBF, _ := new(big.Float).SetString(avgTransPrice)
	if feeBF == nil {
		return "", errors.Wrap(ErrFeeBFEmpty, "set string fail")
	}

	// Преобразуем string в big.Int
	valueBI, _ := new(big.Int).SetString(value, 10)
	if valueBI == nil {
		return "", errors.Wrap(ErrValueBIEmpty, "set string fail")
	}

	// Необходимое количество газа в big.Float и bit.Int
	estimateGasBF := new(big.Float).SetUint64(estimateGas)

	// Формула:
	// 	gasPriceBF = feeBF / estimateGasBF
	gasPriceBF := new(big.Float).Quo(feeBF, estimateGasBF)

	// Приводим big.Float в big.Int
	gasPriceBI := new(big.Int)
	gasPriceBF.Int(gasPriceBI)

	txid, err := e.client().PersonalSendTransaction(ethrpc.T{
		From:     walletFrom,
		To:       walletTo,
		Value:    valueBI,
		Gas:      estimateGas,
		GasPrice: gasPriceBI,
		Data:     "",
		Nonce:    0,
	}, passphrase)
	if err != nil {
		return "", errors.Wrap(err, "personal send transaction fail")
	}

	log.Info().
		Str("walletFrom", walletFrom).
		Str("walletTo", walletTo).
		Str("value", value).
		Str("txid", txid).
		Msg("eth coins is sending")

	return txid, nil
}

func (e *eth) SummaryBalance(_ context.Context) (string, error) {
	addr, err := e.conf.GetAddress(consts.CoinTypeETH)
	if err != nil {
		return "", errors.Wrap(err, "get address fail")
	}

	balance, err := e.client().EthGetBalance(addr)
	if err != nil {
		return "", errors.Wrap(err, "eth get balance fail")
	}

	return balance.String(), nil
}

func (e *eth) client() *Geth {
	return NewGeth(
		e.conf.GetNodes().ETH.Host,
		e.conf.GetNodes().ETH.Port,
	)
}
