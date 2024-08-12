package grpcservice

import (
	"github.com/cronfy/trainer/internal/grpcservice/domain"
	"github.com/cronfy/trainer/internal/grpcservice/server"
)

func Build(multiplyTaskUC domain.MultiplyTaskUseCase) *server.Server {
	return server.New(multiplyTaskUC)
}
