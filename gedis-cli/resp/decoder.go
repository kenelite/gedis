package resp

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func Decode(r *bufio.Reader) (string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	switch line[0] {
	case '+', '-', ':':
		return strings.TrimSpace(line[1:]), nil
	case '$':
		length, _ := strconv.Atoi(strings.TrimSpace(line[1:]))
		if length == -1 {
			return "(nil)", nil
		}
		buf := make([]byte, length+2)
		r.Read(buf)
		return string(buf[:length]), nil
	case '*':
		count, _ := strconv.Atoi(strings.TrimSpace(line[1:]))
		result := []string{}
		for i := 0; i < count; i++ {
			item, _ := Decode(r)
			result = append(result, item)
		}
		return fmt.Sprintf("[%s]", strings.Join(result, ", ")), nil
	default:
		return "", fmt.Errorf("unknown response: %s", line)
	}
}
