package ddsrpc

import (
	"context"
	"github.com/892294101/dds-rpc/pcb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"time"
)

type RpcClient struct {
	rc   pcb.GreeterClient
	conn *grpc.ClientConn
}

func (r *RpcClient) newRpcClient(port string) error {
	// 连接grpc服务器
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	ctx1, canal := context.WithTimeout(context.Background(), time.Second*5)
	defer canal()
	conn, err := grpc.DialContext(ctx1, net.JoinHostPort("localhost", port), dopts...)
	if err != nil {
		return errors.Errorf("did not connect: %v", err)
	}

	// 初始化Greeter服务客户端
	c := pcb.NewGreeterClient(conn)
	r.rc = c
	r.conn = conn
	return nil

}

func (r *RpcClient) Stop() (*pcb.StopReply, error) {
	// 初始化上下文，设置请求超时时间为15秒
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	// 调用Stop接口，发送一条消息
	sr, err := r.rc.Stop(ctx, &pcb.StopCommand{Instruction: 1})
	if err != nil {
		return nil, err
	}
	return sr, nil
}

func (r *RpcClient) StopRpc() error {
	return r.conn.Close()
}

func NewRpcClient(port string) (*RpcClient, error) {
	c := new(RpcClient)
	if err := c.newRpcClient(port); err != nil {
		return nil, err
	}
	return c, nil
}
