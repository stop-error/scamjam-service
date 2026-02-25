package rpc

import (
	"net"
	"context"
	"os"
	"crypto/tls"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/billgraziano/dpapi"
	pbalerts "github.com/stop-error/scamjam-service/alerts_proto"
)

var rpcLogger *zerolog.Logger

var programData = os.Getenv("ProgramData")
var scamJamProgramDataCerts = programData + "\\ScamJam\\grpc\\dns\\tls"

type server struct {
    pbalerts.UnsafeDnsAlertsServer
}


func (s *server) DnsAlertFound(ctx context.Context, threat *pbalerts.DnsThreat) (*pbalerts.DnsReply, error) {
   	rpcLogger.Info().Msg("Received DnsThreat with details- Hostname: " + threat.Hostname + " Source: " + threat.Source + " ThreatType: " + threat.Category)
    return &pbalerts.DnsReply{
        Ack: true,
    }, nil //add error checking for invalid submissions
}


func InitDnsAlertsRpc(logger *zerolog.Logger) { //loop with exit channel, run as goroutine? TODO: clean up, better error handling

	rpcLogger = logger
    const port = ":50051"

	if _, err := os.Stat(scamJamProgramDataCerts); err != nil {

		err := os.MkdirAll(scamJamProgramDataCerts, 0600)
		if err != nil {
			logger.Error().Err(err).Msg("Error creating ScamJam ProgramData directory")
		}
	}

	caCertPEMBytes, caPrivateKeyPEMBytes, leafCertPEMBytes, leafPrivateKeyPEMBytes := InitTls(logger)

	caCertPEMBytes, err :=EncryptPEMBytesDpapi(logger, caCertPEMBytes)
	if err != nil {
		logger.Error().Err(err).Msg("Error encrypting during SaveEncryptedToDisk!")
	}
	err = os.WriteFile(scamJamProgramDataCerts + "\\scamjam-ca.pem", caCertPEMBytes, 0600)
	if err != nil {
		logger.Error().Err(err).Msg("Error writing encrypted file to ProgramData!")
	}

	caPrivateKeyPEMBytes, err = EncryptPEMBytesDpapi(logger, caPrivateKeyPEMBytes)
	if err != nil {
		logger.Error().Err(err).Msg("Error encrypting caPrivateKeyPEMBytes!")
	}


	serverCert, err := tls.X509KeyPair(leafCertPEMBytes, leafPrivateKeyPEMBytes) //error
	if err != nil {
		logger.Error().Err(err).Msg("Error creating server certificate for grpc!")
	}

	tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{serverCert},
        ClientAuth:   tls.NoClientCert,
    }

	tlsCreds := credentials.NewTLS(tlsConfig)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal().Msg("DnsAlerts grpc server failed to listen on port: " + port + err.Error()) //TODO: Try another pot number
	}

	s := grpc.NewServer(grpc.Creds(tlsCreds))
	pbalerts.RegisterDnsAlertsServer(s, &server{})

	logger.Info().Msg("DnsAlerts grpc server listening at" + lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		logger.Fatal().Msg("DnsAlerts grpc server failed to serve: " + err.Error())
	}
}


func EncryptPEMBytesDpapi(logger *zerolog.Logger, pemBytes []byte) (pemBytesDpapi []byte, err error) {
	pemBytes, err = dpapi.EncryptBytes(pemBytes)
	if err != nil {
		logger.Error().Err(err).Msg("Error protecting PEMBytes with DPAPI!")
		byteError := []byte{}
		return byteError, err
	}
	return pemBytes, nil
}

