package server

import (
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	bid "github.com/andreyAKor/otus_arch_project/schema/bid"
)

var ErrCoinTypeParsingFail = errors.New("coin type parsing fail")

// Create new bid.
func (s *Server) create(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	coinType, err := strconv.Atoi(r.PostFormValue("coin_type"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		//nolint:wrapcheck
		return nil, ErrCoinTypeParsingFail
	}

	out, err := s.bidClient.Create(&bid.CreateIn{
		CoinType: uint64(coinType),
		Value:    r.PostFormValue("value"),
		Address:  r.PostFormValue("address"),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return nil, errors.Wrap(err, "can't create new bid")
	}
	if out.Header != nil && out.Header.Status == bid.OutStatus_VALIDATION {
		w.WriteHeader(http.StatusBadRequest)

		return nil, errors.New(out.Header.Message)
	}

	w.WriteHeader(http.StatusCreated)

	return out, nil
}
