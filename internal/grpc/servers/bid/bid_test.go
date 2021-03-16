package bid

import (
	"math/big"
	"testing"

	"github.com/andreyAKor/otus_arch_project/internal/consts"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest
func TestNormalize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("btc", func(t *testing.T) {
			res, err := normalize(consts.CoinTypeBTC, "100000000") // 1 BTC
			require.NoError(t, err)

			require.Equal(t, "1", res.String())
		})
		t.Run("eth", func(t *testing.T) {
			res, err := normalize(consts.CoinTypeETH, "1000000000000000000") // 1 ETH
			require.NoError(t, err)

			require.Equal(t, "1", res.String())
		})
	})
}

//nolint:paralleltest
func TestDenormalize(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("btc", func(t *testing.T) {
			res, err := denormalize(consts.CoinTypeBTC, new(big.Float).SetFloat64(1))
			require.NoError(t, err)

			require.Equal(t, "100000000", res) // 1 BTC
		})
		t.Run("btc", func(t *testing.T) {
			res, err := denormalize(consts.CoinTypeETH, new(big.Float).SetFloat64(1))
			require.NoError(t, err)

			require.Equal(t, "1000000000000000000", res) // 1 ETH
		})
	})
}

//nolint:paralleltest
func TestCalcExchange(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("btc to eth", func(t *testing.T) {
			res, err := calcExchange(consts.CoinTypeBTC, consts.CoinTypeETH, "20000000", 0.5) // 0.2 BTC by rate 0.5 BTC-ETH
			require.NoError(t, err)

			require.Equal(t, "100000000000000000", res) // 0.1 ETH
		})
		t.Run("eth to btc", func(t *testing.T) {
			res, err := calcExchange(consts.CoinTypeETH, consts.CoinTypeBTC, "200000000000000000", 0.5) // 0.2 ETH by rate 0.5 BTC-ETH
			require.NoError(t, err)

			require.Equal(t, "40000000", res) // 0.1 BTC
		})
	})
}
