package rpc

import (
	"github.com/rs/zerolog"
	"github.com/stop-error/scamjam-tls"
)

func InitTls(logger *zerolog.Logger) (caCertPEMBytes []byte, caPrivateKeyPEMBytes []byte, leafCertPEMBytes []byte, leafPrivateKeyPEMBytes []byte){
	caCertPEMBytes, caPrivateKeyPEMBytes, err := scamjamtls.GetRootCa(logger, "scamjam-service")
	if err != nil {
		logger.Error().Err(err).Msg("Error getting CA cert and private key!")
	}

	leafCertPEMBytes, leafPrivateKeyPEMBytes, err = scamjamtls.GetLeaf(logger, "scamjam-service", caCertPEMBytes, caPrivateKeyPEMBytes)

	return caCertPEMBytes, caPrivateKeyPEMBytes, leafCertPEMBytes, leafPrivateKeyPEMBytes
	
}

