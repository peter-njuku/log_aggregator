package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <dir>", os.Args[0])
	}

	dir := os.Args[1]

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fd, err := unix.InotifyInit()
	if err != nil {
		log.Fatal(err)
	}
	defer unix.Close(fd)

	_, err = unix.InotifyAddWatch(fd, dir, unix.IN_CREATE|unix.IN_DELETE)
	if err != nil {
		log.Fatal(err)
	}

	active := make(map[string]context.CancelFunc)
	var wg sync.WaitGroup

	files, err := filepath.Glob(filepath.Join(dir, "*.log"))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		name := filepath.Base(f)
		startTailer(f, name, &wg, active)
	}

	go func() {
		<-ctx.Done()
		unix.Close(fd)
	}()

	buf := make([]byte, 4096)
	for {
		n, err := unix.Read(fd, buf)
		if err != nil {
			if err == syscall.EBADF {
				break
			}
			if err == syscall.EINTR {
				continue
			}
			log.Fatal(err)
		}

		var i uint32
		for i < uint32(n) {
			event := (*unix.InotifyEvent)(unsafe.Pointer(&buf[i]))

			var name string
			if event.Len > 0 {
				nameBytes := buf[i+unix.SizeofInotifyEvent : i+unix.SizeofInotifyEvent+uint32(event.Len)]
				name = string(nameBytes[:len(nameBytes)-1])
			}
			if event.Mask&unix.IN_CREATE != 0 {
				if matched, _ := filepath.Match(".log", name); matched {
					if _, exists := active[name]; !exists {
						startTailer(dir, name, &wg, active)
					}
				}
			}

			if event.Mask&unix.IN_DELETE != 0 {
				if cancel, ok := active[name]; ok {
					cancel()
					delete(active, name)
				}
			}
			i += unix.SizeofInotifyEvent + uint32(event.Len)
		}
	}

	fmt.Println("\nShuttin down...")
	for _, cancel := range active {
		cancel()
	}
	wg.Wait()
	fmt.Println("All trailers stopped")
}
