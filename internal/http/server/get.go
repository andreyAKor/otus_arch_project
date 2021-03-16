package server

import (
	"net/http"

	"github.com/pkg/errors"

	schemaBid "github.com/andreyAKor/otus_arch_project/schema/bid"
)

// Get bid info.
func (s *Server) get(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	out, err := s.bidClient.Get(&schemaBid.GetIn{
		Id: r.URL.Query().Get("id"),
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return nil, errors.Wrap(err, "can't get bid info")
	}

	return out, nil
}
