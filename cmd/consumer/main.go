package main

import "github.com/third-place/user-service/internal/kafka"

func main() {
	kafka.InitializeAndRunLoop()
}
