package main

import (
	"fmt"
	"os"

	"example.com/mdns-broadcast/services"
)

func help() {
	fmt.Println("Should pass a param to specify which service will be running")
	fmt.Println("\n s : stands for source")
	fmt.Println("\n d : stands for destination")
}

// Main function
func main() {
	if len(os.Args) != 2 {
		help()
	}

	switch os.Args[1] {
	case "s":
		services.SourceMain()
	case "d":
		services.DestinationMain()
	default:
		help()
	}

}
