package gedis_go_sdk

import (
	"fmt"
	"strings"
)

func formatRespArray(command string, args ...string) string {
	all := append([]string{command}, args...)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(all)))
	for _, arg := range all {
		sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}
	return sb.String()
}
