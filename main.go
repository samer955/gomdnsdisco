package main

import (
	"fmt"
	"os"
	"os/signal"
	"p2p-example/service"
	"syscall"
)

func main() {

	agent := service.NewSender()
	agent.Start()

	//Run the program till its stopped (forced)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")

}
