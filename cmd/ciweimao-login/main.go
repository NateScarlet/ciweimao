package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/NateScarlet/ciweimao/pkg/client"
)

func main() {
	var username = os.Getenv("CIWEIMAO_USERNAME")
	var password = os.Getenv("CIWEIMAO_PASSWORD")

	if username == "" {
		log.Fatalf("environment variable CIWEIMAO_USERNAME is not defined")
	}
	if password == "" {
		log.Fatalf("environment variable CIWEIMAO_PASSWORD is not defined")
	}

	var ctx = context.Background()
	var err error
	_, err = client.Default.Login(ctx, username, password)
	if err != nil {
		log.Fatalf("login failed: %s", err)
	}

	fmt.Printf("CIWEIMAO_ACCOUNT=%s\n", client.Default.Account)
	fmt.Printf("CIWEIMAO_LOGIN_TOKEN=%s\n", client.Default.LoginToken)
	fmt.Printf("CIWEIMAO_DEVICE_TOKEN=%s\n", client.Default.DeviceToken)
}
