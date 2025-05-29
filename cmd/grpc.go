package cmd

import (
	"ewallet-transaction/helpers"
	"log"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func ServeGRPC() {
	lis, err := net.Listen("tcp", ":"+helpers.GetEnv("GRPC_PORT", "7000"))
	if err != nil {
		log.Fatal("Failed to serve grpc port:", err)
	}

	s := grpc.NewServer()
	logrus.Info("Starting listening grpc on port:", helpers.GetEnv("GRPC_PORT", "7000"))
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve grpc port:", err)
	}
}
