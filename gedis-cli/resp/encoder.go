package resp

import (
	"bufio"
	"fmt"
)

func EncodeArray(w *bufio.Writer, args []string) error {
	_, err := fmt.Fprintf(w, "*%d\r\n", len(args))
	if err != nil {
		return err
	}
	for _, arg := range args {
		_, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len(arg), arg)
		if err != nil {
			return err
		}
	}
	return nil
}
