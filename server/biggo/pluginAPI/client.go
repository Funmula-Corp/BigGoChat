package pluginAPI

import (
	"context"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/shared/mlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newClient() (client PluginAPIClient, connection *grpc.ClientConn, err error) {
	if connection, err = grpc.NewClient("localhost:9999", grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return
	}
	client = NewPluginAPIClient(connection)
	return
}

func GetUserIdByAuthData(authData string) (result *string, err error) {
	var (
		client     PluginAPIClient
		connection *grpc.ClientConn
	)

	if client, connection, err = newClient(); err != nil {
		return
	}
	defer connection.Close()

	var reply *UserIdByAuthDataReply
	if reply, err = client.GetUserIdByAuthData(context.Background(), &UserIdByAuthDataRequest{AuthData: authData}); err != nil {
		mlog.Error("GetUserIdByAuthData", mlog.Err(err))
		return
	}
	result = &reply.UserId
	return
}
