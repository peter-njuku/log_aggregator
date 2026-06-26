package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <logFile.log>", os.Args[0])
	}

	filename := os.Args[1]

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := tailFile(ctx, filename)
	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		log.Fatal(err)
	}

	fmt.Println("\nShuttin down")
}
