package main

import (
	"abc/test"
	"context"
	"fmt"
)

func main() {
	endpoint := []string{"52.11.26.186:2379", "52.11.26.186:22379", "52.11.26.186:3237"}
	rpcClient, err := test.NewUserServiceClient(endpoint)
	if err != nil {
		panic(err)
	}

	req := new(test.LoginRequest)
	req.Name = []byte("1234")
	req.Pwd = []byte("abc")
	resp, err := rpcClient.Login(context.Background(), req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result is %v\n", resp.Success)

}
