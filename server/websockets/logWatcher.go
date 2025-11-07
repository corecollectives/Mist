package websockets

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

func WatcherLogs(ctx context.Context, filePath string, send chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		send <- line
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			} else if err != nil {
				return err
			}
			send <- line
		}
	}
}
