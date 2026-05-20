package grpc_client

import (
	"context"
	"fmt"
	"time"

	corepb "proto/core"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CoreClient struct {
	client corepb.CoreServiceClient
	conn   *grpc.ClientConn
}

func NewCoreClient(host string, port int) (*CoreClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &CoreClient{
		client: corepb.NewCoreServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *CoreClient) Close() error {
	return c.conn.Close()
}

func (c *CoreClient) RegisterUser(
	ctx context.Context,
	name string,
	email string,
	password string,
) (*corepb.User, error) {
	resp, err := c.client.RegisterUser(ctx, &corepb.RegisterUserRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

func (c *CoreClient) LoginUser(
	ctx context.Context,
	email string,
	password string,
) (*corepb.User, error) {
	resp, err := c.client.LoginUser(ctx, &corepb.LoginUserRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}

func (c *CoreClient) GetUser(ctx context.Context, id int64) (*corepb.User, error) {
	resp, err := c.client.GetUser(ctx, &corepb.GetUserRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.User, nil
}
