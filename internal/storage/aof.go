package storage

import (
	"bufio"
	"fmt"
	"github.com/kenelite/gedis/internal/handle"
	"github.com/kenelite/gedis/internal/response"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.mu.Lock()

			aof.file.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *Aof) Write(value response.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Load() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	// Seek to the beginning of the file
	if _, err := aof.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	aof.rd = bufio.NewReader(aof.file)
	respReader := response.NewResp(aof.rd)

	for {
		val, err := respReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read AOF entry: %w", err)
		}

		// Expecting top-level value to be an array: [command, arg1, arg2, ...]
		if val.Typ != "array" || len(val.Array) == 0 {
			continue // skip invalid or malformed entries
		}

		cmdVal := val.Array[0]
		if cmdVal.Typ != "bulk" {
			continue // skip if command is not a bulk string
		}

		cmd := strings.ToUpper(cmdVal.Bulk)
		handler, ok := handle.Handlers[cmd]
		if !ok {
			continue // unknown command
		}

		// Execute the command to restore in-memory state
		handler(val.Array[1:])
	}

	return nil
}
