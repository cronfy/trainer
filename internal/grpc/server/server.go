package server

import (
	"context"
	"math/rand"

	"github.com/golang/protobuf/ptypes/empty"

	pb "oap/trainer/internal/grpc/generated"
)

type Server struct {
	pb.UnimplementedTrainerServer
}

func New() *Server {
	return &Server{}
}

func (s *Server) GetMultiplyTask(ctx context.Context, req *empty.Empty) (*pb.MultiplyTask, error) {
	operands := make([]int32, 2)

	operands[0] = randomInt31Between(1, 20)
	operands[1] = randomInt31Between(1, 11)

	if rand.Int31n(2) == 0 {
		operands[0], operands[1] = operands[1], operands[0]
	}

	return &pb.MultiplyTask{
		Operands: operands,
	}, nil
}

func (s *Server) SolveMultiplyTask(ctx context.Context, solution *pb.MultiplyTaskSolution) (*pb.SolutionResult, error) {
	correctSolution := solution.Task.Operands[0] * solution.Task.Operands[1]

	return &pb.SolutionResult{
		Correct: solution.Solution == correctSolution,
	}, nil
}

func randomInt31Between(min, max int32) int32 {
	return rand.Int31n(max-min) + min
}
