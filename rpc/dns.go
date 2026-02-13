package rpc

import (
	"net"
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
	pbalerts "github.com/stop-error/scamjam-service/alerts_proto"
)

var rpcLogger *zerolog.Logger

type server struct {
    pbalerts.UnsafeDnsAlertsServer
}


func (s *server) DnsAlertFound(ctx context.Context, threat *pbalerts.DnsThreat) (*pbalerts.DnsReply, error) {
   	rpcLogger.Info().Msg("Received DnsThreat with details- Hostname: " + threat.Hostname + " Source: " + threat.Source + " ThreatType: " + threat.Category)
    return &pbalerts.DnsReply{
        Ack: true,
    }, nil //add error checking for invalid submissions
}



func InitDnsAlertsRpc(logger *zerolog.Logger) { //loop with exit channel, run as goroutine?

	rpcLogger = logger
    const port = ":50051"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal().Msg("DnsAlerts grpc server failed to listen on port: " + port + err.Error()) //TODO: Try another pot number
	}

	altsCreds := alts.NewServerCreds(alts.DefaultServerOptions())
	s := grpc.NewServer(grpc.Creds(altsCreds))
	pbalerts.RegisterDnsAlertsServer(s, &server{})

	logger.Info().Msg("DnsAlerts grpc server listening at" + lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		logger.Fatal().Msg("failed to serve: " + err.Error())
	}
}