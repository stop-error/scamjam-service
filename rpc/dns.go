package rpc

import (
	"crypto/tls"
	"io"
	"net"
	"os"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pbalerts "github.com/stop-error/scamjam-service/alerts_proto"
	"github.com/stop-error/scamjam-service/certs"
)

var rpcLogger *zerolog.Logger

type DnsAlertsServer struct {
    pbalerts.UnsafeDnsAlertServiceServer
}


func (server DnsAlertsServer) Dns(srv pbalerts.DnsAlertService_DnsServer) error {
	rpcLogger.Info().Msg("Starting DnsAlertsServer grpc stream loop")

	ctx := srv.Context()

	for {

		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			rpcLogger.Error().Err(ctx.Err()).Msg("Stream loop executed ctx.Done case! (Check the context deadline?)")
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := srv.Recv()
		if err == io.EOF {
			// return will close stream from server side
			rpcLogger.Info().Msg("Recieved io.EOF from client, shutting down grpc stream.")
			return nil
		}
		if err != nil {
			rpcLogger.Error().Err(err).Msg("Recieved error from client!")
			continue
		}

		rpcLogger.Info().Msg("Recieved message from client (blocky): " + req.Category + " " + req.Hostname + " " + req.Source)
		
	}
}


func InitDnsAlertsRpc(logger *zerolog.Logger) { //loop with exit channel, run as goroutine? TODO: clean up, better error handling

	programData := os.Getenv("ProgramData")
	scamJamProgramDataCerts := programData + "\\ScamJam\\grpc\\dns"


	rpcLogger = logger
    const port = ":50051"

	if _, err := os.Stat(scamJamProgramDataCerts); err != nil {

		err := os.MkdirAll(scamJamProgramDataCerts, 0600)
		if err != nil {
			logger.Error().Err(err).Msg("Error creating ScamJam ProgramData directory")
		}
	}

	caCertPEMBytes, caPrivateKeyPEMBytes, leafCertPEMBytes, leafPrivateKeyPEMBytes := InitDnsTls(logger)

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
	pbalerts.RegisterDnsAlertServiceServer(s, &DnsAlertsServer{})

	logger.Info().Msg("DnsAlerts grpc server listening at" + lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		logger.Fatal().Msg("DnsAlerts grpc server failed to serve: " + err.Error())
	}
}

func InitDnsTls(logger *zerolog.Logger) (caCertPEMBytes []byte, caPrivateKeyPEMBytes []byte, leafCertPEMBytes []byte, leafPrivateKeyPEMBytes []byte){
	caCertPEMBytes, caPrivateKeyPEMBytes, err := certs.GetRootCa(logger, "scamjam-service")
	if err != nil {
		logger.Error().Err(err).Msg("Error getting CA cert and private key!")
	}

	leafCertPEMBytes, leafPrivateKeyPEMBytes, err = certs.GetLeaf(logger, "scamjam-service", caCertPEMBytes, caPrivateKeyPEMBytes)

	return caCertPEMBytes, caPrivateKeyPEMBytes, leafCertPEMBytes, leafPrivateKeyPEMBytes
	
}




