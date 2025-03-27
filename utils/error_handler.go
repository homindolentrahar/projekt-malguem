package utils

import (
	"os"

	"github.com/rs/zerolog/log"
)

func HandleErrorReturn(err error) error {
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	return nil
}

func HandleErrorExit(err error) {
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}
