package main

import (
	"fmt"
	"time"

	"log-receiver/pkg/auth"
)

func main() {
	now := time.Now().UTC().Unix()

	token := auth.GenIDPJWTToken("sao", "123", "789", "456", now+3600*24)
	fmt.Print(token)
}
