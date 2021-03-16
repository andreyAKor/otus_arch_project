package psql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib" // Register pgx
	"github.com/pkg/errors"

	"github.com/andreyAKor/otus_arch_project/internal/repository"
)

var _ repository.DBBidsRepo = (*repo)(nil)

type repo struct {
	db *sql.DB
}

func New() (repository.DBBidsRepo, error) {
	return &repo{}, nil
}

// Create connection pool.
func (r *repo) Connect(ctx context.Context, dsn string) (err error) {
	r.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return
	}

	return r.db.PingContext(ctx)
}

// Close connection pool.
func (r *repo) Close() error {
	return r.db.Close()
}

// Add new bid.
func (r *repo) Create(ctx context.Context, bid repository.Bid) (string, error) {
	var id string
	query := `INSERT INTO
bids (id, status, received_coin_type, received_value, received_address, given_coin_type, given_value, given_address)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (id) DO NOTHING
RETURNING id`
	err := r.db.QueryRowContext(
		ctx, query,
		uuid.New().String(),
		bid.Status,
		bid.ReceivedCoinType,
		bid.ReceivedValue,
		bid.ReceivedAddress,
		bid.GivenCoinType,
		bid.GivenValue,
		bid.GivenAddress,
	).Scan(&id)
	if err != nil {
		return "", errors.Wrap(err, "create fail")
	}

	return id, nil
}

// Update bid status by id.
func (r *repo) UpdateStatus(ctx context.Context, id string, status uint64) error {
	query := "UPDATE bids SET status = $1, updated_at = $2 WHERE id = $3"
	res, err := r.db.ExecContext(
		ctx, query,
		status,
		"now()",
		id,
	)
	if err != nil {
		return errors.Wrap(err, "exec context fail")
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "update status fail")
	}
	if ra == 0 {
		return errors.Wrap(repository.ErrNotFound, "update status fail")
	}

	return nil
}

// Get bid by id.
func (r *repo) Get(ctx context.Context, id string) (repository.Bid, error) {
	var bid repository.Bid
	query := `
SELECT
	id,
	status,
	received_coin_type,
	received_value,
	received_address,
	given_coin_type,
	given_value,
	given_address,
	created_at,
	updated_at
FROM
	bids
WHERE
	id = $1`
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(
			&bid.ID,
			&bid.Status,
			&bid.ReceivedCoinType,
			&bid.ReceivedValue,
			&bid.ReceivedAddress,
			&bid.GivenCoinType,
			&bid.GivenValue,
			&bid.GivenAddress,
			&bid.CreatedAt,
			&bid.UpdatedAt,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return bid, errors.Wrap(repository.ErrNotFound, "get fail")
	} else if err != nil {
		return bid, errors.Wrap(err, "get fail")
	}

	return bid, nil
}

// Get list by status.
func (r *repo) GetListByStatus(ctx context.Context, status uint64) ([]repository.Bid, error) {
	query := `
SELECT
	id,
	status,
	received_coin_type,
	received_value,
	received_address,
	given_coin_type,
	given_value,
	given_address,
	created_at,
	updated_at
FROM
	bids
WHERE
	status = $1`
	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, errors.Wrap(err, "query context fail")
	}
	defer rows.Close()

	var events []repository.Bid

	for rows.Next() {
		var bid repository.Bid

		if err := rows.Scan(
			&bid.ID,
			&bid.Status,
			&bid.ReceivedCoinType,
			&bid.ReceivedValue,
			&bid.ReceivedAddress,
			&bid.GivenCoinType,
			&bid.GivenValue,
			&bid.GivenAddress,
			&bid.CreatedAt,
			&bid.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "scan fail")
		}

		events = append(events, bid)
	}

	return events, rows.Err()
}
