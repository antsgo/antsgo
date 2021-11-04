package server

import (
	"fmt"
	"net"

	"github.com/antsgo/antsgo/conf"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func NewGrpc(c conf.Conf, logger *logrus.Logger, registerGrpcHandler func(*grpc.Server)) {
	addr := fmt.Sprintf(":%d", c.ServerGrpc.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("Faild to listen", err)
	}
	s := grpc.NewServer()
	registerGrpcHandler(s)
	fmt.Printf("Serving Grpc on 0.0.0.0%s\n", addr)
	logger.Fatal(s.Serve(lis))
}
