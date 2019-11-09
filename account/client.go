package account

import (
	"github.com/micro-kit/microkit-demo/internal/pb"
	"github.com/micro-kit/microkit/client"
	"google.golang.org/grpc"
)

var (
	svcName = "account"
)

// NewClient 创建客户端
func NewClient() (accountClient pb.AccountClient, err error) {
	c, err := client.NewDefaultClient(client.ServiceName(svcName))
	if err != nil {
		return
	}
	// 连接服务端
	c.Dial(func(cc *grpc.ClientConn) {
		accountClient = pb.NewAccountClient(cc)
	})
	return
}
