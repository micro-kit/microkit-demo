package test

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/micro-kit/microkit-client/client/account"
	"github.com/micro-kit/microkit-client/proto/accountpb"
)

var (
	cl accountpb.AccountClient
)

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var err error
	cl, err = account.NewClient()
	if err != nil {
		log.Panicln(err)
	}
	m.Run()
}

// 测试登录
func TestGRPCLogin(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	u, err := cl.Login(ctx, &accountpb.LoginRequest{
		Username: "admin",
		Password: "111111",
	})
	if err != nil {
		log.Println(err)
		return
	}
	js, _ := json.Marshal(u)
	log.Println(string(js))

}

// 测试多个请求
func TestMultipleGRPCLogin(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(100)
	for j := 0; j < 100; j++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				u, err := cl.Login(ctx, &accountpb.LoginRequest{
					Username: "admin",
					Password: "111111",
				})
				defer cancel()
				if err != nil {
					log.Println(err)
					return
				}
				js, _ := json.Marshal(u)
				log.Println(string(js))
			}
		}()
	}
	wg.Wait()
}
