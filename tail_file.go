package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

func tailFile(ctx context.Context, path string, prefix string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Could not open file: %w", err)
	}
	defer f.Close()

	go func() {
		<-ctx.Done()
		f.Close()
	}()

	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("Could not find end realtive to start: %w", err)
	}

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if err == io.EOF {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if errors.Is(err, os.ErrClosed) {
				return nil
			}
			return fmt.Errorf("Could not read string line: %w", err)
		}
		fmt.Printf("[%s] - %s", prefix, line)
	}
}

func startTailer(ctx context.Context, path, prefix string, wg *sync.WaitGroup, active map[string]context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	active[prefix] = cancel
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := tailFile(ctx, path, prefix); err != nil && err != context.Canceled {
			log.Printf("Tailer for %s: %v", prefix, err)
		}
	}()
}
