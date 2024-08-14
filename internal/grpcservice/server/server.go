package server

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"

	app "github.com/cronfy/trainer/internal/app/domain"
	grpc "github.com/cronfy/trainer/internal/grpcservice/domain"
	pb "github.com/cronfy/trainer/internal/grpcservice/generated"
)

type Server struct {
	pb.UnimplementedTrainerServer
	multiplyTaskUC grpc.MultiplyTaskUseCase
}

func New(multiplyTaskUC grpc.MultiplyTaskUseCase) *Server {
	return &Server{
		multiplyTaskUC: multiplyTaskUC,
	}
}

func (s *Server) GetMultiplyTask(ctx context.Context, req *empty.Empty) (*pb.MultiplyTask, error) {
	task := s.multiplyTaskUC.Get()

	return &pb.MultiplyTask{
		Operands: []int32{int32(task.Operands[0]), int32(task.Operands[1])},
	}, nil
}

func (s *Server) SolveMultiplyTask(ctx context.Context, solution *pb.MultiplyTaskSolution) (*pb.SolutionResult, error) {
	if len(solution.Task.Operands) != 2 {
		return nil, fmt.Errorf("invalid number of task operands, expected 2, got %d", len(solution.Task.Operands))
	}

	task := app.MultiplyTask{Operands: []int{int(solution.Task.Operands[0]), int(solution.Task.Operands[1])}}
	correct := s.multiplyTaskUC.Solve(task, int(solution.Solution))

	return &pb.SolutionResult{
		Correct: correct,
	}, nil
}
