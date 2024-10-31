package client

import (
	"context"
	"time"

	"github.com/ammm56/lings/cmd/lingswallet/daemon/server"

	"github.com/pkg/errors"

	"github.com/ammm56/lings/cmd/lingswallet/daemon/pb"
	"google.golang.org/grpc"
)

// Connect connects to the lingswalletd server, and returns the client instance
func Connect(address string) (pb.HtnwalletdClient, func(), error) {
	// Connection is local, so 1 second timeout is sufficient
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(server.MaxDaemonSendMsgSize)))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, errors.New("lingswallet daemon is not running, start it with `lingswallet start-daemon`")
		}
		return nil, nil, err
	}

	return pb.NewHtnwalletdClient(conn), func() {
		conn.Close()
	}, nil
}
