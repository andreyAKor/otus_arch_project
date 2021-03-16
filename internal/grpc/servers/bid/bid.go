package bid

import (
	"context"
	"math/big"

	"github.com/micro/go-micro/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/andreyAKor/otus_arch_project/internal/blockchains"
	configs "github.com/andreyAKor/otus_arch_project/internal/configs/bid"
	"github.com/andreyAKor/otus_arch_project/internal/consts"
	"github.com/andreyAKor/otus_arch_project/internal/repository"
	schemaBid "github.com/andreyAKor/otus_arch_project/schema/bid"
)

const rateFromBTCtoETH = 32.127414 // Курс валют BTC-ETH на момент 14.03.2021

var (
	ErrBlockchainNotFound = errors.New("blockchain not found")
	ErrTPValueBIEmpty     = errors.New("transaction price value big int is empty")
	ErrValueBFEmpty       = errors.New("value big float is empty")
	ErrSummaryBIEmpty     = errors.New("summary big int is empty")
	ErrValueLessSummary   = errors.New("value less than summary")
	ErrDivisorBFEmpty     = errors.New("divisor big float is empty")
)

type Bid struct {
	srv         server.Server
	repo        repository.DBBidsRepo
	conf        *configs.Config
	blockchains map[uint64]blockchains.Blockchainer
}

func New(
	srv server.Server,
	repo repository.DBBidsRepo,
	conf *configs.Config,
	blockchains map[uint64]blockchains.Blockchainer,
) (*Bid, error) {
	return &Bid{srv, repo, conf, blockchains}, nil
}

func (b *Bid) Run(_ context.Context) error {
	if err := schemaBid.RegisterBidHandler(b.srv, b); err != nil {
		return errors.Wrap(err, "register bid handler fail")
	}

	return nil
}

//nolint:funlen
func (b *Bid) Create(ctx context.Context, in *schemaBid.CreateIn, out *schemaBid.CreateOut) error {
	log.Info().
		Str("value", in.Value).
		Str("address", in.Address).
		Uint64("coinType", in.CoinType).
		Msg("creating a bid")

	givenCoinType := b.getGivenCoinType(in.CoinType)
	log.Info().
		Uint64("givenCoinType", givenCoinType).
		Msg("get given coin type result")

	ratedValue, err := calcExchange(in.CoinType, givenCoinType, in.Value, rateFromBTCtoETH)
	if err != nil {
		return errors.Wrap(err, "calculate exchange value fail")
	}
	log.Info().
		Str("ratedValue", ratedValue).
		Msg("calc exchange result")

	givenValue, err := b.calcGivenValue(ratedValue)
	if err != nil {
		return errors.Wrap(err, "calculate received value fail")
	}
	log.Info().
		Str("givenValue", givenValue).
		Msg("calc given value result")

	valid, err := b.validateReceivedValue(ctx, givenCoinType, givenValue)
	if err != nil {
		log.Info().
			Err(err).
			Msg("creating validation fail")

		if valid {
			return errors.Wrap(err, "validate received value fail")
		}

		out.Header = &schemaBid.OutHeader{Status: schemaBid.OutStatus_VALIDATION, Message: err.Error()}

		return nil
	}

	receivedAddress, err := b.getReceivedAddress(in.CoinType)
	if err != nil {
		return errors.Wrap(err, "get received address fail")
	}

	id, err := b.repo.Create(ctx, repository.Bid{
		Status:           repository.StatusPending,
		ReceivedCoinType: in.CoinType,
		ReceivedValue:    in.Value,
		ReceivedAddress:  receivedAddress,
		GivenCoinType:    givenCoinType,
		GivenValue:       givenValue,
		GivenAddress:     in.Address,
	})
	if err != nil {
		return errors.Wrap(err, "create repository fail")
	}

	out.Header = &schemaBid.OutHeader{Status: schemaBid.OutStatus_OK}
	out.Id = id

	return nil
}

func (b *Bid) Get(ctx context.Context, in *schemaBid.GetIn, out *schemaBid.GetOut) error {
	bid, err := b.repo.Get(ctx, in.Id)
	if err != nil {
		return errors.Wrap(err, "create repository fail")
	}

	out.Header = &schemaBid.OutHeader{Status: schemaBid.OutStatus_OK}
	out.Status = bid.Status
	out.ReceivedCoinType = bid.ReceivedCoinType
	out.ReceivedValue = bid.ReceivedValue
	out.ReceivedAddress = bid.ReceivedAddress
	out.GivenCoinType = bid.GivenCoinType
	out.GivenValue = bid.GivenValue
	out.GivenAddress = bid.GivenAddress

	return nil
}

