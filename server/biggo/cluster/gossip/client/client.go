package client

import (
	"fmt"
	"time"

	"git.biggo.com/Funmula/BigGoChat/server/v8/biggo/cluster/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClusterConnection(hostname string, port int32) (*grpc.ClientConn, error) {
	// address-naming: https://github.com/grpc/grpc/blob/master/doc/naming.md
	return grpc.NewClient(
		fmt.Sprintf("%s:%d", hostname, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			MinConnectTimeout: time.Second,
		}),
	)
}

var NewClusterClient = proto.NewClusterClient
