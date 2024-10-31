package server

import (
	"context"

	"github.com/ammm56/lings/cmd/lingswallet/daemon/pb"
	"github.com/ammm56/lings/version"
)

func (s *server) GetVersion(_ context.Context, _ *pb.GetVersionRequest) (*pb.GetVersionResponse, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return &pb.GetVersionResponse{
		Version: version.Version(),
	}, nil
}
