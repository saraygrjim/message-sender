package main

import (
	"fmt"
	"message-sender/cmd"
	"os"
)

func startReceiver() {
	cmd.Receiver()
}

func startBroadcaster() {
	cmd.Broadcaster()
}

func startSubscriber() {
	cmd.StartSubscriber()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <command>")
		fmt.Println("Available commands:")
		fmt.Println("  receiver - Starts Microservice 1")
		fmt.Println("  broadcaster - Starts Microservice 2")
		fmt.Println("  subscriber    - Start a subscriber")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "receiver":
		startReceiver()
	case "broadcaster":
		startBroadcaster()
	case "subscriber":
		startSubscriber()
	default:
		fmt.Println("❌ Comand not supported:", command)
		fmt.Println("Available commands: receiver, broadcaster, subscriber")
		os.Exit(1)
	}
}
