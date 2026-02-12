package rpc

import (
	"net"
	"context"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	pbalerts "github.com/stop-error/scamjam-service/alerts_proto"
)



type server struct {
    pbalerts.UnsafeDnsAlertsServer
}


func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    log.Printf("Received: %v\n", req.GetName())
    return &pb.HelloResponse{
        Message: fmt.Sprintf("Hello %s!", req.GetName()),
    }, nil
}



func InitDnsAlertsRpc(logger *zerolog.Logger) { //loop with exit channel, run as goroutine?

    const port = ":50051"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatal().Msg("failed to listen on port: " + port + err.Error())
	}

	s := grpc.NewServer()
	pbalerts.RegisterDnsAlertsServer(s, &server{})

	logger.Info().Msg("server listening at" + lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		logger.Fatal().Msg("failed to serve: " + err.Error())
	}
}