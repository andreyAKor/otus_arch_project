package btc

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

type Bitcoin struct {
	host string
	port int64
	user string
	pass string
}

func NewBitcoin(host string, port int64, user, pass string) *Bitcoin {
	return &Bitcoin{host, port, user, pass}
}

// Установка соединения с RPC шлюзом bitcoin.
func (b *Bitcoin) connect() (*rpcclient.Client, error) {
	// Если работаем с тестовой сетью битков
	chaincfg.MainNetParams = chaincfg.TestNet3Params

	return rpcclient.New(&rpcclient.ConnConfig{
		HTTPPostMode: true,
		DisableTLS:   true,
		Host:         fmt.Sprintf("%s:%d", b.host, b.port),
		User:         b.user,
		Pass:         b.pass,
	}, nil)
}

// Устанавливаем комсу для кошелька.
func (b *Bitcoin) SetTxFee(fee int64) error {
	client, err := b.connect()
	if err != nil {
		return errors.Wrap(err, "connect fail")
	}

	defer client.Shutdown()

	if err := client.SetTxFee(btcutil.Amount(fee)); err != nil {
		return errors.Wrap(err, "set tx fee fail")
	}

	return nil
}

// Проверка баланса на сервере с минимальным количеством подтверждений.
func (b *Bitcoin) GetBalanceMinConf(account string, minConfirms int) (int64, error) {
	client, err := b.connect()
	if err != nil {
		return 0, errors.Wrap(err, "connect fail")
	}
	defer client.Shutdown()

	amount, err := client.GetBalanceMinConf(account, minConfirms)
	if err != nil {
		return 0, errors.Wrap(err, "get balance min conf fail")
	}

	return int64(amount), nil
}

// Отправляет сумму на указанный адрес.
func (b *Bitcoin) SendToAddress(address string, amount int64) (string, error) {
	client, err := b.connect()
	if err != nil {
		return "", errors.Wrap(err, "connect fail")
	}
	defer client.Shutdown()

	res, err := client.RawRequest(
		"sendtoaddress",
		[]json.RawMessage{
			json.RawMessage(`"` + address + `"`),
			json.RawMessage(fmt.Sprintf("%f", btcutil.Amount(amount).ToBTC())),
			json.RawMessage(`""`),
			json.RawMessage(`""`),
			json.RawMessage(`true`),
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "raw request fail")
	}

	var txid string

	if err := json.Unmarshal(res, &txid); err != nil {
		return "", errors.Wrap(err, "json unmarshal fail")
	}

	return txid, nil
}
