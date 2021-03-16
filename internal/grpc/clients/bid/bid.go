package bid

import (
	"context"
	"time"

	"github.com/micro/go-micro/client"
	"github.com/pkg/errors"

	"github.com/andreyAKor/otus_arch_project/internal/consts"
	bid "github.com/andreyAKor/otus_arch_project/schema/bid"
)

const requestTimeout = 15

type Bid struct {
	cl client.Client
}

func New(cl client.Client) (*Bid, error) {
	return &Bid{cl}, nil
}

func (b *Bid) Create(in *bid.CreateIn) (*bid.CreateOut, error) {
	out, err := bid.NewBidService(consts.BidService, b.cl).Create(
		context.Background(),
		in,
		client.WithRequestTimeout(time.Second*requestTimeout),
	)
	if err != nil {
		return nil, errors.Wrap(err, "create request fail")
	}

	return out, nil
}

func (b *Bid) Get(in *bid.GetIn) (*bid.GetOut, error) {
	out, err := bid.NewBidService(consts.BidService, b.cl).Get(
		context.Background(),
		in,
		client.WithRequestTimeout(time.Second*requestTimeout),
	)
	if err != nil {
		return nil, errors.Wrap(err, "get request fail")
	}

	return out, nil
}
