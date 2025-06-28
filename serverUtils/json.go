package serverUtils

import (
	"encoding/json"
	"io"

	"github.com/Alquama00s/serverEngine/lib"
	"github.com/rs/zerolog/log"
)

func Unmarshal[T any](r *lib.Request) (*T, error) {
	var result T

	if len(r.RawBody) == 0 && r.RawRequest.Body != nil {

		bodyByte, err := io.ReadAll(r.RawRequest.Body)

		if err != nil {
			log.Error().Err(err).Msg("Failed to read request body")
			return nil, err
		}
		r.RawBody = bodyByte

	}

	if len(r.RawBody) > 0 {
		json.Unmarshal(r.RawBody, &result)

		r.Body = &result
	}

	log.Info().Interface("result", result)
	return &result, nil
}
