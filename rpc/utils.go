package rpc

import (

	"github.com/billgraziano/dpapi"
	"github.com/rs/zerolog"
) 

func EncryptPEMBytesDpapi(logger *zerolog.Logger, pemBytes []byte) (pemBytesDpapi []byte, err error) {
	pemBytes, err = dpapi.EncryptBytes(pemBytes)
	if err != nil {
		logger.Error().Err(err).Msg("Error protecting PEMBytes with DPAPI!")
		byteError := []byte{}
		return byteError, err
	}
	return pemBytes, nil
}