func (b *Bid) validateReceivedValue(ctx context.Context, coinType uint64, value string) (bool, error) {
	cl, err := b.getBlockchain(coinType)
	if err != nil {
		return true, errors.Wrap(err, "get blockchain fail")
	}

	tp, err := cl.CalcTransaction(ctx, value)
	if err != nil {
		return true, errors.Wrap(err, "calc transaction fail")
	}

	log.Info().
		Str("value", tp.Value).
		Str("fee", tp.Fee).
		Int("avgTxSize", tp.AvgTxSize).
		Msg("calc transaction")

	tpValueBI, _ := new(big.Int).SetString(tp.Value, 10)
	if tpValueBI == nil {
		return true, errors.Wrap(ErrTPValueBIEmpty, "set string fail")
	}

	summary, err := cl.SummaryBalance(ctx)
	if err != nil {
		return true, errors.Wrap(err, "blockchain summary balance fail")
	}

	summaryBI, _ := new(big.Int).SetString(summary, 10)
	if summaryBI == nil {
		return true, errors.Wrap(ErrSummaryBIEmpty, "set string fail")
	}

	if tpValueBI.Cmp(summaryBI) > 0 {
		//nolint:wrapcheck
		return false, ErrValueLessSummary
	}

	return true, nil
}

func (b Bid) getGivenCoinType(coinType uint64) uint64 {
	if coinType == consts.CoinTypeETH {
		return consts.CoinTypeBTC
	}

	return consts.CoinTypeETH
}

func (b *Bid) getReceivedAddress(coinType uint64) (string, error) {
	addr, err := b.conf.GetAddress(coinType)
	if err != nil {
		return "", errors.Wrap(err, "get address fail")
	}

	return addr, nil
}

func (b *Bid) getBlockchain(coinType uint64) (blockchains.Blockchainer, error) {
	cl, ok := b.blockchains[coinType]
	if !ok {
		//nolint:wrapcheck
		return nil, ErrBlockchainNotFound
	}

	return cl, nil
}

// calcGivenValue расчет отдаваемой валюты с учетом процента обменника.
func (b Bid) calcGivenValue(receivedValue string) (string, error) {
	valueBF, _ := new(big.Float).SetString(receivedValue)
	if valueBF == nil {
		return "", errors.Wrap(ErrValueBFEmpty, "set string fail")
	}

	basisBF := new(big.Float).SetUint64(uint64(100))
	percentBF := new(big.Float).SetUint64(uint64(b.conf.Fee))

	// Формулы:
	// 	valueByPercentBF = valueBF / 100 - считаем сколько крипты приходится на 1 процент
	// 	feeBF = valueByPercentBF * percentBF - счтаем комисию обменника от количества крипты
	// 	totalBF = valueBF - feeBF - считаем итоговое количество крипты на отправку юзеру
	valueByPercentBF := new(big.Float).Quo(valueBF, basisBF)
	feeBF := new(big.Float).Mul(valueByPercentBF, percentBF)
	totalBF := new(big.Float).Sub(valueBF, feeBF)

	// Приводим big.Float в big.Int
	totalBI := new(big.Int)
	totalBF.Int(totalBI)

	return totalBI.String(), nil
}

// calcExchange расчет отдаваемой валюты с учетом курска BTC-ETH.
func calcExchange(fromCoinType, toCoinType uint64, fromValue string, rate float64) (string, error) {
	valueBF, err := normalize(fromCoinType, fromValue)
	if err != nil {
		return "", errors.Wrap(err, "normalize fail")
	}

	rateBF := new(big.Float).SetFloat64(rate)

	ratedValueBF := new(big.Float)
	if fromCoinType == consts.CoinTypeBTC {
		ratedValueBF = ratedValueBF.Mul(valueBF, rateBF)
	} else {
		ratedValueBF = ratedValueBF.Quo(valueBF, rateBF)
	}

	resValue, err := denormalize(toCoinType, ratedValueBF)
	if err != nil {
		return "", errors.Wrap(err, "denormalize fail")
	}

	return resValue, nil
}

// normalize возвращает значение монеты в нормализованном формате.
func normalize(coinType uint64, value string) (*big.Float, error) {
	divisor := getDivisor(coinType)

	divisorBF, _ := new(big.Float).SetString(divisor)
	if divisorBF == nil {
		return nil, errors.Wrap(ErrDivisorBFEmpty, "set string fail")
	}

	valueBF, _ := new(big.Float).SetString(value)
	if valueBF == nil {
		return nil, errors.Wrap(ErrValueBFEmpty, "set string fail")
	}

	normalizedValueBF := new(big.Float).Quo(valueBF, divisorBF)

	return normalizedValueBF, nil
}

// denormalize возвращает значение монеты в денормализованном формате.
func denormalize(coinType uint64, valueBF *big.Float) (string, error) {
	if valueBF == nil {
		//nolint:wrapcheck
		return "", ErrValueBFEmpty
	}

	divisor := getDivisor(coinType)

	divisorBF, _ := new(big.Float).SetString(divisor)
	if divisorBF == nil {
		return "", errors.Wrap(ErrDivisorBFEmpty, "set string fail")
	}

	normalizedValueBF := new(big.Float).Mul(valueBF, divisorBF)

	// Приводим big.Float в big.Int
	normalizedValueBI := new(big.Int)
	normalizedValueBF.Int(normalizedValueBI)

	return normalizedValueBI.String(), nil
}

func getDivisor(coinType uint64) string {
	if coinType == consts.CoinTypeBTC {
		return "100000000"
	}

	return "1000000000000000000"
}
