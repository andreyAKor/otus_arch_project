package repository

import (
	"context"
	"database/sql"
	"io"
	"time"

	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("bid not found")

type BidsRepo interface {
	Create(ctx context.Context, bid Bid) (string, error)
	Get(ctx context.Context, id string) (Bid, error)
	GetListByStatus(ctx context.Context, status uint64) ([]Bid, error)
	UpdateStatus(ctx context.Context, id string, status uint64) error
}

type DBBidsRepo interface {
	Connect(ctx context.Context, dsn string) error

	io.Closer
	BidsRepo
}

type Bid struct {
	ID               string
	Status           uint64
	ReceivedCoinType uint64
	ReceivedValue    string
	ReceivedAddress  string
	GivenCoinType    uint64
	GivenValue       string
	GivenAddress     string
	CreatedAt        time.Time
	UpdatedAt        sql.NullTime
}
