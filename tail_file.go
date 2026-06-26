package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

func tailFile(ctx context.Context, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Could not open file: %w", err)
	}
	defer f.Close()

	// TODO : seek end of file
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("Could not find end realtive to start: %w", err)
	}

	reader := bufio.NewReader(f)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// TODO: Reach EOF do something not to crash our program
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return fmt.Errorf("Could not read string line: %w", err)
		}
		os.Stdout.WriteString(line)
	}
}
