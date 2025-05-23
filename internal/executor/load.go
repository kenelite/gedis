package executor

import (
	"fmt"
	"github.com/kenelite/gedis/internal/response"
	"io"
	"strings"
)

func LoadAof(respReader *response.Resp) error {

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
		//
		cmd := strings.ToUpper(cmdVal.Bulk)
		_ = Execute(cmd, val.Array[1:])
	}

	return nil

}
