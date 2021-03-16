-- +goose Up
-- +goose StatementBegin
CREATE TABLE bids
(
    id                 varchar(64)  NOT NULL,
    status             int2         NOT NULL,
    received_coin_type int2         NOT NULL,
    received_value     varchar(64)  NOT NULL,
    received_address   varchar(128) NOT NULL,
    given_coin_type    int2         NOT NULL,
    given_value        varchar(64)  NOT NULL,
    given_address      varchar(128) NOT NULL,
    created_at         timestamptz  NOT NULL DEFAULT now(),
    updated_at         timestamptz  NULL
);
CREATE UNIQUE INDEX bids_id_idx ON bids USING btree (id);
CREATE INDEX bids_status_idx ON bids USING btree (status);
CREATE INDEX bids_received_coin_type_idx ON bids USING btree (received_coin_type);
CREATE INDEX bids_given_coin_type_idx ON bids USING btree (given_coin_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bids;
-- +goose StatementEnd
