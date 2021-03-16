package eth

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"
)

type Geth struct {
	client *ethrpc.EthRPC
}

func NewGeth(ip string, port int64) *Geth {
	return &Geth{
		client: ethrpc.New(fmt.Sprintf("http://%s:%d", ip, port)),
	}
}

func (g *Geth) call(method string, target interface{}, params ...interface{}) error {
	result, err := g.client.Call(method, params...)
	if err != nil {
		return errors.Wrap(err, "call method fail")
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}

// Баланс аккаунта.
func (g *Geth) EthGetBalance(account string) (big.Int, error) {
	return g.client.EthGetBalance(account, "latest")
}

// Посылаем транзакцию (personal).
func (g *Geth) PersonalSendTransaction(transaction ethrpc.T, passphrase string) (string, error) {
	var r string
	err := g.call("personal_sendTransaction", &r, transaction, passphrase)

	return r, err
}
