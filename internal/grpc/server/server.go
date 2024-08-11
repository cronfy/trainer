package server

import (
	pb "oap/trainer/internal/grpc/generated"
)

type Server struct {
	pb.UnimplementedTrainerServer
}

func New() *Server {
	return &Server{}
}
