package utils

import (
	"os"

	"github.com/rs/zerolog/log"
) 

func HandleErrorExit(err error) {
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}